package repositories

import (
	"database/sql"
)

type ProductRepository struct {
	db *sql.DB
}

// NewProductRepository returns a minimal product repository stub.
func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}
