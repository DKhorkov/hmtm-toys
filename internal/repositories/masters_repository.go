package repositories

import (
	"context"
	"log/slog"

	"github.com/DKhorkov/libs/db"
	"github.com/DKhorkov/libs/logging"
	"github.com/DKhorkov/libs/tracing"

	"github.com/DKhorkov/hmtm-toys/internal/entities"
)

func NewMastersRepository(
	dbConnector db.Connector,
	logger *slog.Logger,
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

type MastersRepository struct {
	dbConnector   db.Connector
	logger        *slog.Logger
	traceProvider tracing.Provider
	spanConfig    tracing.SpanConfig
}

func (repo *MastersRepository) GetAllMasters(ctx context.Context) ([]entities.Master, error) {
	ctx, span := repo.traceProvider.Span(ctx, tracing.CallerName(tracing.DefaultSkipLevel))
	defer span.End()

	span.AddEvent(repo.spanConfig.Events.Start.Name, repo.spanConfig.Events.Start.Opts...)
	defer span.AddEvent(repo.spanConfig.Events.End.Name, repo.spanConfig.Events.End.Opts...)

	connection, err := repo.dbConnector.Connection(ctx)
	if err != nil {
		return nil, err
	}

	defer db.CloseConnectionContext(ctx, connection, repo.logger)

	rows, err := connection.QueryContext(
		ctx,
		`
			SELECT * 
			FROM masters
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

func (repo *MastersRepository) GetMasterByUserID(ctx context.Context, userID uint64) (*entities.Master, error) {
	ctx, span := repo.traceProvider.Span(ctx, tracing.CallerName(tracing.DefaultSkipLevel))
	defer span.End()

	span.AddEvent(repo.spanConfig.Events.Start.Name, repo.spanConfig.Events.Start.Opts...)
	defer span.AddEvent(repo.spanConfig.Events.End.Name, repo.spanConfig.Events.End.Opts...)

	connection, err := repo.dbConnector.Connection(ctx)
	if err != nil {
		return nil, err
	}

	defer db.CloseConnectionContext(ctx, connection, repo.logger)

	master := &entities.Master{}
	columns := db.GetEntityColumns(master)
	err = connection.QueryRowContext(
		ctx,
		`
			SELECT * 
			FROM masters AS m
			WHERE m.user_id = $1
		`,
		userID,
	).Scan(columns...)

	if err != nil {
		return nil, err
	}

	return master, nil
}

func (repo *MastersRepository) GetMasterByID(ctx context.Context, id uint64) (*entities.Master, error) {
	ctx, span := repo.traceProvider.Span(ctx, tracing.CallerName(tracing.DefaultSkipLevel))
	defer span.End()

	span.AddEvent(repo.spanConfig.Events.Start.Name, repo.spanConfig.Events.Start.Opts...)
	defer span.AddEvent(repo.spanConfig.Events.End.Name, repo.spanConfig.Events.End.Opts...)

	connection, err := repo.dbConnector.Connection(ctx)
	if err != nil {
		return nil, err
	}

	defer db.CloseConnectionContext(ctx, connection, repo.logger)

	master := &entities.Master{}
	columns := db.GetEntityColumns(master)
	err = connection.QueryRowContext(
		ctx,
		`
			SELECT * 
			FROM masters AS m
			WHERE m.id = $1
		`,
		id,
	).Scan(columns...)

	if err != nil {
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

	var masterID uint64
	err = connection.QueryRowContext(
		ctx,
		`
			INSERT INTO masters (user_id, info) 
			VALUES ($1, $2)
			RETURNING masters.id
		`,
		masterData.UserID,
		masterData.Info,
	).Scan(&masterID)

	if err != nil {
		return 0, err
	}

	return masterID, nil
}
