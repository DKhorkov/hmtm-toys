package repositories

import (
	"context"
	"database/sql"
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
		columns = columns[:len(columns)-2]   // Not to paste Tags and Attachments fields to Scan function.
		err = rows.Scan(columns...)
		if err != nil {
			return nil, err
		}

		toys = append(toys, toy)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Reading Tags and Attachments for each Toy in new circle due
	// to next error: https://github.com/lib/pq/issues/635
	// Using toyIndex to avoid range iter semantics error, via using copied variable.
	for toyIndex := range toys {
		if err = repo.processToyTags(ctx, &toys[toyIndex], connection); err != nil {
			return nil, err
		}

		if err = repo.processToyAttachments(ctx, &toys[toyIndex], connection); err != nil {
			return nil, err
		}
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
		columns = columns[:len(columns)-2]   // Not to paste Tags and Attachments fields to Scan function.
		err = rows.Scan(columns...)
		if err != nil {
			return nil, err
		}

		toys = append(toys, toy)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Reading Tags and Attachments for each Toy in new circle due
	// to next error: https://github.com/lib/pq/issues/635
	// Using toyIndex to avoid range iter semantics error, via using copied variable.
	for toyIndex := range toys {
		if err = repo.processToyTags(ctx, &toys[toyIndex], connection); err != nil {
			return nil, err
		}

		if err = repo.processToyAttachments(ctx, &toys[toyIndex], connection); err != nil {
			return nil, err
		}
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
	columns = columns[:len(columns)-2] // Not to paste Tags and Attachments fields to Scan function.
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

	if err = repo.processToyTags(ctx, toy, connection); err != nil {
		return nil, err
	}

	if err = repo.processToyAttachments(ctx, toy, connection); err != nil {
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
				INSERT INTO toys_tags_associations (toy_id, tag_id)
				VALUES 
			`+strings.Join(toysAndTagsInsertPlaceholders, ","),
			toysAndTagsInsertValues...,
		)

		if err != nil {
			return 0, err
		}
	}

	if len(toyData.Attachments) > 0 {
		// Bulk insert of Toy's Attachments.
		toyAttachmentsInsertPlaceholders := make([]string, 0, len(toyData.Attachments))
		toyAttachmentsInsertValues := make([]interface{}, 0, len(toyData.Attachments))
		for index, attachment := range toyData.Attachments {
			toyAttachmentsInsertPlaceholder := fmt.Sprintf("($%d,$%d)",
				index*2+1, // (*2) - where 2 is number of inserted params.
				index*2+2,
			)

			toyAttachmentsInsertPlaceholders = append(
				toyAttachmentsInsertPlaceholders,
				toyAttachmentsInsertPlaceholder,
			)

			toyAttachmentsInsertValues = append(toyAttachmentsInsertValues, toyID, attachment)
		}

		_, err = transaction.Exec(
			`
				INSERT INTO toys_attachments_associations (toy_id, link)
				VALUES 
			`+strings.Join(toyAttachmentsInsertPlaceholders, ","),
			toyAttachmentsInsertValues...,
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

func (repo *CommonToysRepository) processToyTags(
	ctx context.Context,
	toy *entities.Toy,
	connection *sql.Conn,
) error {
	rows, err := connection.QueryContext(
		ctx,
		`
			SELECT * 
			FROM tags AS t
			WHERE t.id IN (
			    SELECT tta.tag_id
				FROM toys_tags_associations AS tta
				WHERE tta.toy_id = $1
			)
		`,
		toy.ID,
	)

	if err != nil {
		return err
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
		var tag entities.Tag
		columns := db.GetEntityColumns(&tag) // Only pointer to use rows.Scan() successfully
		err = rows.Scan(columns...)
		if err != nil {
			return err
		}

		tags = append(tags, tag)
	}

	if err = rows.Err(); err != nil {
		return err
	}

	toy.Tags = tags
	return nil
}

func (repo *CommonToysRepository) processToyAttachments(
	ctx context.Context,
	toy *entities.Toy,
	connection *sql.Conn,
) error {
	rows, err := connection.QueryContext(
		ctx,
		`
			SELECT *
			FROM toys_attachments_associations AS taa
			WHERE taa.toy_id = $1
		`,
		toy.ID,
	)

	if err != nil {
		return err
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

	var attachments []entities.Attachment
	for rows.Next() {
		var attachment entities.Attachment
		columns := db.GetEntityColumns(&attachment) // Only pointer to use rows.Scan() successfully
		err = rows.Scan(columns...)
		if err != nil {
			return err
		}

		attachments = append(attachments, attachment)
	}

	if err = rows.Err(); err != nil {
		return err
	}

	toy.Attachments = attachments
	return nil
}
