//go:build integration

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

func TestToysRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(ToysRepositoryTestSuite))
}

type ToysRepositoryTestSuite struct {
	suite.Suite

	cwd            string
	ctx            context.Context
	dbConnector    db.Connector
	connection     *sql.Conn
	toysRepository *repositories.ToysRepository
	logger         *mocklogging.MockLogger
	traceProvider  *mocktracing.MockProvider
	spanConfig     tracing.SpanConfig
}

func (s *ToysRepositoryTestSuite) SetupSuite() {
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
	s.toysRepository = repositories.NewToysRepository(s.dbConnector, s.logger, s.traceProvider, s.spanConfig)
}

func (s *ToysRepositoryTestSuite) SetupTest() {
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

func (s *ToysRepositoryTestSuite) TearDownTest() {
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

func (s *ToysRepositoryTestSuite) TearDownSuite() {
	s.NoError(s.dbConnector.Close())
}

func (s *ToysRepositoryTestSuite) TestGetToysWithExistingToys() {
	s.traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(context.Background(), mocktracing.NewMockSpan()).
		Times(5) // Основной + 2x(getToyTags + getToyAttachments)

	createdAt := time.Now().UTC()
	_, err := s.connection.ExecContext(
		s.ctx,
		"INSERT INTO toys (id, master_id, category_id, name, description, price, quantity, created_at, updated_at) "+
			"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?), (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		1, 1, 2, "Toy 1", "Desc 1", 99.99, 5, createdAt, createdAt,
		2, 2, 3, "Toy 2", "Desc 2", 49.99, 3, createdAt, createdAt,
	)
	s.NoError(err)

	_, err = s.connection.ExecContext(
		s.ctx,
		"INSERT INTO tags (id, name, created_at, updated_at) "+
			"VALUES (?, ?, ?, ?), (?, ?, ?, ?)",
		10, "Tag 1", createdAt, createdAt,
		20, "Tag 2", createdAt, createdAt,
	)
	s.NoError(err)

	_, err = s.connection.ExecContext(
		s.ctx,
		"INSERT INTO toys_tags_associations (id, toy_id, tag_id) "+
			"VALUES (?, ?, ?), (?, ?, ?)",
		1, 1, 10,
		2, 1, 20,
	)
	s.NoError(err)

	_, err = s.connection.ExecContext(
		s.ctx,
		"INSERT INTO toys_attachments (id, toy_id, link, created_at, updated_at) "+
			"VALUES (?, ?, ?, ?, ?)",
		1, 1, "file1.jpg", createdAt, createdAt,
	)
	s.NoError(err)

	toys, err := s.toysRepository.GetToys(s.ctx, nil)
	s.NoError(err)
	s.NotEmpty(toys)
	s.Equal(2, len(toys))
	s.Equal(uint64(1), toys[0].MasterID)
	s.Equal(uint32(2), toys[0].CategoryID)
	s.Equal("Toy 1", toys[0].Name)
	s.Equal("Desc 1", toys[0].Description)
	s.InDelta(99.99, toys[0].Price, 0.01)
	s.Equal(uint32(5), toys[0].Quantity)
	s.Equal(2, len(toys[0].Tags))
	s.Contains([]string{toys[0].Tags[0].Name, toys[0].Tags[1].Name}, "Tag 1")
	s.Contains([]string{toys[0].Tags[0].Name, toys[0].Tags[1].Name}, "Tag 2")
	s.Equal(1, len(toys[0].Attachments))
	s.Equal("file1.jpg", toys[0].Attachments[0].Link)
}

func (s *ToysRepositoryTestSuite) TestGetToysWithExistingToysAndPagination() {
	s.traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(context.Background(), mocktracing.NewMockSpan()).
		Times(1)

	createdAt := time.Now().UTC()
	_, err := s.connection.ExecContext(
		s.ctx,
		"INSERT INTO toys (id, master_id, category_id, name, description, price, quantity, created_at, updated_at) "+
			"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?), (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		1, 1, 2, "Toy 1", "Desc 1", 99.99, 5, createdAt, createdAt,
		2, 2, 3, "Toy 2", "Desc 2", 49.99, 3, createdAt, createdAt,
	)
	s.NoError(err)

	_, err = s.connection.ExecContext(
		s.ctx,
		"INSERT INTO tags (id, name, created_at, updated_at) "+
			"VALUES (?, ?, ?, ?), (?, ?, ?, ?)",
		10, "Tag 1", createdAt, createdAt,
		20, "Tag 2", createdAt, createdAt,
	)
	s.NoError(err)

	_, err = s.connection.ExecContext(
		s.ctx,
		"INSERT INTO toys_tags_associations (id, toy_id, tag_id) "+
			"VALUES (?, ?, ?), (?, ?, ?)",
		1, 1, 10,
		2, 1, 20,
	)
	s.NoError(err)

	_, err = s.connection.ExecContext(
		s.ctx,
		"INSERT INTO toys_attachments (id, toy_id, link, created_at, updated_at) "+
			"VALUES (?, ?, ?, ?, ?)",
		1, 1, "file1.jpg", createdAt, createdAt,
	)
	s.NoError(err)

	pagination := &entities.Pagination{
		Limit:  pointers.New[uint64](1),
		Offset: pointers.New[uint64](2),
	}

	toys, err := s.toysRepository.GetToys(s.ctx, pagination)
	s.NoError(err)
	s.Empty(toys)
}

func (s *ToysRepositoryTestSuite) TestGetToysWithoutExisting() {
	s.traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(context.Background(), mocktracing.NewMockSpan()).
		Times(1)

	toys, err := s.toysRepository.GetToys(s.ctx, nil)
	s.NoError(err)
	s.Empty(toys)
}

func (s *ToysRepositoryTestSuite) TestCountToysWithExistingToys() {
	s.traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(context.Background(), mocktracing.NewMockSpan()).
		Times(1)

	createdAt := time.Now().UTC()
	_, err := s.connection.ExecContext(
		s.ctx,
		"INSERT INTO toys (id, master_id, category_id, name, description, price, quantity, created_at, updated_at) "+
			"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?), (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		1, 1, 2, "Toy 1", "Desc 1", 99.99, 5, createdAt, createdAt,
		2, 2, 3, "Toy 2", "Desc 2", 49.99, 3, createdAt, createdAt,
	)
	s.NoError(err)

	count, err := s.toysRepository.CountToys(s.ctx)
	s.NoError(err)
	s.Equal(uint64(2), count)
}

func (s *ToysRepositoryTestSuite) TestCountToysWithoutExisting() {
	s.traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(context.Background(), mocktracing.NewMockSpan()).
		Times(1)

	count, err := s.toysRepository.CountToys(s.ctx)
	s.NoError(err)
	s.Zero(count)
}

func (s *ToysRepositoryTestSuite) TestGetMasterToysWithExistingToys() {
	s.traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(context.Background(), mocktracing.NewMockSpan()).
		Times(5) // Основной + 2x(getToyTags + getToyAttachments)

	masterID := uint64(1)
	createdAt := time.Now().UTC()
	_, err := s.connection.ExecContext(
		s.ctx,
		"INSERT INTO toys (id, master_id, category_id, name, description, price, quantity, created_at, updated_at) "+
			"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?), (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		1, masterID, 2, "Toy 1", "Desc 1", 99.99, 5, createdAt, createdAt,
		2, masterID, 3, "Toy 2", "Desc 2", 49.99, 3, createdAt, createdAt,
	)
	s.NoError(err)

	toys, err := s.toysRepository.GetMasterToys(s.ctx, masterID, nil)
	s.NoError(err)
	s.NotEmpty(toys)
	s.Equal(2, len(toys))
	s.Equal(masterID, toys[0].MasterID)
	s.Equal(masterID, toys[1].MasterID)
}

func (s *ToysRepositoryTestSuite) TestGetMasterToysWithoutExistingToys() {
	s.traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(context.Background(), mocktracing.NewMockSpan()).
		Times(1)

	toys, err := s.toysRepository.GetMasterToys(s.ctx, 999, nil)
	s.NoError(err)
	s.Empty(toys)
}

func (s *ToysRepositoryTestSuite) TestGetMasterToysWithExistingToysAndPagination() {
	s.traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(context.Background(), mocktracing.NewMockSpan()).
		Times(1)

	masterID := uint64(1)
	createdAt := time.Now().UTC()
	_, err := s.connection.ExecContext(
		s.ctx,
		"INSERT INTO toys (id, master_id, category_id, name, description, price, quantity, created_at, updated_at) "+
			"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?), (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		1, masterID, 2, "Toy 1", "Desc 1", 99.99, 5, createdAt, createdAt,
		2, masterID, 3, "Toy 2", "Desc 2", 49.99, 3, createdAt, createdAt,
	)
	s.NoError(err)

	pagination := &entities.Pagination{
		Limit:  pointers.New[uint64](1),
		Offset: pointers.New[uint64](2),
	}

	toys, err := s.toysRepository.GetMasterToys(s.ctx, masterID, pagination)
	s.NoError(err)
	s.Empty(toys)
}

func (s *ToysRepositoryTestSuite) TestGetToyByIDExisting() {
	s.traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(context.Background(), mocktracing.NewMockSpan()).
		Times(3) // Основной + getToyTags + getToyAttachments

	createdAt := time.Now().UTC()
	_, err := s.connection.ExecContext(
		s.ctx,
		"INSERT INTO toys (id, master_id, category_id, name, description, price, quantity, created_at, updated_at) "+
			"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		1, 1, 2, "Test Toy", "Test Description", 99.99, 5, createdAt, createdAt,
	)
	s.NoError(err)
	_, err = s.connection.ExecContext(
		s.ctx,
		"INSERT INTO tags (id, name, created_at, updated_at) "+
			"VALUES (?, ?, ?, ?), (?, ?, ?, ?)",
		10, "Tag 1", createdAt, createdAt,
		20, "Tag 2", createdAt, createdAt,
	)
	s.NoError(err)
	_, err = s.connection.ExecContext(
		s.ctx,
		"INSERT INTO toys_tags_associations (id, toy_id, tag_id) "+
			"VALUES (?, ?, ?), (?, ?, ?)",
		1, 1, 10,
		2, 1, 20,
	)
	s.NoError(err)
	_, err = s.connection.ExecContext(
		s.ctx,
		"INSERT INTO toys_attachments (id, toy_id, link, created_at, updated_at) "+
			"VALUES (?, ?, ?, ?, ?)",
		1, 1, "file1.jpg", createdAt, createdAt,
	)
	s.NoError(err)

	toy, err := s.toysRepository.GetToyByID(s.ctx, 1)
	s.NoError(err)
	s.NotNil(toy)
	s.Equal(uint64(1), toy.ID)
	s.Equal(uint64(1), toy.MasterID)
	s.Equal(uint32(2), toy.CategoryID)
	s.Equal("Test Toy", toy.Name)
	s.Equal("Test Description", toy.Description)
	s.InDelta(99.99, toy.Price, 0.01)
	s.Equal(uint32(5), toy.Quantity)
	s.Equal(2, len(toy.Tags))
	s.Contains([]string{toy.Tags[0].Name, toy.Tags[1].Name}, "Tag 1")
	s.Contains([]string{toy.Tags[0].Name, toy.Tags[1].Name}, "Tag 2")
	s.Equal(1, len(toy.Attachments))
	s.Equal("file1.jpg", toy.Attachments[0].Link)
}

func (s *ToysRepositoryTestSuite) TestGetToyByIDNonExisting() {
	s.traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(context.Background(), mocktracing.NewMockSpan()).
		Times(1)

	toy, err := s.toysRepository.GetToyByID(s.ctx, 999)
	s.Error(err)
	s.Nil(toy)
}

func (s *ToysRepositoryTestSuite) TestAddToySuccess() {
	s.traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(context.Background(), mocktracing.NewMockSpan()).
		Times(1)

	toyData := entities.AddToyDTO{
		MasterID:    1,
		CategoryID:  2,
		Name:        "Test Toy",
		Description: "Test Description",
		Price:       99.99,
		Quantity:    5,
		TagIDs:      []uint32{10, 20},
		Attachments: []string{"file1.jpg", "file2.pdf"},
	}

	// Error and zero id due to returning nil ID after insert operation
	// SQLite inner realization without AUTO_INCREMENT for SERIAL PRIMARY KEY
	id, err := s.toysRepository.AddToy(s.ctx, toyData)
	s.Error(err)
	s.Zero(id)
}

func (s *ToysRepositoryTestSuite) TestDeleteToySuccess() {
	s.traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(context.Background(), mocktracing.NewMockSpan()).
		Times(1)

	createdAt := time.Now().UTC()
	_, err := s.connection.ExecContext(
		s.ctx,
		"INSERT INTO toys (id, master_id, category_id, name, description, price, quantity, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		1, 1, 2, "Test Toy", "Test Description", 99.99, 5, createdAt, createdAt,
	)
	s.NoError(err)

	err = s.toysRepository.DeleteToy(s.ctx, 1)
	s.NoError(err)

	rows, err := s.connection.QueryContext(s.ctx, "SELECT id FROM toys WHERE id = ?", 1)
	s.NoError(err)

	defer func() {
		s.NoError(rows.Close())
	}()

	s.False(rows.Next())
}

func (s *ToysRepositoryTestSuite) TestUpdateToyFullUpdate() {
	s.traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(context.Background(), mocktracing.NewMockSpan()).
		Times(1)

	s.logger.
		EXPECT().
		ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
		Times(1)

	createdAt := time.Now().UTC()
	_, err := s.connection.ExecContext(
		s.ctx,
		"INSERT INTO toys (id, master_id, category_id, name, description, price, quantity, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		1, 1, 2, "Old Toy", "Old Desc", 50.00, 1, createdAt, createdAt,
	)
	s.NoError(err)
	_, err = s.connection.ExecContext(
		s.ctx,
		"INSERT INTO toys_tags_associations (toy_id, tag_id) VALUES (?, ?)",
		1, 10,
	)
	s.NoError(err)
	_, err = s.connection.ExecContext(
		s.ctx,
		"INSERT INTO toys_attachments (id, toy_id, link, created_at, updated_at) "+
			"VALUES (?, ?, ?, ?, ?)",
		1, 1, "oldfile.jpg", createdAt, createdAt,
	)
	s.NoError(err)

	newCategoryID := uint32(3)
	newName := "Updated Toy"
	newDesc := "Updated Desc"
	newPrice := pointers.New[float32](150.00)
	newQuantity := uint32(10)
	toyData := entities.UpdateToyDTO{
		ID:                    1,
		CategoryID:            &newCategoryID,
		Name:                  &newName,
		Description:           &newDesc,
		Price:                 newPrice,
		Quantity:              &newQuantity,
		TagIDsToAdd:           []uint32{30, 40},
		TagIDsToDelete:        []uint32{10},
		AttachmentsToAdd:      []string{"newfile.jpg"},
		AttachmentIDsToDelete: []uint64{1},
	}

	err = s.toysRepository.UpdateToy(s.ctx, toyData)
	s.NoError(err)

	// Проверка toys
	rows, err := s.connection.QueryContext(s.ctx, "SELECT master_id, category_id, name, description, price, quantity FROM toys WHERE id = ?", 1)
	s.NoError(err)

	defer func() {
		s.NoError(rows.Close())
	}()

	s.True(rows.Next())
	var masterID uint64
	var categoryID uint32
	var name, description string
	var price float32
	var quantity uint32
	s.NoError(rows.Scan(&masterID, &categoryID, &name, &description, &price, &quantity))
	s.Equal(uint64(1), masterID)
	s.Equal(newCategoryID, categoryID)
	s.Equal(newName, name)
	s.Equal(newDesc, description)
	s.InDelta(*newPrice, price, 0.01)
	s.Equal(newQuantity, quantity)
}

func (s *ToysRepositoryTestSuite) TestUpdateToyMinimalUpdate() {
	s.traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(context.Background(), mocktracing.NewMockSpan()).
		Times(1)

	s.logger.
		EXPECT().
		ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
		Times(1)

	createdAt := time.Now().UTC()
	_, err := s.connection.ExecContext(
		s.ctx,
		"INSERT INTO toys (id, master_id, category_id, name, description, price, quantity, created_at, updated_at) "+
			"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		1, 1, 2, "Old Toy", "Old Desc", 50.00, 1, createdAt, createdAt,
	)
	s.NoError(err)

	newPrice := pointers.New[float32](75.00)
	toyData := entities.UpdateToyDTO{
		ID:    1,
		Price: newPrice,
	}

	err = s.toysRepository.UpdateToy(s.ctx, toyData)
	s.NoError(err)

	// Проверка toys
	rows, err := s.connection.QueryContext(s.ctx, "SELECT master_id, category_id, name, description, price, quantity FROM toys WHERE id = ?", 1)
	s.NoError(err)

	defer func() {
		s.NoError(rows.Close())
	}()

	s.True(rows.Next())
	var masterID uint64
	var categoryID uint32
	var name, description string
	var price float32
	var quantity uint32
	s.NoError(rows.Scan(&masterID, &categoryID, &name, &description, &price, &quantity))
	s.Equal(uint64(1), masterID)
	s.Equal(uint32(2), categoryID)
	s.Equal("Old Toy", name)
	s.Equal("Old Desc", description)
	s.InDelta(*newPrice, price, 0.01)
	s.Equal(uint32(1), quantity)
}
