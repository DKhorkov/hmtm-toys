package repositories

import (
	"context"

	"github.com/DKhorkov/libs/db"
	"github.com/DKhorkov/libs/logging"
	"github.com/DKhorkov/libs/tracing"

	sq "github.com/Masterminds/squirrel"

	"github.com/DKhorkov/hmtm-toys/internal/entities"
)

const (
	tagsTableName     = "tags"
	tagNameColumnName = "name"
)

func NewTagsRepository(
	dbConnector db.Connector,
	logger logging.Logger,
	traceProvider tracing.Provider,
	spanConfig tracing.SpanConfig,
) *TagsRepository {
	return &TagsRepository{
		dbConnector:   dbConnector,
		logger:        logger,
		traceProvider: traceProvider,
		spanConfig:    spanConfig,
	}
}

type TagsRepository struct {
	dbConnector   db.Connector
	logger        logging.Logger
	traceProvider tracing.Provider
	spanConfig    tracing.SpanConfig
}

func (repo *TagsRepository) GetAllTags(ctx context.Context) ([]entities.Tag, error) {
	ctx, span := repo.traceProvider.Span(ctx, tracing.CallerName(tracing.DefaultSkipLevel))
	defer span.End()

	span.AddEvent(repo.spanConfig.Events.Start.Name, repo.spanConfig.Events.Start.Opts...)
	defer span.AddEvent(repo.spanConfig.Events.End.Name, repo.spanConfig.Events.End.Opts...)

	connection, err := repo.dbConnector.Connection(ctx)
	if err != nil {
		return nil, err
	}

	defer db.CloseConnectionContext(ctx, connection, repo.logger)

	stmt, params, err := sq.
		Select(selectAllColumns).
		From(tagsTableName).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := connection.QueryContext(
		ctx,
		stmt,
		params...,
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

func (repo *TagsRepository) GetTagByID(ctx context.Context, id uint32) (*entities.Tag, error) {
	ctx, span := repo.traceProvider.Span(ctx, tracing.CallerName(tracing.DefaultSkipLevel))
	defer span.End()

	span.AddEvent(repo.spanConfig.Events.Start.Name, repo.spanConfig.Events.Start.Opts...)
	defer span.AddEvent(repo.spanConfig.Events.End.Name, repo.spanConfig.Events.End.Opts...)

	connection, err := repo.dbConnector.Connection(ctx)
	if err != nil {
		return nil, err
	}

	defer db.CloseConnectionContext(ctx, connection, repo.logger)

	stmt, params, err := sq.
		Select(selectAllColumns).
		From(tagsTableName).
		Where(sq.Eq{idColumnName: id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	tag := &entities.Tag{}

	columns := db.GetEntityColumns(tag)
	if err = connection.QueryRowContext(ctx, stmt, params...).Scan(columns...); err != nil {
		return nil, err
	}

	return tag, nil
}

func (repo *TagsRepository) CreateTags(
	ctx context.Context,
	tagsData []entities.CreateTagDTO,
) ([]uint32, error) {
	ctx, span := repo.traceProvider.Span(ctx, tracing.CallerName(tracing.DefaultSkipLevel))
	defer span.End()

	span.AddEvent(repo.spanConfig.Events.Start.Name, repo.spanConfig.Events.Start.Opts...)
	defer span.AddEvent(repo.spanConfig.Events.End.Name, repo.spanConfig.Events.End.Opts...)

	transaction, err := repo.dbConnector.Transaction(ctx)
	if err != nil {
		return nil, err
	}

	// Rollback transaction according Go best practises https://go.dev/doc/database/execute-transactions.
	defer func() {
		if err = transaction.Rollback(); err != nil {
			logging.LogErrorContext(ctx, repo.logger, "failed to rollback db transaction", err)
		}
	}()

	tagIDs := make([]uint32, len(tagsData))

	for i, tag := range tagsData {
		stmt, params, err := sq.
			Insert(tagsTableName).
			Columns(tagNameColumnName).
			Values(tag.Name).
			Suffix(returningIDSuffix).
			PlaceholderFormat(sq.Dollar). // pq postgres driver works only with $ placeholders
			ToSql()
		if err != nil {
			return nil, err
		}

		var tagID uint32
		if err = transaction.QueryRowContext(ctx, stmt, params...).Scan(&tagID); err != nil {
			return nil, err
		}

		tagIDs[i] = tagID
	}

	err = transaction.Commit()
	if err != nil {
		return nil, err
	}

	return tagIDs, nil
}
