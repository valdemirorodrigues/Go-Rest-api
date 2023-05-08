package service

import (
	"database/sql"
	"go-api-meli/database"
	"go-api-meli/model"
	"go-api-meli/repository"
)

type products struct {
	db *sql.DB
}

func CreateProductCreateProduct(product model.Product) (uint64, error) {
	db, err := database.Connection()
	if err != nil {
		return 0, err
	}

	return repository.RepositoryProduct(db).CreateProduct(product)
}
