package repositories

import (
	"context"
	"log/slog"

	"github.com/DKhorkov/hmtm-toys/internal/entities"
	"github.com/DKhorkov/libs/db"
	"github.com/DKhorkov/libs/logging"
)

func NewCommonTagsRepository(
	dbConnector db.Connector,
	logger *slog.Logger,
) *CommonTagsRepository {
	return &CommonTagsRepository{
		dbConnector: dbConnector,
		logger:      logger,
	}
}

type CommonTagsRepository struct {
	dbConnector db.Connector
	logger      *slog.Logger
}

func (repo *CommonTagsRepository) GetAllTags(ctx context.Context) ([]entities.Tag, error) {
	connection, err := repo.dbConnector.Connection(ctx)
	if err != nil {
		return nil, err
	}

	defer db.CloseConnectionContext(ctx, connection, repo.logger)

	rows, err := connection.QueryContext(
		ctx,
		`
			SELECT * 
			FROM tags
		`,
	)

	if err != nil {
		return nil, err
	}

	defer func() {
		if err = rows.Close(); err != nil {
			logging.LogErrorContext(
				ctx,
				repo.logger,
				"error during closing SQL rows",
				err,
			)
		}
	}()

	var tags []entities.Tag
	for rows.Next() {
		tag := entities.Tag{}
		columns := db.GetEntityColumns(&tag) // Only pointer to use rows.Scan() successfully
		err = rows.Scan(columns...)
		if err != nil {
			return nil, err
		}

		tags = append(tags, tag)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tags, nil
}

func (repo *CommonTagsRepository) GetTagByID(ctx context.Context, id uint32) (*entities.Tag, error) {
	connection, err := repo.dbConnector.Connection(ctx)
	if err != nil {
		return nil, err
	}

	defer db.CloseConnectionContext(ctx, connection, repo.logger)

	tag := &entities.Tag{}
	columns := db.GetEntityColumns(tag)
	err = connection.QueryRowContext(
		ctx,
		`
			SELECT * 
			FROM tags AS t
			WHERE t.id = $1
		`,
		id,
	).Scan(columns...)

	if err != nil {
		return nil, err
	}

	return tag, nil
}
