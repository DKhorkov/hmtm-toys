package repositories

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/DKhorkov/hmtm-toys/internal/entities"
	"github.com/DKhorkov/libs/db"
	"github.com/DKhorkov/libs/logging"
)

func NewCommonToysRepository(
	dbConnector db.Connector,
	logger *slog.Logger,
) *CommonToysRepository {
	return &CommonToysRepository{
		dbConnector: dbConnector,
		logger:      logger,
	}
}

type CommonToysRepository struct {
	dbConnector db.Connector
	logger      *slog.Logger
}

func (repo *CommonToysRepository) GetAllToys(ctx context.Context) ([]entities.Toy, error) {
	connection, err := repo.dbConnector.Connection(ctx)
	if err != nil {
		return nil, err
	}

	defer db.CloseConnectionContext(ctx, connection, repo.logger)

	rows, err := connection.QueryContext(
		ctx,
		`
			SELECT * 
			FROM toys
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

	var toys []entities.Toy
	for rows.Next() {
		toy := entities.Toy{}
		columns := db.GetEntityColumns(&toy) // Only pointer to use rows.Scan() successfully
		columns = columns[:len(columns)-1]   // not to paste tags field ([]Tag) to Scan function.
		err = rows.Scan(columns...)
		if err != nil {
			return nil, err
		}

		toys = append(toys, toy)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return toys, nil
}

func (repo *CommonToysRepository) GetMasterToys(ctx context.Context, masterID uint64) ([]entities.Toy, error) {
	connection, err := repo.dbConnector.Connection(ctx)
	if err != nil {
		return nil, err
	}

	defer db.CloseConnectionContext(ctx, connection, repo.logger)

	rows, err := connection.QueryContext(
		ctx,
		`
			SELECT * 
			FROM toys AS t
			WHERE t.master_id = $1
		`,
		masterID,
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

	var toys []entities.Toy
	for rows.Next() {
		toy := entities.Toy{}
		columns := db.GetEntityColumns(&toy) // Only pointer to use rows.Scan() successfully
		columns = columns[:len(columns)-1]   // not to paste tags field ([]Tag) to Scan function.
		err = rows.Scan(columns...)
		if err != nil {
			return nil, err
		}

		toys = append(toys, toy)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return toys, nil
}

func (repo *CommonToysRepository) GetToyByID(ctx context.Context, id uint64) (*entities.Toy, error) {
	connection, err := repo.dbConnector.Connection(ctx)
	if err != nil {
		return nil, err
	}

	defer db.CloseConnectionContext(ctx, connection, repo.logger)

	toy := &entities.Toy{}
	columns := db.GetEntityColumns(toy)
	columns = columns[:len(columns)-1] // not to paste tags field ([]Tag) to Scan function.
	err = connection.QueryRowContext(
		ctx,
		`
			SELECT * 
			FROM toys AS t
			WHERE t.id = $1
		`,
		id,
	).Scan(columns...)

	if err != nil {
		return nil, err
	}

	return toy, nil
}

func (repo *CommonToysRepository) AddToy(ctx context.Context, toyData entities.AddToyDTO) (uint64, error) {
	transaction, err := repo.dbConnector.Transaction(ctx)
	if err != nil {
		return 0, err
	}

	// Rollback transaction according Go best practises https://go.dev/doc/database/execute-transactions.
	defer func() {
		if err = transaction.Rollback(); err != nil {
			logging.LogErrorContext(ctx, repo.logger, "failed to rollback db transaction", err)
		}
	}()

	var toyID uint64
	err = transaction.QueryRow(
		`
			INSERT INTO toys (master_id, category_id, name, description, price, quantity) 
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING toys.id
		`,
		toyData.MasterID,
		toyData.CategoryID,
		toyData.Name,
		toyData.Description,
		toyData.Price,
		toyData.Quantity,
	).Scan(&toyID)

	if err != nil {
		return 0, err
	}

	if len(toyData.TagIDs) > 0 {
		// Bulk insert of Toy's Tags.
		toysAndTagsInsertPlaceholders := make([]string, 0, len(toyData.TagIDs))
		toysAndTagsInsertValues := make([]interface{}, 0, len(toyData.TagIDs))
		for index, tagID := range toyData.TagIDs {
			toysAndTagsInsertPlaceholder := fmt.Sprintf("($%d,$%d)",
				index*2+1, // (*2) - where 2 is number of inserted params.
				index*2+2,
			)

			toysAndTagsInsertPlaceholders = append(toysAndTagsInsertPlaceholders, toysAndTagsInsertPlaceholder)
			toysAndTagsInsertValues = append(toysAndTagsInsertValues, toyID, tagID)
		}

		_, err = transaction.Exec(
			`
			INSERT INTO toys_and_tags_associations (toy_id, tag_id)
			VALUES 
		`+strings.Join(toysAndTagsInsertPlaceholders, ","),
			toysAndTagsInsertValues...,
		)

		if err != nil {
			return 0, err
		}
	}

	err = transaction.Commit()
	if err != nil {
		return 0, err
	}

	return toyID, nil
}
