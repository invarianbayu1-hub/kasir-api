package repositories

import (
	"database/sql"
	"errors"
	"kasir-api/models"
)

type InventoryRepository struct {
	db *sql.DB
}

func NewInventoryRepository(db *sql.DB) *InventoryRepository {
	return &InventoryRepository{db: db}
}

func (repo *InventoryRepository) GetAll() ([]models.Inventory, error) {
	query := "SELECT id, name, description FROM inventory"
	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	inventories := make([]models.Inventory, 0)
	for rows.Next() {
		var p models.Inventory
		err := rows.Scan(&p.ID, &p.Name, &p.Description)
		if err != nil {
			return nil, err
		}
		inventories = append(inventories, p)
	}

	return inventories, nil
}

func (repo *InventoryRepository) Create(inventory *models.Inventory) error {
	query := "INSERT INTO inventory (name, description) VALUES ($1, $2) RETURNING id"
	err := repo.db.QueryRow(query, inventory.Name, inventory.Description).Scan(&inventory.ID)
	return err
}

// GetByID - ambil inventory by ID
func (repo *InventoryRepository) GetByID(id int) (*models.Inventory, error) {
	query := "SELECT id, name, description FROM inventory WHERE id = $1"

	var p models.Inventory
	err := repo.db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Description)
	if err == sql.ErrNoRows {
		return nil, errors.New("inventory tidak ditemukan")
	}
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (repo *InventoryRepository) Update(inventory *models.Inventory) error {
	query := "UPDATE inventory SET name = $1, description = $2 WHERE id = $3"
	result, err := repo.db.Exec(query, inventory.Name, inventory.Description, inventory.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("inventory tidak ditemukan")
	}

	return nil
}

func (repo *InventoryRepository) Delete(id int) error {
	query := "DELETE FROM inventory WHERE id = $1"
	result, err := repo.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("inventory tidak ditemukan")
	}

	return err
}
