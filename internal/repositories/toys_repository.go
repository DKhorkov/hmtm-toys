package repositories

import (
	"fmt"
	"strings"

	"github.com/DKhorkov/hmtm-toys/pkg/entities"
	"github.com/DKhorkov/libs/db"
)

type CommonToysRepository struct {
	dbConnector db.Connector
}

func (repo *CommonToysRepository) GetAllToys() ([]*entities.Toy, error) {
	connection := repo.dbConnector.GetConnection()
	rows, err := connection.Query(
		`
			SELECT * 
			FROM toys
		`,
	)

	if err != nil {
		return nil, err
	}

	var toys []*entities.Toy
	for rows.Next() {
		toy := &entities.Toy{}
		columns := db.GetEntityColumns(toy)
		columns = columns[:len(columns)-1] // not to paste tags field ([]*Tag) to Scan function.
		err = rows.Scan(columns...)
		if err != nil {
			return nil, err
		}

		toys = append(toys, toy)
	}

	if err = rows.Close(); err != nil {
		return nil, err
	}

	return toys, nil
}

func (repo *CommonToysRepository) GetMasterToys(masterID uint64) ([]*entities.Toy, error) {
	connection := repo.dbConnector.GetConnection()
	rows, err := connection.Query(
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

	var toys []*entities.Toy
	for rows.Next() {
		toy := &entities.Toy{}
		columns := db.GetEntityColumns(toy)
		columns = columns[:len(columns)-1] // not to paste tags field ([]*Tag) to Scan function.
		err = rows.Scan(columns...)
		if err != nil {
			return nil, err
		}

		toys = append(toys, toy)
	}

	if err = rows.Close(); err != nil {
		return nil, err
	}

	return toys, nil
}

func (repo *CommonToysRepository) GetToyByID(id uint64) (*entities.Toy, error) {
	toy := &entities.Toy{}
	columns := db.GetEntityColumns(toy)
	columns = columns[:len(columns)-1] // not to paste tags field ([]*Tag) to Scan function.
	connection := repo.dbConnector.GetConnection()
	err := connection.QueryRow(
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

func (repo *CommonToysRepository) AddToy(toyData entities.AddToyDTO) (uint64, error) {
	var toyID uint64
	transaction, err := repo.dbConnector.GetTransaction()
	if err != nil {
		return 0, err
	}

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

	// Bulk insert of Toy's Tags.
	toysAndTagsInsertPlaceholders := make([]string, 0, len(toyData.TagsIDs))
	toysAndTagsInsertValues := make([]interface{}, 0, len(toyData.TagsIDs))
	for index, tagID := range toyData.TagsIDs {
		toysAdnTagsInsertPlaceholder := fmt.Sprintf("($%d,$%d)",
			index*2+1,
			index*2+2,
		)

		toysAndTagsInsertPlaceholders = append(toysAndTagsInsertPlaceholders, toysAdnTagsInsertPlaceholder)
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

	err = transaction.Commit()
	if err != nil {
		return 0, err
	}

	return toyID, nil
}

func NewCommonToysRepository(dbConnector db.Connector) *CommonToysRepository {
	return &CommonToysRepository{dbConnector: dbConnector}
}
