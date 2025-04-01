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
	categoriesTableName = "categories"
)

func NewCategoriesRepository(
	dbConnector db.Connector,
	logger logging.Logger,
	traceProvider tracing.Provider,
	spanConfig tracing.SpanConfig,
) *CategoriesRepository {
	return &CategoriesRepository{
		dbConnector:   dbConnector,
		logger:        logger,
		traceProvider: traceProvider,
		spanConfig:    spanConfig,
	}
}

type CategoriesRepository struct {
	dbConnector   db.Connector
	logger        logging.Logger
	traceProvider tracing.Provider
	spanConfig    tracing.SpanConfig
}

func (repo *CategoriesRepository) GetAllCategories(
	ctx context.Context,
) ([]entities.Category, error) {
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
		From(categoriesTableName).
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

	var categories []entities.Category

	for rows.Next() {
		category := entities.Category{}
		columns := db.GetEntityColumns(&category) // Only pointer to use rows.Scan() successfully

		err = rows.Scan(columns...)
		if err != nil {
			return nil, err
		}

		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (repo *CategoriesRepository) GetCategoryByID(
	ctx context.Context,
	id uint32,
) (*entities.Category, error) {
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
		From(categoriesTableName).
		Where(sq.Eq{idColumnName: id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	category := &entities.Category{}

	columns := db.GetEntityColumns(category)
	if err = connection.QueryRowContext(ctx, stmt, params...).Scan(columns...); err != nil {
		return nil, err
	}

	return category, nil
}
