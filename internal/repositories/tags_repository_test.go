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
)

func TestTagsRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(TagsRepositoryTestSuite))
}

type TagsRepositoryTestSuite struct {
	suite.Suite

	cwd            string
	ctx            context.Context
	dbConnector    db.Connector
	connection     *sql.Conn
	tagsRepository *repositories.TagsRepository
	logger         *mocklogging.MockLogger
	traceProvider  *mocktracing.MockProvider
	spanConfig     tracing.SpanConfig
}

func (s *TagsRepositoryTestSuite) SetupSuite() {
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
	s.tagsRepository = repositories.NewTagsRepository(s.dbConnector, s.logger, s.traceProvider, s.spanConfig)
}

func (s *TagsRepositoryTestSuite) SetupTest() {
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

func (s *TagsRepositoryTestSuite) TearDownTest() {
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

func (s *TagsRepositoryTestSuite) TearDownSuite() {
	s.NoError(s.dbConnector.Close())
}

func (s *TagsRepositoryTestSuite) TestGetAllTagsWithExisting() {
	s.traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(context.Background(), mocktracing.NewMockSpan()).
		Times(1)

	createdAt := time.Now().UTC()
	_, err := s.connection.ExecContext(
		s.ctx,
		"INSERT INTO tags (id, name, created_at, updated_at) VALUES (?, ?, ?, ?), (?, ?, ?, ?)",
		1, "Tag 1", createdAt, createdAt,
		2, "Tag 2", createdAt, createdAt,
	)
	s.NoError(err)

	tags, err := s.tagsRepository.GetAllTags(s.ctx)
	s.NoError(err)
	s.NotEmpty(tags)
	s.Equal(2, len(tags))
	s.Equal(uint32(1), tags[0].ID)
	s.Equal("Tag 1", tags[0].Name)
	s.WithinDuration(createdAt, tags[0].CreatedAt, time.Second)
	s.WithinDuration(createdAt, tags[0].UpdatedAt, time.Second)
	s.Equal(uint32(2), tags[1].ID)
	s.Equal("Tag 2", tags[1].Name)
}

func (s *TagsRepositoryTestSuite) TestGetAllTagsWithoutExisting() {
	s.traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(context.Background(), mocktracing.NewMockSpan()).
		Times(1)

	tags, err := s.tagsRepository.GetAllTags(s.ctx)
	s.NoError(err)
	s.Empty(tags)
}

func (s *TagsRepositoryTestSuite) TestGetTagByIDExisting() {
	s.traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(context.Background(), mocktracing.NewMockSpan()).
		Times(1)

	createdAt := time.Now().UTC()
	_, err := s.connection.ExecContext(
		s.ctx,
		"INSERT INTO tags (id, name, created_at, updated_at) VALUES (?, ?, ?, ?)",
		1, "Test Tag", createdAt, createdAt,
	)
	s.NoError(err)

	tag, err := s.tagsRepository.GetTagByID(s.ctx, 1)
	s.NoError(err)
	s.NotNil(tag)
	s.Equal(uint32(1), tag.ID)
	s.Equal("Test Tag", tag.Name)
	s.WithinDuration(createdAt, tag.CreatedAt, time.Second)
	s.WithinDuration(createdAt, tag.UpdatedAt, time.Second)
}

func (s *TagsRepositoryTestSuite) TestGetTagByIDNonExisting() {
	s.traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(context.Background(), mocktracing.NewMockSpan()).
		Times(1)

	tag, err := s.tagsRepository.GetTagByID(s.ctx, 999)
	s.Error(err)
	s.Nil(tag)
}

func (s *TagsRepositoryTestSuite) TestCreateTagsSuccess() {
	s.traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(context.Background(), mocktracing.NewMockSpan()).
		Times(1)

	tagsData := []entities.CreateTagDTO{
		{Name: "Tag 1"},
		{Name: "Tag 2"},
	}

	// Error and zero id due to returning nil ID after insert operation
	// SQLite inner realization without AUTO_INCREMENT for SERIAL PRIMARY KEY
	tagIDs, err := s.tagsRepository.CreateTags(s.ctx, tagsData)
	s.Error(err)
	s.Nil(tagIDs)
}

func (s *TagsRepositoryTestSuite) TestCreateTagsEmpty() {
	s.traceProvider.
		EXPECT().
		Span(gomock.Any(), gomock.Any()).
		Return(context.Background(), mocktracing.NewMockSpan()).
		Times(1)

	s.logger.
		EXPECT().
		ErrorContext(gomock.Any(), gomock.Any(), gomock.Any()).
		Times(1)

	var tagsData []entities.CreateTagDTO

	// Error and zero id due to returning nil ID after insert operation
	// SQLite inner realization without AUTO_INCREMENT for SERIAL PRIMARY KEY
	tagIDs, err := s.tagsRepository.CreateTags(s.ctx, tagsData)
	s.NoError(err)
	s.Empty(tagIDs)

	rows, err := s.connection.QueryContext(s.ctx, "SELECT id FROM tags")
	s.NoError(err)

	defer func() {
		s.NoError(rows.Close())
	}()

	s.False(rows.Next())
}
