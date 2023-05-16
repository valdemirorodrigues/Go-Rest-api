package repository

import (
	"database/sql"
	"fmt"
	"go-api-meli/model"
)

var QuantityInItems, QuantityInStock, Result int64

type cart struct {
	db *sql.DB
}

type CartRepository interface {
	AddProductToCart(cart model.Cart) (uint64, error)
	GetCartById(ID uint64) (model.CartFinallity, error)
	CartFinallity(ID uint64) (model.Purchase, error)
	Purchase(Result, ID uint64) error
}

func NewCartRepository(db *sql.DB) *cart {
	return &cart{db}
}

func (c cart) AddProductToCart(cart model.Cart) (uint64, error) {

	for _, products := range cart.Products {
		statement, _ := c.db.Prepare("insert into tb_cart (idtb_product, quantity) values (?,?)")

		defer statement.Close()

		result, err := statement.Exec(products.ID_product, products.Quantity)
		if err != nil {
			return 0, err
		}

		ID, err := result.LastInsertId()
		if err != nil {
			return 0, err
		}
		fmt.Println(products.ID_product)

		return uint64(ID), nil
	}
	return 0, nil

}

func (c cart) GetCartById(ID uint64) (model.CartFinallity, error) {
	row, err := c.db.Query(`select 
	p.idtb_product, 
	c.idtb_cart,
	p.title,
	c.quantity as QuantityInItems,
	c.date 
	from 
	tb_product p
	join tb_cart_tb_product cp
	on cp.codetb_product = p.idtb_product
	join tb_cart c
	on c.idtb_cart = cp.codetb_cart
	where cp.idtb_cart_tb_produc = ?`, ID)
	if err != nil {
		return model.CartFinallity{}, err
	}
	defer row.Close()

	var cart model.CartFinallity
	if row.Next() {
		row.Scan(&cart.IDProduct, &cart.IDCat, &cart.Item, &cart.QuantityInItems, &cart.DateOfPurchase)
	}

	fmt.Println(cart.IDCat, cart.IDProduct)

	return cart, nil

}
func (c cart) CartFinallity(ID uint64) (model.Purchase, error) {

	row, err := c.db.Query(`select 
	p.quantity as QuantityStock,
	c.quantity as QuantityItems
	from 
	tb_product p
	join tb_cart_tb_product cp
	on cp.codetb_product = p.idtb_product
	join tb_cart c
	on c.idtb_cart = cp.codetb_cart
	where cp.idtb_cart_tb_produc = ?`, ID)
	if err != nil {
		return model.Purchase{}, err
	}
	defer row.Close()

	var cart model.Purchase
	if row.Next() {
		row.Scan(&cart.QuantityStock, &cart.QuantityItems)
	}
	Result = cart.QuantityStock - cart.QuantityItems
	fmt.Println(Result)
	c.Purchase(uint64(Result), ID)

	return cart, nil

}

func (c cart) Purchase(Result, ID uint64) error {
	statement, err := c.db.Prepare(
		`update tb_product p
	join tb_cart_tb_product cp
	on cp.codetb_product = p.idtb_product
	join tb_cart c
	on c.idtb_cart = cp.codetb_cart
	set p.quantity = ?
	where cp.idtb_cart_tb_produc = ? `)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(Result, ID)
	if err != nil {
		return err

	}
	return nil

}

// vai receber o id do produto e do carrinho via postman
func (c cart) InsertTbcartTbProduct(codeTbProduct int64, codeTbCart int64) (uint64, error) {

	statement, _ := c.db.Prepare("insert into tb_cart_tb_product (codetb_product, codetb_cart) values (?,?)")

	defer statement.Close()
	result, err := statement.Exec(codeTbProduct, codeTbCart)
	if err != nil {
		return 0, err
	}

	ID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(ID), nil

}
