package repositories_test

import (
	"context"
	"database/sql"
	"os"
	"path"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3" // Must be imported for correct work

	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"github.com/DKhorkov/libs/db"
	mocklogging "github.com/DKhorkov/libs/logging/mocks"
	"github.com/DKhorkov/libs/tracing"
	mocktracing "github.com/DKhorkov/libs/tracing/mocks"

	"github.com/DKhorkov/hmtm-toys/internal/entities"
	"github.com/DKhorkov/hmtm-toys/internal/repositories"
	"github.com/DKhorkov/libs/pointers"
)

func TestMastersRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(MastersRepositoryTestSuite))
}

type MastersRepositoryTestSuite struct {
	suite.Suite

	cwd               string
	ctx               context.Context
	dbConnector       db.Connector
	connection        *sql.Conn
	mastersRepository *repositories.MastersRepository
	logger            *mocklogging.MockLogger
	traceProvider     *mocktracing.MockProvider
	spanConfig        tracing.SpanConfig
}

func (s *MastersRepositoryTestSuite) SetupSuite() {
	s.NoError(goose.SetDialect(driver))

	ctrl := gomock.NewController(s.T())
	s.ctx = context.Background()
	s.logger = mocklogging.NewMockLogger(ctrl)
	dbConnector, err := db.New(dsn, driver, s.logger)
	s.NoError(err)

	cwd, err := os.Getwd()
	s.NoError(err)

	s.cwd = cwd
	s.dbConnector = dbConnector
	s.traceProvider = mocktracing.NewMockProvider(ctrl)
	s.spanConfig = tracing.SpanConfig{}
	s.mastersRepository = repositories.NewMastersRepository(s.dbConnector, s.logger, s.traceProvider, s.spanConfig)
}

func (s *MastersRepositoryTestSuite) SetupTest() {
	s.NoError(
		goose.Up(
			s.dbConnector.Pool(),
			path.Dir(
				path.Dir(s.cwd),
			)+migrationsDir,
		),
	)

	connection, err := s.dbConnector.Connection(s.ctx)
	s.NoError(err)

	s.connection = connection
}

func (s *MastersRepositoryTestSuite) TearDownTest() {
	s.NoError(
		goose.DownTo(
			s.dbConnector.Pool(),
			path.Dir(
				path.Dir(s.cwd),
			)+migrationsDir,
			gooseZeroVersion,
		),
	)

	s.NoError(s.connection.Close())
}

func (s *MastersRepositoryTestSuite) TearDownSuite() {
	s.NoError(s.dbConnector.Close())
}

func (s *MastersRepositoryTestSuite) TestGetAllMastersWithExisting() {
	s.traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(context.Background(), mocktracing.NewMockSpan()).
		Times(1)

	createdAt := time.Now().UTC()
	info1 := pointers.New("Master Info 1")
	info2 := pointers.New("Master Info 2")
	_, err := s.connection.ExecContext(
		s.ctx,
		"INSERT INTO masters (id, user_id, info, created_at, updated_at) "+
			"VALUES (?, ?, ?, ?, ?), (?, ?, ?, ?, ?)",
		1, 1, info1, createdAt, createdAt,
		2, 2, info2, createdAt, createdAt,
	)
	s.NoError(err)

	masters, err := s.mastersRepository.GetAllMasters(s.ctx)
	s.NoError(err)
	s.NotEmpty(masters)
	s.Equal(2, len(masters))
	s.Equal(uint64(1), masters[0].UserID)
	s.NotNil(masters[0].Info)
	s.Equal(*info1, *masters[0].Info)
	s.WithinDuration(createdAt, masters[0].CreatedAt, time.Second)
	s.WithinDuration(createdAt, masters[0].UpdatedAt, time.Second)
	s.Equal(uint64(2), masters[1].UserID)
	s.NotNil(masters[1].Info)
	s.Equal(*info2, *masters[1].Info)
}

func (s *MastersRepositoryTestSuite) TestGetAllMastersWithoutExisting() {
	s.traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(context.Background(), mocktracing.NewMockSpan()).
		Times(1)

	masters, err := s.mastersRepository.GetAllMasters(s.ctx)
	s.NoError(err)
	s.Empty(masters)
}

func (s *MastersRepositoryTestSuite) TestGetMasterByIDExisting() {
	s.traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(context.Background(), mocktracing.NewMockSpan()).
		Times(1)

	createdAt := time.Now().UTC()
	info := pointers.New("Test Master Info")
	_, err := s.connection.ExecContext(
		s.ctx,
		"INSERT INTO masters (id, user_id, info, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		1, 100, info, createdAt, createdAt,
	)
	s.NoError(err)

	master, err := s.mastersRepository.GetMasterByID(s.ctx, 1)
	s.NoError(err)
	s.NotNil(master)
	s.Equal(uint64(1), master.ID)
	s.Equal(uint64(100), master.UserID)
	s.NotNil(master.Info)
	s.Equal(*info, *master.Info)
	s.WithinDuration(createdAt, master.CreatedAt, time.Second)
	s.WithinDuration(createdAt, master.UpdatedAt, time.Second)
}

func (s *MastersRepositoryTestSuite) TestGetMasterByIDNonExisting() {
	s.traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(context.Background(), mocktracing.NewMockSpan()).
		Times(1)

	master, err := s.mastersRepository.GetMasterByID(s.ctx, 999)
	s.Error(err)
	s.Nil(master)
}

func (s *MastersRepositoryTestSuite) TestGetMasterByUserIDExisting() {
	s.traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(context.Background(), mocktracing.NewMockSpan()).
		Times(1)

	createdAt := time.Now().UTC()
	info := pointers.New("Test Master Info")
	masterID := uint64(1)
	userID := uint64(1)
	_, err := s.connection.ExecContext(
		s.ctx,
		"INSERT INTO masters (id, user_id, info, created_at, updated_at) "+
			"VALUES (?, ?, ?, ?, ?)",
		masterID, userID, info, createdAt, createdAt,
	)
	s.NoError(err)

	master, err := s.mastersRepository.GetMasterByUserID(s.ctx, userID)
	s.NoError(err)
	s.NotNil(master)
	s.Equal(userID, master.UserID)
	s.NotNil(master.Info)
	s.Equal(*info, *master.Info)
	s.WithinDuration(createdAt, master.CreatedAt, time.Second)
	s.WithinDuration(createdAt, master.UpdatedAt, time.Second)
}

func (s *MastersRepositoryTestSuite) TestGetMasterByUserIDNonExisting() {
	s.traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(context.Background(), mocktracing.NewMockSpan()).
		Times(1)

	master, err := s.mastersRepository.GetMasterByUserID(s.ctx, 999)
	s.Error(err)
	s.Nil(master)
}

func (s *MastersRepositoryTestSuite) TestRegisterMasterSuccess() {
	s.traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(context.Background(), mocktracing.NewMockSpan()).
		Times(1)

	info := pointers.New("New Master Info")
	masterData := entities.RegisterMasterDTO{
		UserID: 100,
		Info:   info,
	}

	// Error and zero id due to returning nil ID after insert operation
	// SQLite inner realization without AUTO_INCREMENT for SERIAL PRIMARY KEY
	id, err := s.mastersRepository.RegisterMaster(s.ctx, masterData)
	s.Error(err)
	s.Zero(id)
}

func (s *MastersRepositoryTestSuite) TestUpdateMasterSuccess() {
	s.traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(context.Background(), mocktracing.NewMockSpan()).
		Times(1)

	createdAt := time.Now().UTC()
	_, err := s.connection.ExecContext(
		s.ctx,
		"INSERT INTO masters (id, user_id, info, created_at, updated_at) "+
			"VALUES (?, ?, ?, ?, ?)",
		1, 100, "Old Info", createdAt, createdAt,
	)
	s.NoError(err)

	newInfo := pointers.New("Updated Master Info")
	masterData := entities.UpdateMasterDTO{
		ID:   1,
		Info: newInfo,
	}

	err = s.mastersRepository.UpdateMaster(s.ctx, masterData)
	s.NoError(err)

	// Проверка masters
	rows, err := s.connection.QueryContext(s.ctx, "SELECT info FROM masters WHERE id = ?", 1)
	s.NoError(err)

	defer func() {
		s.NoError(rows.Close())
	}()

	s.True(rows.Next())
	var infoVal sql.NullString
	s.NoError(rows.Scan(&infoVal))
	s.True(infoVal.Valid)
	s.Equal(*newInfo, infoVal.String)
}

func (s *MastersRepositoryTestSuite) TestUpdateMasterNullInfo() {
	s.traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(context.Background(), mocktracing.NewMockSpan()).
		Times(1)

	createdAt := time.Now().UTC()
	_, err := s.connection.ExecContext(
		s.ctx,
		"INSERT INTO masters (id, user_id, info, created_at, updated_at) "+
			"VALUES (?, ?, ?, ?, ?)",
		1, 100, "Old Info", createdAt, createdAt,
	)
	s.NoError(err)

	masterData := entities.UpdateMasterDTO{
		ID:   1,
		Info: nil,
	}

	err = s.mastersRepository.UpdateMaster(s.ctx, masterData)
	s.NoError(err)

	// Проверка masters
	rows, err := s.connection.QueryContext(s.ctx, "SELECT info FROM masters WHERE id = ?", 1)
	s.NoError(err)

	defer func() {
		s.NoError(rows.Close())
	}()

	s.True(rows.Next())
	var infoVal sql.NullString
	s.NoError(rows.Scan(&infoVal))
	s.False(infoVal.Valid)
}
