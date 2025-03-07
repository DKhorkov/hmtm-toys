package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"

	"github.com/DKhorkov/libs/db"
	"github.com/DKhorkov/libs/logging"
	"github.com/DKhorkov/libs/tracing"

	"github.com/DKhorkov/hmtm-toys/internal/entities"
)

const (
	toysTableName                   = "toys"
	toysAndTagsAssociationTableName = "toys_tags_associations"
	toysAttachmentsTableName        = "toys_attachments"
	idColumnName                    = "id"
	toyCategoryIDColumnName         = "category_id"
	toyNameColumnName               = "name"
	toyDescriptionColumnName        = "description"
	toyPriceColumnName              = "price"
	toyQuantityColumnName           = "quantity"
	toyIDColumnName                 = "toy_id"
	tagIDColumnName                 = "tag_id"
	attachmentLinkColumnName        = "link"
)

func NewToysRepository(
	dbConnector db.Connector,
	logger logging.Logger,
	traceProvider tracing.Provider,
	spanConfig tracing.SpanConfig,
) *ToysRepository {
	return &ToysRepository{
		dbConnector:   dbConnector,
		logger:        logger,
		traceProvider: traceProvider,
		spanConfig:    spanConfig,
	}
}

type ToysRepository struct {
	dbConnector   db.Connector
	logger        logging.Logger
	traceProvider tracing.Provider
	spanConfig    tracing.SpanConfig
}

func (repo *ToysRepository) GetAllToys(ctx context.Context) ([]entities.Toy, error) {
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
	// Using toy index to avoid range iter semantics error, via using copied variable.
	for i, toy := range toys {
		tags, err := repo.getToyTags(ctx, toy.ID, connection)
		if err != nil {
			return nil, err
		}

		toys[i].Tags = tags

		attachments, err := repo.getToyAttachments(ctx, toy.ID, connection)
		if err != nil {
			return nil, err
		}

		toys[i].Attachments = attachments
	}

	return toys, nil
}

func (repo *ToysRepository) GetMasterToys(ctx context.Context, masterID uint64) ([]entities.Toy, error) {
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
	// Using toy index to avoid range iter semantics error, via using copied variable.
	for i, toy := range toys {
		tags, err := repo.getToyTags(ctx, toy.ID, connection)
		if err != nil {
			return nil, err
		}

		toys[i].Tags = tags

		attachments, err := repo.getToyAttachments(ctx, toy.ID, connection)
		if err != nil {
			return nil, err
		}

		toys[i].Attachments = attachments
	}

	return toys, nil
}

func (repo *ToysRepository) GetToyByID(ctx context.Context, id uint64) (*entities.Toy, error) {
	ctx, span := repo.traceProvider.Span(ctx, tracing.CallerName(tracing.DefaultSkipLevel))
	defer span.End()

	span.AddEvent(repo.spanConfig.Events.Start.Name, repo.spanConfig.Events.Start.Opts...)
	defer span.AddEvent(repo.spanConfig.Events.End.Name, repo.spanConfig.Events.End.Opts...)

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

	tags, err := repo.getToyTags(ctx, toy.ID, connection)
	if err != nil {
		return nil, err
	}

	toy.Tags = tags

	attachments, err := repo.getToyAttachments(ctx, toy.ID, connection)
	if err != nil {
		return nil, err
	}

	toy.Attachments = attachments
	return toy, nil
}

func (repo *ToysRepository) AddToy(ctx context.Context, toyData entities.AddToyDTO) (uint64, error) {
	ctx, span := repo.traceProvider.Span(ctx, tracing.CallerName(tracing.DefaultSkipLevel))
	defer span.End()

	span.AddEvent(repo.spanConfig.Events.Start.Name, repo.spanConfig.Events.Start.Opts...)
	defer span.AddEvent(repo.spanConfig.Events.End.Name, repo.spanConfig.Events.End.Opts...)

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
	err = transaction.QueryRowContext(
		ctx,
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
				INSERT INTO toys_attachments (toy_id, link)
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

func (repo *ToysRepository) getToyTags(
	ctx context.Context,
	toyID uint64,
	connection *sql.Conn,
) ([]entities.Tag, error) {
	ctx, span := repo.traceProvider.Span(ctx, tracing.CallerName(tracing.DefaultSkipLevel))
	defer span.End()

	span.AddEvent(repo.spanConfig.Events.Start.Name, repo.spanConfig.Events.Start.Opts...)
	defer span.AddEvent(repo.spanConfig.Events.End.Name, repo.spanConfig.Events.End.Opts...)

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
		toyID,
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
		var tag entities.Tag
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

func (repo *ToysRepository) getToyAttachments(
	ctx context.Context,
	toyID uint64,
	connection *sql.Conn,
) ([]entities.Attachment, error) {
	ctx, span := repo.traceProvider.Span(ctx, tracing.CallerName(tracing.DefaultSkipLevel))
	defer span.End()

	span.AddEvent(repo.spanConfig.Events.Start.Name, repo.spanConfig.Events.Start.Opts...)
	defer span.AddEvent(repo.spanConfig.Events.End.Name, repo.spanConfig.Events.End.Opts...)

	rows, err := connection.QueryContext(
		ctx,
		`
			SELECT *
			FROM toys_attachments AS ta
			WHERE ta.toy_id = $1
		`,
		toyID,
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

	var attachments []entities.Attachment
	for rows.Next() {
		var attachment entities.Attachment
		columns := db.GetEntityColumns(&attachment) // Only pointer to use rows.Scan() successfully
		err = rows.Scan(columns...)
		if err != nil {
			return nil, err
		}

		attachments = append(attachments, attachment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return attachments, nil
}

func (repo *ToysRepository) DeleteToy(ctx context.Context, id uint64) error {
	ctx, span := repo.traceProvider.Span(ctx, tracing.CallerName(tracing.DefaultSkipLevel))
	defer span.End()

	span.AddEvent(repo.spanConfig.Events.Start.Name, repo.spanConfig.Events.Start.Opts...)
	defer span.AddEvent(repo.spanConfig.Events.End.Name, repo.spanConfig.Events.End.Opts...)

	connection, err := repo.dbConnector.Connection(ctx)
	if err != nil {
		return err
	}

	defer db.CloseConnectionContext(ctx, connection, repo.logger)

	_, err = connection.ExecContext(
		ctx,
		`
			DELETE 
			FROM toys AS t
			WHERE t.id = $1
		`,
		id,
	)

	return err
}

func (repo *ToysRepository) UpdateToy(ctx context.Context, toyData entities.UpdateToyDTO) error {
	ctx, span := repo.traceProvider.Span(ctx, tracing.CallerName(tracing.DefaultSkipLevel))
	defer span.End()

	span.AddEvent(repo.spanConfig.Events.Start.Name, repo.spanConfig.Events.Start.Opts...)
	defer span.AddEvent(repo.spanConfig.Events.End.Name, repo.spanConfig.Events.End.Opts...)

	transaction, err := repo.dbConnector.Transaction(ctx)
	if err != nil {
		return err
	}

	// Rollback transaction according Go best practises https://go.dev/doc/database/execute-transactions.
	defer func() {
		if err = transaction.Rollback(); err != nil {
			logging.LogErrorContext(ctx, repo.logger, "failed to rollback db transaction", err)
		}
	}()

	stmt, params, err := sq.
		Update(toysTableName).
		Where(sq.Eq{idColumnName: toyData.ID}).
		Set(toyCategoryIDColumnName, toyData.CategoryID).
		Set(toyNameColumnName, toyData.Name).
		Set(toyDescriptionColumnName, toyData.Description).
		Set(toyPriceColumnName, toyData.Price).
		Set(toyQuantityColumnName, toyData.Quantity).
		PlaceholderFormat(sq.Dollar). // pq postgres driver works only with $ placeholders
		ToSql()

	if err != nil {
		return err
	}

	if _, err = transaction.ExecContext(ctx, stmt, params...); err != nil {
		return err
	}

	if len(toyData.TagIDsToAdd) > 0 {
		builder := sq.Insert(toysAndTagsAssociationTableName).Columns(toyIDColumnName, tagIDColumnName)
		for _, tagID := range toyData.TagIDsToAdd {
			builder = builder.Values(toyData.ID, tagID)
		}

		if stmt, params, err = builder.PlaceholderFormat(sq.Dollar).ToSql(); err != nil {
			return err
		}

		if _, err = transaction.ExecContext(ctx, stmt, params...); err != nil {
			return err
		}
	}

	if len(toyData.TagIDsToDelete) > 0 {
		stmt, params, err = sq.
			Delete(toysAndTagsAssociationTableName).
			Where(
				sq.And{
					sq.Eq{toyIDColumnName: toyData.ID},
					sq.Eq{tagIDColumnName: toyData.TagIDsToDelete},
				},
			).
			PlaceholderFormat(sq.Dollar).
			ToSql()

		if err != nil {
			return err
		}

		if _, err = transaction.ExecContext(ctx, stmt, params...); err != nil {
			return err
		}
	}

	if len(toyData.AttachmentsToAdd) > 0 {
		builder := sq.Insert(toysAttachmentsTableName).Columns(toyIDColumnName, attachmentLinkColumnName)
		for _, attachment := range toyData.AttachmentsToAdd {
			builder = builder.Values(toyData.ID, attachment)
		}

		if stmt, params, err = builder.PlaceholderFormat(sq.Dollar).ToSql(); err != nil {
			return err
		}

		if _, err = transaction.ExecContext(ctx, stmt, params...); err != nil {
			return err
		}
	}

	if len(toyData.AttachmentIDsToDelete) > 0 {
		stmt, params, err = sq.
			Delete(toysAttachmentsTableName).
			Where(sq.Eq{idColumnName: toyData.AttachmentIDsToDelete}).
			PlaceholderFormat(sq.Dollar).
			ToSql()

		if err != nil {
			return err
		}

		if _, err = transaction.ExecContext(ctx, stmt, params...); err != nil {
			return err
		}
	}

	return transaction.Commit()
}
