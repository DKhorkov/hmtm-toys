package repositories

import (
	"context"
	"log/slog"

	"github.com/DKhorkov/libs/tracing"

	"github.com/DKhorkov/hmtm-toys/internal/entities"
	"github.com/DKhorkov/libs/db"
	"github.com/DKhorkov/libs/logging"
)

func NewCommonCategoriesRepository(
	dbConnector db.Connector,
	logger *slog.Logger,
	traceProvider tracing.TraceProvider,
	spanConfig tracing.SpanConfig,
) *CommonCategoriesRepository {
	return &CommonCategoriesRepository{
		dbConnector:   dbConnector,
		logger:        logger,
		traceProvider: traceProvider,
		spanConfig:    spanConfig,
	}
}

type CommonCategoriesRepository struct {
	dbConnector   db.Connector
	logger        *slog.Logger
	traceProvider tracing.TraceProvider
	spanConfig    tracing.SpanConfig
}

func (repo *CommonCategoriesRepository) GetAllCategories(ctx context.Context) ([]entities.Category, error) {
	ctx, span := repo.traceProvider.Span(ctx, tracing.CallerName(tracing.DefaultSkipLevel))
	defer span.End()

	span.AddEvent(repo.spanConfig.Events.Start.Name, repo.spanConfig.Events.Start.Opts...)

	connection, err := repo.dbConnector.Connection(ctx)
	if err != nil {
		return nil, err
	}

	defer db.CloseConnectionContext(ctx, connection, repo.logger)

	rows, err := connection.QueryContext(
		ctx,
		`
			SELECT * 
			FROM categories
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

	span.AddEvent(repo.spanConfig.Events.End.Name, repo.spanConfig.Events.End.Opts...)
	return categories, nil
}

func (repo *CommonCategoriesRepository) GetCategoryByID(ctx context.Context, id uint32) (*entities.Category, error) {
	ctx, span := repo.traceProvider.Span(ctx, tracing.CallerName(tracing.DefaultSkipLevel))
	defer span.End()

	span.AddEvent(repo.spanConfig.Events.Start.Name, repo.spanConfig.Events.Start.Opts...)

	connection, err := repo.dbConnector.Connection(ctx)
	if err != nil {
		return nil, err
	}

	defer db.CloseConnectionContext(ctx, connection, repo.logger)

	category := &entities.Category{}
	columns := db.GetEntityColumns(category)
	err = connection.QueryRowContext(
		ctx,
		`
			SELECT * 
			FROM categories AS c
			WHERE c.id = $1
		`,
		id,
	).Scan(columns...)

	if err != nil {
		return nil, err
	}

	span.AddEvent(repo.spanConfig.Events.End.Name, repo.spanConfig.Events.End.Opts...)
	return category, nil
}
