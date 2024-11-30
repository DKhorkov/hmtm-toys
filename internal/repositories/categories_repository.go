package repositories

import (
	"github.com/DKhorkov/hmtm-toys/pkg/entities"
	"github.com/DKhorkov/libs/db"
)

type CommonCategoriesRepository struct {
	dbConnector db.Connector
}

func (repo *CommonCategoriesRepository) GetAllCategories() ([]*entities.Category, error) {
	connection := repo.dbConnector.GetConnection()
	rows, err := connection.Query(
		`
			SELECT * 
			FROM categories
		`,
	)

	if err != nil {
		return nil, err
	}

	var categories []*entities.Category
	for rows.Next() {
		category := &entities.Category{}
		columns := db.GetEntityColumns(category)
		err = rows.Scan(columns...)
		if err != nil {
			return nil, err
		}

		categories = append(categories, category)
	}

	if err = rows.Close(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (repo *CommonCategoriesRepository) GetCategoryByID(id uint32) (*entities.Category, error) {
	category := &entities.Category{}
	columns := db.GetEntityColumns(category)
	connection := repo.dbConnector.GetConnection()
	err := connection.QueryRow(
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

	return category, nil
}

func NewCommonCategoriesRepository(dbConnector db.Connector) *CommonCategoriesRepository {
	return &CommonCategoriesRepository{dbConnector: dbConnector}
}
