package repositories

import (
	"context"
	"fmt"

	"github.com/DKhorkov/libs/db"
	"github.com/DKhorkov/libs/logging"
	"github.com/DKhorkov/libs/tracing"

	sq "github.com/Masterminds/squirrel"

	"github.com/DKhorkov/hmtm-toys/internal/entities"
)

const (
	mastersTableName     = "masters"
	userIDColumnName     = "user_id"
	masterInfoColumnName = "info"
)

type MastersRepository struct {
	dbConnector   db.Connector
	logger        logging.Logger
	traceProvider tracing.Provider
	spanConfig    tracing.SpanConfig
}

func NewMastersRepository(
	dbConnector db.Connector,
	logger logging.Logger,
	traceProvider tracing.Provider,
	spanConfig tracing.SpanConfig,
) *MastersRepository {
	return &MastersRepository{
		dbConnector:   dbConnector,
		logger:        logger,
		traceProvider: traceProvider,
		spanConfig:    spanConfig,
	}
}

func (repo *MastersRepository) GetMasters(
	ctx context.Context,
	pagination *entities.Pagination,
) ([]entities.Master, error) {
	ctx, span := repo.traceProvider.Span(ctx, tracing.CallerName(tracing.DefaultSkipLevel))
	defer span.End()

	span.AddEvent(repo.spanConfig.Events.Start.Name, repo.spanConfig.Events.Start.Opts...)
	defer span.AddEvent(repo.spanConfig.Events.End.Name, repo.spanConfig.Events.End.Opts...)

	connection, err := repo.dbConnector.Connection(ctx)
	if err != nil {
		return nil, err
	}

	defer db.CloseConnectionContext(ctx, connection, repo.logger)

	builder := sq.
		Select(selectAllColumns).
		From(mastersTableName).
		OrderBy(fmt.Sprintf("%s %s", idColumnName, DESC)).
		PlaceholderFormat(sq.Dollar)

	if pagination != nil && pagination.Limit != nil {
		builder = builder.Limit(*pagination.Limit)
	}

	if pagination != nil && pagination.Offset != nil {
		builder = builder.Offset(*pagination.Offset)
	}

	stmt, params, err := builder.ToSql()
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

	var masters []entities.Master

	for rows.Next() {
		master := entities.Master{}
		columns := db.GetEntityColumns(&master) // Only pointer to use rows.Scan() successfully

		err = rows.Scan(columns...)
		if err != nil {
			return nil, err
		}

		masters = append(masters, master)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return masters, nil
}

func (repo *MastersRepository) GetMasterByUserID(
	ctx context.Context,
	userID uint64,
) (*entities.Master, error) {
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
		From(mastersTableName).
		Where(sq.Eq{userIDColumnName: userID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	master := &entities.Master{}

	columns := db.GetEntityColumns(master)
	if err = connection.QueryRowContext(ctx, stmt, params...).Scan(columns...); err != nil {
		return nil, err
	}

	return master, nil
}

func (repo *MastersRepository) GetMasterByID(
	ctx context.Context,
	id uint64,
) (*entities.Master, error) {
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
		From(mastersTableName).
		Where(sq.Eq{idColumnName: id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	master := &entities.Master{}

	columns := db.GetEntityColumns(master)
	if err = connection.QueryRowContext(ctx, stmt, params...).Scan(columns...); err != nil {
		return nil, err
	}

	return master, nil
}

func (repo *MastersRepository) RegisterMaster(
	ctx context.Context,
	masterData entities.RegisterMasterDTO,
) (uint64, error) {
	ctx, span := repo.traceProvider.Span(ctx, tracing.CallerName(tracing.DefaultSkipLevel))
	defer span.End()

	span.AddEvent(repo.spanConfig.Events.Start.Name, repo.spanConfig.Events.Start.Opts...)
	defer span.AddEvent(repo.spanConfig.Events.End.Name, repo.spanConfig.Events.End.Opts...)

	connection, err := repo.dbConnector.Connection(ctx)
	if err != nil {
		return 0, err
	}

	defer db.CloseConnectionContext(ctx, connection, repo.logger)

	stmt, params, err := sq.
		Insert(mastersTableName).
		Columns(
			userIDColumnName,
			masterInfoColumnName,
		).
		Values(
			masterData.UserID,
			masterData.Info,
		).
		Suffix(returningIDSuffix).
		PlaceholderFormat(sq.Dollar). // pq postgres driver works only with $ placeholders
		ToSql()
	if err != nil {
		return 0, err
	}

	var masterID uint64
	if err = connection.QueryRowContext(ctx, stmt, params...).Scan(&masterID); err != nil {
		return 0, err
	}

	return masterID, nil
}

func (repo *MastersRepository) UpdateMaster(
	ctx context.Context,
	masterData entities.UpdateMasterDTO,
) error {
	ctx, span := repo.traceProvider.Span(ctx, tracing.CallerName(tracing.DefaultSkipLevel))
	defer span.End()

	span.AddEvent(repo.spanConfig.Events.Start.Name, repo.spanConfig.Events.Start.Opts...)
	defer span.AddEvent(repo.spanConfig.Events.End.Name, repo.spanConfig.Events.End.Opts...)

	connection, err := repo.dbConnector.Connection(ctx)
	if err != nil {
		return err
	}

	defer db.CloseConnectionContext(ctx, connection, repo.logger)

	stmt, params, err := sq.
		Update(mastersTableName).
		Where(sq.Eq{idColumnName: masterData.ID}).
		Set(masterInfoColumnName, masterData.Info).
		PlaceholderFormat(sq.Dollar). // pq postgres driver works only with $ placeholders
		ToSql()
	if err != nil {
		return err
	}

	_, err = connection.ExecContext(
		ctx,
		stmt,
		params...,
	)

	return err
}
