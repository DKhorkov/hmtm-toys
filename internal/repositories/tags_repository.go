package repositories

import (
	"github.com/DKhorkov/hmtm-toys/internal/entities"
	"github.com/DKhorkov/libs/db"
)

type CommonTagsRepository struct {
	dbConnector db.Connector
}

func (repo *CommonTagsRepository) GetAllTags() ([]entities.Tag, error) {
	connection := repo.dbConnector.GetConnection()
	rows, err := connection.Query(
		`
			SELECT * 
			FROM tags
		`,
	)

	if err != nil {
		return nil, err
	}

	var tags []entities.Tag
	for rows.Next() {
		tag := entities.Tag{}
		columns := db.GetEntityColumns(&tag) // Only pointer to use rows.Scan() successfully
		err = rows.Scan(columns...)
		if err != nil {
			return nil, err
		}

		tags = append(tags, tag)
	}

	if err = rows.Close(); err != nil {
		return nil, err
	}

	return tags, nil
}

func (repo *CommonTagsRepository) GetToyTags(toyID uint64) ([]entities.Tag, error) {
	connection := repo.dbConnector.GetConnection()
	rows, err := connection.Query(
		`
			SELECT * 
			FROM tags AS t
			WHERE t.id IN (
			    SELECT ta.tag_id
				FROM toys_and_tags_associations AS ta
				WHERE ta.toy_id = $1
			)
		`,
		toyID,
	)

	if err != nil {
		return nil, err
	}

	var tags []entities.Tag
	for rows.Next() {
		tag := entities.Tag{}
		columns := db.GetEntityColumns(&tag) // Only pointer to use rows.Scan() successfully
		err = rows.Scan(columns...)
		if err != nil {
			return nil, err
		}

		tags = append(tags, tag)
	}

	if err = rows.Close(); err != nil {
		return nil, err
	}

	return tags, nil
}

func (repo *CommonTagsRepository) GetTagByID(id uint32) (*entities.Tag, error) {
	tag := &entities.Tag{}
	columns := db.GetEntityColumns(tag)
	connection := repo.dbConnector.GetConnection()
	err := connection.QueryRow(
		`
			SELECT * 
			FROM tags AS t
			WHERE t.id = $1
		`,
		id,
	).Scan(columns...)

	if err != nil {
		return nil, err
	}

	return tag, nil
}

func NewCommonTagsRepository(dbConnector db.Connector) *CommonTagsRepository {
	return &CommonTagsRepository{dbConnector: dbConnector}
}
