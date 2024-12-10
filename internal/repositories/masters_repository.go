package repositories

import (
	"github.com/DKhorkov/hmtm-toys/internal/entities"
	"github.com/DKhorkov/libs/db"
)

type CommonMastersRepository struct {
	dbConnector db.Connector
}

func (repo *CommonMastersRepository) GetAllMasters() ([]entities.Master, error) {
	connection := repo.dbConnector.GetConnection()
	rows, err := connection.Query(
		`
			SELECT * 
			FROM masters
		`,
	)

	if err != nil {
		return nil, err
	}

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

	if err = rows.Close(); err != nil {
		return nil, err
	}

	return masters, nil
}

func (repo *CommonMastersRepository) GetMasterByUserID(userID uint64) (*entities.Master, error) {
	master := &entities.Master{}
	columns := db.GetEntityColumns(master)
	connection := repo.dbConnector.GetConnection()
	err := connection.QueryRow(
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

func (repo *CommonMastersRepository) GetMasterByID(id uint64) (*entities.Master, error) {
	master := &entities.Master{}
	columns := db.GetEntityColumns(master)
	connection := repo.dbConnector.GetConnection()
	err := connection.QueryRow(
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

func (repo *CommonMastersRepository) RegisterMaster(masterData entities.RegisterMasterDTO) (uint64, error) {
	var masterID uint64
	connection := repo.dbConnector.GetConnection()
	err := connection.QueryRow(
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

func NewCommonMastersRepository(dbConnector db.Connector) *CommonMastersRepository {
	return &CommonMastersRepository{dbConnector: dbConnector}
}
