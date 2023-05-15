package repository

import (
	"database/sql"
	"fmt"
	"go-api-meli/model"
)

var QuantityInItems, QuantityInStock, result int64

type cart struct {
	db *sql.DB
}

type CartRepository interface {
	AddProductToCart(cart model.Cart) (uint64, error)
	GetCartById(ID uint64) (model.CartFinallity, error)
	CartFinallity(ID uint64) (model.Purchase, error)
	//Purchase(QuantityInItems uint64, ID uint64) error
}

func NewCartRepository(db *sql.DB) *cart {
	return &cart{db}
}

func (c cart) AddProductToCart(cart model.Cart) (uint64, error) {

	for _, products := range cart.Products {
		statement, _ := c.db.Prepare("insert into tb_cart (idtb_product, quantity) values (?,?)")

		defer statement.Close()
		statement.Exec(products.ID_product, products.Quantity)
	}
	return 0, nil

}

/*
	func (c cart) GetCartById(ID uint64) (model.Detail, error) {
		row, err := c.db.Query("select idtb_cart, quantity from tb_cart where idtb_cart = ?", ID)
		if err != nil {
			return model.Detail{}, err
		}
		defer row.Close()

		var cart model.Detail
		if row.Next() {
			row.Scan(&cart.ID_product, &cart.Quantity)
		}

		return cart, nil

}
*/
func (c cart) GetCartById(ID uint64) (model.CartFinallity, error) {
	row, err := c.db.Query(`select
	p.title,
	p.quantity AS qtd_estoque,
	c.quantity AS qtd_vendida,
	c.date
	from tb_product p
	join tb_cart c
	on p.idtb_product = c.idtb_product
	where c.idtb_cart = ?`, ID)
	if err != nil {
		return model.CartFinallity{}, err
	}
	defer row.Close()

	var cart model.CartFinallity
	if row.Next() {
		row.Scan(&cart.Item, &cart.QuantityInStock, &cart.QuantityInItems, &cart.DateOfPurchase)
	}

	return cart, nil

}
func (c cart) CartFinallity(ID uint64) (model.Purchase, error) {

	row, err := c.db.Query(`select
	p.quantity AS qtd_estoque,
	c.quantity AS qtd_estoque
	from tb_product p
	join tb_cart c
	on p.idtb_product = c.idtb_product
	where c.idtb_cart = ?`, ID)
	if err != nil {
		return model.Purchase{}, err
	}
	defer row.Close()

	var cart model.Purchase
	if row.Next() {
		row.Scan(&cart.QuantityStock, &cart.QuantityItems)
	}
	result = cart.QuantityStock - cart.QuantityItems
	fmt.Println(result)

	return cart, nil

}

/*

func (p products) Purchase(QuantityInItems uint64, ID uint64) error {
	statement, err := p.db.Prepare("update tb_product quantity = ? where idtb_product = ? ")
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(QuantityInItems, ID)
	if err != nil {
		return err

	}
	return nil

}
*/
