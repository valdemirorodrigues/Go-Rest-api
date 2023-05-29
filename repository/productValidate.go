package repository

import (
	"database/sql"
	"go-api-meli/model"
)

type product struct {
	db *sql.DB
}

type RepositoryProductValidation interface {
	ValidateProduct(productID uint64) error
}

func NewProductRepositoryValidation(db *sql.DB) *products {
	return &products{db}
}

func (p products) ValidateProduct(productID uint64) error {
	row, err := p.db.Query("select idtb_product from tb_product where idtb_product = ?", productID)
	if err != nil {
		return err
	}
	defer row.Close()

	var product model.Product

	if row.Next() {
		if err = row.Scan(
			&product.ID,
		); err != nil {
			return err
		}
	}

	return err
}
