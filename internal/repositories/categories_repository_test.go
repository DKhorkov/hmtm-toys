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

	"github.com/DKhorkov/hmtm-toys/internal/repositories"
)

const (
	driver = "sqlite3"
	//dsn    = "file::memory:?cache=shared"
	dsn              = "../../test.db"
	migrationsDir    = "/migrations"
	gooseZeroVersion = 0
)

func TestCategoriesRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(CategoriesRepositoryTestSuite))
}

type CategoriesRepositoryTestSuite struct {
	suite.Suite

	cwd                  string
	ctx                  context.Context
	dbConnector          db.Connector
	connection           *sql.Conn
	categoriesRepository *repositories.CategoriesRepository
	logger               *mocklogging.MockLogger
	traceProvider        *mocktracing.MockProvider
	spanConfig           tracing.SpanConfig
}

func (s *CategoriesRepositoryTestSuite) SetupSuite() {
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
	s.categoriesRepository = repositories.NewCategoriesRepository(s.dbConnector, s.logger, s.traceProvider, s.spanConfig)
}

func (s *CategoriesRepositoryTestSuite) SetupTest() {
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

func (s *CategoriesRepositoryTestSuite) TearDownTest() {
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

func (s *CategoriesRepositoryTestSuite) TearDownSuite() {
	s.NoError(s.dbConnector.Close())
}

func (s *CategoriesRepositoryTestSuite) TestGetAllCategoriesWithExisting() {
	s.traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(context.Background(), mocktracing.NewMockSpan()).
		Times(1)

	createdAt := time.Now().UTC()
	_, err := s.connection.ExecContext(
		s.ctx,
		"INSERT INTO categories (id, name, created_at, updated_at) VALUES (?, ?, ?, ?), (?, ?, ?, ?)",
		1, "Category 1", createdAt, createdAt,
		2, "Category 2", createdAt, createdAt,
	)
	s.NoError(err)

	categories, err := s.categoriesRepository.GetAllCategories(s.ctx)
	s.NoError(err)
	s.NotEmpty(categories)
	s.Equal(2, len(categories))
	s.Equal(uint32(1), categories[0].ID)
	s.Equal("Category 1", categories[0].Name)
	s.WithinDuration(createdAt, categories[0].CreatedAt, time.Second)
	s.WithinDuration(createdAt, categories[0].UpdatedAt, time.Second)
	s.Equal(uint32(2), categories[1].ID)
	s.Equal("Category 2", categories[1].Name)
}

func (s *CategoriesRepositoryTestSuite) TestGetAllCategoriesWithoutExisting() {
	s.traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(context.Background(), mocktracing.NewMockSpan()).
		Times(1)

	categories, err := s.categoriesRepository.GetAllCategories(s.ctx)
	s.NoError(err)
	s.Empty(categories)
}

func (s *CategoriesRepositoryTestSuite) TestGetCategoryByIDExisting() {
	s.traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(context.Background(), mocktracing.NewMockSpan()).
		Times(1)

	createdAt := time.Now().UTC()
	_, err := s.connection.ExecContext(
		s.ctx,
		"INSERT INTO categories (id, name, created_at, updated_at) VALUES (?, ?, ?, ?)",
		1, "Test Category", createdAt, createdAt,
	)
	s.NoError(err)

	category, err := s.categoriesRepository.GetCategoryByID(s.ctx, 1)
	s.NoError(err)
	s.NotNil(category)
	s.Equal(uint32(1), category.ID)
	s.Equal("Test Category", category.Name)
	s.WithinDuration(createdAt, category.CreatedAt, time.Second)
	s.WithinDuration(createdAt, category.UpdatedAt, time.Second)
}

func (s *CategoriesRepositoryTestSuite) TestGetCategoryByIDNonExisting() {
	s.traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(context.Background(), mocktracing.NewMockSpan()).
		Times(1)

	category, err := s.categoriesRepository.GetCategoryByID(s.ctx, 999)
	s.Error(err)
	s.Nil(category)
}
