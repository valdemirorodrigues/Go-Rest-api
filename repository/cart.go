package repository

import (
	"database/sql"
	"go-api-meli/model"
)

type cart struct {
	db *sql.DB
}

type CartRepository interface {
	AddProductToCart(cart model.Cart) (uint64, error)
	GetCartById(ID uint64) (model.Detail, error)
}

func NewCartRepository(db *sql.DB) *cart {
	return &cart{db}
}

func (c cart) AddProductToCart(cart model.Cart) (uint64, error) {

	for _, CarProduct := range cart.Products {
		statement, err := c.db.Prepare("insert into tb_cart (idtb_product, quantity) values (?,?)")
		if err != nil {
			return 0, err
		}
		defer statement.Close()

		result, err := statement.Exec(CarProduct.IdProduct, CarProduct.Quantity)
		if err != nil {
			return 0, err
		}
		ID, err := result.LastInsertId()
		if err != nil {
			return 0, err
		}
		return uint64(ID), nil

	}
	return 0, nil

}

func (c cart) GetCartById(ID uint64) (model.Detail, error) {
	row, err := c.db.Query("select idtb_cart, quantity from tb_cart where idtb_cart = ?", ID)
	if err != nil {
		return model.Detail{}, err
	}
	defer row.Close()

	var cart model.Detail
	if row.Next() {
		row.Scan(&cart.IdProduct, &cart.Quantity)
	}

	return cart, nil

}
