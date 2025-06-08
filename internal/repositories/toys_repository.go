package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/DKhorkov/libs/db"
	"github.com/DKhorkov/libs/logging"
	"github.com/DKhorkov/libs/tracing"

	sq "github.com/Masterminds/squirrel"

	"github.com/DKhorkov/hmtm-toys/internal/entities"
)

const (
	selectAllColumns                = "*"
	selectCount                     = "COUNT(*)"
	toysTableName                   = "toys"
	toysAndTagsAssociationTableName = "toys_tags_associations"
	toysAttachmentsTableName        = "toys_attachments"
	idColumnName                    = "id"
	categoryIDColumnName            = "category_id"
	toyNameColumnName               = "name"
	toyDescriptionColumnName        = "description"
	toyPriceColumnName              = "price"
	toyQuantityColumnName           = "quantity"
	toyIDColumnName                 = "toy_id"
	tagIDColumnName                 = "tag_id"
	masterIDColumnName              = "master_id"
	attachmentLinkColumnName        = "link"
	returningIDSuffix               = "RETURNING id"
	createdAtColumnName             = "created_at"
	updatedAtColumnName             = "updated_at"
	desc                            = "DESC"
	asc                             = "ASC"
)

type ToysRepository struct {
	dbConnector   db.Connector
	logger        logging.Logger
	traceProvider tracing.Provider
	spanConfig    tracing.SpanConfig
}

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

func (repo *ToysRepository) GetToys(
	ctx context.Context,
	pagination *entities.Pagination,
	filters *entities.ToysFilters,
) ([]entities.Toy, error) {
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
		From(toysTableName).
		PlaceholderFormat(sq.Dollar)

	if filters != nil && filters.Search != nil && *filters.Search != "" {
		searchTerm := "%" + strings.ToLower(*filters.Search) + "%"
		builder = builder.
			Where(
				sq.Or{
					sq.ILike{
						fmt.Sprintf(
							"%s.%s",
							toysTableName,
							toyNameColumnName,
						): searchTerm,
					},
					sq.ILike{
						fmt.Sprintf(
							"%s.%s",
							toysTableName,
							toyDescriptionColumnName,
						): searchTerm,
					},
				},
			)
	}

	if filters != nil && (filters.PriceFloor != nil || filters.PriceCeil != nil) {
		priceConditions := sq.And{}
		if filters.PriceFloor != nil {
			priceConditions = append(
				priceConditions,
				sq.GtOrEq{
					fmt.Sprintf(
						"%s.%s",
						toysTableName,
						toyPriceColumnName,
					): *filters.PriceFloor,
				},
			)
		}

		if filters.PriceCeil != nil {
			priceConditions = append(
				priceConditions,
				sq.LtOrEq{
					fmt.Sprintf(
						"%s.%s",
						toysTableName,
						toyPriceColumnName,
					): *filters.PriceCeil,
				},
			)
		}

		builder = builder.Where(priceConditions)
	}

	if filters != nil && filters.QuantityFloor != nil {
		builder = builder.
			Where(
				sq.GtOrEq{
					fmt.Sprintf(
						"%s.%s",
						toysTableName,
						toyQuantityColumnName,
					): *filters.QuantityFloor,
				},
			)
	}

	if filters != nil && filters.CategoryIDs != nil {
		builder = builder.
			Where(
				sq.Eq{
					fmt.Sprintf(
						"%s.%s",
						toysTableName,
						categoryIDColumnName,
					): filters.CategoryIDs,
				},
			)
	}

	if filters != nil && len(filters.TagIDs) > 0 {
		for _, tagID := range filters.TagIDs {
			builder = builder.
				Where(
					sq.Expr(
						fmt.Sprintf(
							"EXISTS (SELECT 1 FROM %s WHERE %s.%s = %s.%s AND %s.%s = ?)",
							toysAndTagsAssociationTableName,
							toysAndTagsAssociationTableName,
							toyIDColumnName,
							toysTableName,
							idColumnName,
							toysAndTagsAssociationTableName,
							tagIDColumnName,
						),
						tagID,
					),
				)
		}
	}

	createdAtOrder := desc
	if filters != nil && filters.CreatedAtOrderByAsc != nil && *filters.CreatedAtOrderByAsc {
		createdAtOrder = asc
	}

	builder = builder.
		OrderBy(
			fmt.Sprintf(
				"%s.%s %s",
				toysTableName,
				createdAtColumnName,
				createdAtOrder,
			),
		)

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

func (repo *ToysRepository) CountToys(ctx context.Context, filters *entities.ToysFilters) (uint64, error) {
	ctx, span := repo.traceProvider.Span(ctx, tracing.CallerName(tracing.DefaultSkipLevel))
	defer span.End()

	span.AddEvent(repo.spanConfig.Events.Start.Name, repo.spanConfig.Events.Start.Opts...)
	defer span.AddEvent(repo.spanConfig.Events.End.Name, repo.spanConfig.Events.End.Opts...)

	connection, err := repo.dbConnector.Connection(ctx)
	if err != nil {
		return 0, err
	}

	defer db.CloseConnectionContext(ctx, connection, repo.logger)

	builder := sq.
		Select(selectCount).
		From(toysTableName).
		PlaceholderFormat(sq.Dollar)

	if filters != nil && filters.Search != nil && *filters.Search != "" {
		searchTerm := "%" + strings.ToLower(*filters.Search) + "%"
		builder = builder.
			Where(
				sq.Or{
					sq.ILike{
						fmt.Sprintf(
							"%s.%s",
							toysTableName,
							toyNameColumnName,
						): searchTerm,
					},
					sq.ILike{
						fmt.Sprintf(
							"%s.%s",
							toysTableName,
							toyDescriptionColumnName,
						): searchTerm,
					},
				},
			)
	}

	if filters != nil && (filters.PriceFloor != nil || filters.PriceCeil != nil) {
		priceConditions := sq.And{}
		if filters.PriceFloor != nil {
			priceConditions = append(
				priceConditions,
				sq.GtOrEq{
					fmt.Sprintf(
						"%s.%s",
						toysTableName,
						toyPriceColumnName,
					): *filters.PriceFloor,
				},
			)
		}

		if filters.PriceCeil != nil {
			priceConditions = append(
				priceConditions,
				sq.LtOrEq{
					fmt.Sprintf(
						"%s.%s",
						toysTableName,
						toyPriceColumnName,
					): *filters.PriceCeil,
				},
			)
		}

		builder = builder.Where(priceConditions)
	}

	if filters != nil && filters.QuantityFloor != nil {
		builder = builder.
			Where(
				sq.GtOrEq{
					fmt.Sprintf(
						"%s.%s",
						toysTableName,
						toyQuantityColumnName,
					): *filters.QuantityFloor,
				},
			)
	}

	if filters != nil && filters.CategoryIDs != nil {
		builder = builder.
			Where(
				sq.Eq{
					fmt.Sprintf(
						"%s.%s",
						toysTableName,
						categoryIDColumnName,
					): filters.CategoryIDs,
				},
			)
	}

	if filters != nil && len(filters.TagIDs) > 0 {
		for _, tagID := range filters.TagIDs {
			builder = builder.
				Where(
					sq.Expr(
						fmt.Sprintf(
							"EXISTS (SELECT 1 FROM %s WHERE %s.%s = %s.%s AND %s.%s = ?)",
							toysAndTagsAssociationTableName,
							toysAndTagsAssociationTableName,
							toyIDColumnName,
							toysTableName,
							idColumnName,
							toysAndTagsAssociationTableName,
							tagIDColumnName,
						),
						tagID,
					),
				)
		}
	}

	// Для запросов COUNT сортировка не нужна, поэтому параметр CreatedAtOrderByAsc не используется
	stmt, params, err := builder.ToSql()
	if err != nil {
		return 0, err
	}

	var count uint64
	if err = connection.QueryRowContext(ctx, stmt, params...).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (repo *ToysRepository) GetMasterToys(
	ctx context.Context,
	masterID uint64,
	pagination *entities.Pagination,
	filters *entities.ToysFilters,
) ([]entities.Toy, error) {
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
		From(toysTableName).
		Where(sq.Eq{masterIDColumnName: masterID}).
		PlaceholderFormat(sq.Dollar)

	if filters != nil && filters.Search != nil && *filters.Search != "" {
		searchTerm := "%" + strings.ToLower(*filters.Search) + "%"
		builder = builder.
			Where(
				sq.Or{
					sq.ILike{
						fmt.Sprintf(
							"%s.%s",
							toysTableName,
							toyNameColumnName,
						): searchTerm,
					},
					sq.ILike{
						fmt.Sprintf(
							"%s.%s",
							toysTableName,
							toyDescriptionColumnName,
						): searchTerm,
					},
				},
			)
	}

	if filters != nil && (filters.PriceFloor != nil || filters.PriceCeil != nil) {
		priceConditions := sq.And{}
		if filters.PriceFloor != nil {
			priceConditions = append(
				priceConditions,
				sq.GtOrEq{
					fmt.Sprintf(
						"%s.%s",
						toysTableName,
						toyPriceColumnName,
					): *filters.PriceFloor,
				},
			)
		}

		if filters.PriceCeil != nil {
			priceConditions = append(
				priceConditions,
				sq.LtOrEq{
					fmt.Sprintf(
						"%s.%s",
						toysTableName,
						toyPriceColumnName,
					): *filters.PriceCeil,
				},
			)
		}

		builder = builder.Where(priceConditions)
	}

	if filters != nil && filters.QuantityFloor != nil {
		builder = builder.
			Where(
				sq.GtOrEq{
					fmt.Sprintf(
						"%s.%s",
						toysTableName,
						toyQuantityColumnName,
					): *filters.QuantityFloor,
				},
			)
	}

	if filters != nil && filters.CategoryIDs != nil {
		builder = builder.
			Where(
				sq.Eq{
					fmt.Sprintf(
						"%s.%s",
						toysTableName,
						categoryIDColumnName,
					): filters.CategoryIDs,
				},
			)
	}

	if filters != nil && len(filters.TagIDs) > 0 {
		for _, tagID := range filters.TagIDs {
			builder = builder.
				Where(
					sq.Expr(
						fmt.Sprintf(
							"EXISTS (SELECT 1 FROM %s WHERE %s.%s = %s.%s AND %s.%s = ?)",
							toysAndTagsAssociationTableName,
							toysAndTagsAssociationTableName,
							toyIDColumnName,
							toysTableName,
							idColumnName,
							toysAndTagsAssociationTableName,
							tagIDColumnName,
						),
						tagID,
					),
				)
		}
	}

	createdAtOrder := desc
	if filters != nil && filters.CreatedAtOrderByAsc != nil && *filters.CreatedAtOrderByAsc {
		createdAtOrder = asc
	}

	builder = builder.
		OrderBy(
			fmt.Sprintf(
				"%s.%s %s",
				toysTableName,
				createdAtColumnName,
				createdAtOrder,
			),
		)

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

func (repo *ToysRepository) CountMasterToys(
	ctx context.Context,
	masterID uint64,
	filters *entities.ToysFilters,
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

	builder := sq.
		Select(selectCount).
		From(toysTableName).
		Where(sq.Eq{masterIDColumnName: masterID}).
		PlaceholderFormat(sq.Dollar)

	if filters != nil && filters.Search != nil && *filters.Search != "" {
		searchTerm := "%" + strings.ToLower(*filters.Search) + "%"
		builder = builder.
			Where(
				sq.Or{
					sq.ILike{
						fmt.Sprintf(
							"%s.%s",
							toysTableName,
							toyNameColumnName,
						): searchTerm,
					},
					sq.ILike{
						fmt.Sprintf(
							"%s.%s",
							toysTableName,
							toyDescriptionColumnName,
						): searchTerm,
					},
				},
			)
	}

	if filters != nil && (filters.PriceFloor != nil || filters.PriceCeil != nil) {
		priceConditions := sq.And{}
		if filters.PriceFloor != nil {
			priceConditions = append(
				priceConditions,
				sq.GtOrEq{
					fmt.Sprintf(
						"%s.%s",
						toysTableName,
						toyPriceColumnName,
					): *filters.PriceFloor,
				},
			)
		}

		if filters.PriceCeil != nil {
			priceConditions = append(
				priceConditions,
				sq.LtOrEq{
					fmt.Sprintf(
						"%s.%s",
						toysTableName,
						toyPriceColumnName,
					): *filters.PriceCeil,
				},
			)
		}

		builder = builder.Where(priceConditions)
	}

	if filters != nil && filters.QuantityFloor != nil {
		builder = builder.
			Where(
				sq.GtOrEq{
					fmt.Sprintf(
						"%s.%s",
						toysTableName,
						toyQuantityColumnName,
					): *filters.QuantityFloor,
				},
			)
	}

	if filters != nil && filters.CategoryIDs != nil {
		builder = builder.
			Where(
				sq.Eq{
					fmt.Sprintf(
						"%s.%s",
						toysTableName,
						categoryIDColumnName,
					): filters.CategoryIDs,
				},
			)
	}

	if filters != nil && len(filters.TagIDs) > 0 {
		for _, tagID := range filters.TagIDs {
			builder = builder.
				Where(
					sq.Expr(
						fmt.Sprintf(
							"EXISTS (SELECT 1 FROM %s WHERE %s.%s = %s.%s AND %s.%s = ?)",
							toysAndTagsAssociationTableName,
							toysAndTagsAssociationTableName,
							toyIDColumnName,
							toysTableName,
							idColumnName,
							toysAndTagsAssociationTableName,
							tagIDColumnName,
						),
						tagID,
					),
				)
		}
	}

	// Для запросов COUNT сортировка не нужна, поэтому параметр CreatedAtOrderByAsc не используется
	stmt, params, err := builder.ToSql()
	if err != nil {
		return 0, err
	}

	var count uint64
	if err = connection.QueryRowContext(ctx, stmt, params...).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
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

	stmt, params, err := sq.
		Select(selectAllColumns).
		From(toysTableName).
		Where(sq.Eq{idColumnName: id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	toy := &entities.Toy{}
	columns := db.GetEntityColumns(toy)
	columns = columns[:len(columns)-2] // Not to paste Tags and Attachments fields to Scan function.

	if err = connection.QueryRowContext(ctx, stmt, params...).Scan(columns...); err != nil {
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

func (repo *ToysRepository) AddToy(
	ctx context.Context,
	toyData entities.AddToyDTO,
) (uint64, error) {
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

	stmt, params, err := sq.
		Insert(toysTableName).
		Columns(
			masterIDColumnName,
			categoryIDColumnName,
			toyNameColumnName,
			toyDescriptionColumnName,
			toyPriceColumnName,
			toyQuantityColumnName,
		).
		Values(
			toyData.MasterID,
			toyData.CategoryID,
			toyData.Name,
			toyData.Description,
			toyData.Price,
			toyData.Quantity,
		).
		Suffix(returningIDSuffix).
		PlaceholderFormat(sq.Dollar). // pq postgres driver works only with $ placeholders
		ToSql()
	if err != nil {
		return 0, err
	}

	var toyID uint64
	if err = transaction.QueryRowContext(ctx, stmt, params...).Scan(&toyID); err != nil {
		return 0, err
	}

	if err != nil {
		return 0, err
	}

	if len(toyData.TagIDs) > 0 {
		builder := sq.Insert(toysAndTagsAssociationTableName).
			Columns(toyIDColumnName, tagIDColumnName)
		for _, tagID := range toyData.TagIDs {
			builder = builder.Values(toyID, tagID)
		}

		if stmt, params, err = builder.PlaceholderFormat(sq.Dollar).ToSql(); err != nil {
			return 0, err
		}

		if _, err = transaction.ExecContext(ctx, stmt, params...); err != nil {
			return 0, err
		}
	}

	if len(toyData.Attachments) > 0 {
		builder := sq.Insert(toysAttachmentsTableName).
			Columns(toyIDColumnName, attachmentLinkColumnName)
		for _, attachment := range toyData.Attachments {
			builder = builder.Values(toyID, attachment)
		}

		if stmt, params, err = builder.PlaceholderFormat(sq.Dollar).ToSql(); err != nil {
			return 0, err
		}

		if _, err = transaction.ExecContext(ctx, stmt, params...); err != nil {
			return 0, err
		}
	}

	err = transaction.Commit()
	if err != nil {
		return 0, err
	}

	return toyID, nil
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

	stmt, params, err := sq.
		Delete(toysTableName).
		Where(sq.Eq{idColumnName: id}).
		PlaceholderFormat(sq.Dollar).
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

	builder := sq.
		Update(toysTableName).
		Where(sq.Eq{idColumnName: toyData.ID}).
		PlaceholderFormat(sq.Dollar) // pq postgres driver works only with $ placeholders

	if toyData.CategoryID != nil {
		builder = builder.Set(categoryIDColumnName, toyData.CategoryID)
	}

	if toyData.Name != nil {
		builder = builder.Set(toyNameColumnName, toyData.Name)
	}

	if toyData.Description != nil {
		builder = builder.Set(toyDescriptionColumnName, toyData.Description)
	}

	if toyData.Price != nil {
		builder = builder.Set(toyPriceColumnName, toyData.Price)
	}

	if toyData.Quantity != nil {
		builder = builder.Set(toyQuantityColumnName, toyData.Quantity)
	}

	stmt, params, err := builder.ToSql()
	if err != nil {
		return err
	}

	if _, err = transaction.ExecContext(ctx, stmt, params...); err != nil {
		return err
	}

	if len(toyData.TagIDsToAdd) > 0 {
		builder := sq.Insert(toysAndTagsAssociationTableName).
			Columns(toyIDColumnName, tagIDColumnName)
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
		builder := sq.Insert(toysAttachmentsTableName).
			Columns(toyIDColumnName, attachmentLinkColumnName)
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

func (repo *ToysRepository) getToyAttachments(
	ctx context.Context,
	toyID uint64,
	connection *sql.Conn,
) ([]entities.Attachment, error) {
	ctx, span := repo.traceProvider.Span(ctx, tracing.CallerName(tracing.DefaultSkipLevel))
	defer span.End()

	span.AddEvent(repo.spanConfig.Events.Start.Name, repo.spanConfig.Events.Start.Opts...)
	defer span.AddEvent(repo.spanConfig.Events.End.Name, repo.spanConfig.Events.End.Opts...)

	stmt, params, err := sq.
		Select(selectAllColumns).
		From(toysAttachmentsTableName).
		Where(sq.Eq{toyIDColumnName: toyID}).
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

func (repo *ToysRepository) getToyTags(
	ctx context.Context,
	toyID uint64,
	connection *sql.Conn,
) ([]entities.Tag, error) {
	ctx, span := repo.traceProvider.Span(ctx, tracing.CallerName(tracing.DefaultSkipLevel))
	defer span.End()

	span.AddEvent(repo.spanConfig.Events.Start.Name, repo.spanConfig.Events.Start.Opts...)
	defer span.AddEvent(repo.spanConfig.Events.End.Name, repo.spanConfig.Events.End.Opts...)

	stmt, params, err := sq.
		Select(selectAllColumns).
		From(tagsTableName).
		Where(
			sq.Expr(
				idColumnName+" IN (?)",
				sq.Select(tagIDColumnName).
					From(toysAndTagsAssociationTableName).
					Where(sq.Eq{toyIDColumnName: toyID}),
			),
		).
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
