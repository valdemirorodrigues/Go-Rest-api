package repository

import (
	"database/sql"
	"fmt"
	"go-api-meli/model"
)

var QuantityInItems, QuantityInStock, Result int64
var ValueFinal float64

type cart struct {
	db *sql.DB
}

type CartRepository interface {
	AddProductToCart(cart model.Cart) (model.Cart, error)
	GetCartById(ID uint64) ([]model.Detail, error)
	CartFinallity(ID uint64) (model.Purchase, error)
	Purchase(Result, ID uint64) error
	InsertTbProductTbcart(codeTbProduct uint64, codeTbCart uint64) (uint64, error)
}

func NewCartRepository(db *sql.DB) *cart {
	return &cart{db}
}

func (c cart) AddProductToCart(cart model.Cart) (model.Cart, error) {

	for _, products := range cart.Products {
		statement, err := c.db.Prepare("insert into tb_cart (id_product, quantity_of_items) values (?,?)")
		if err != nil {
			return model.Cart{}, err
		}

		defer statement.Close()

		statement.Exec(products.ID_product, products.Quantity)
		if err != nil {
			return model.Cart{}, err
		}
	}
	return cart, nil
}

// vai receber o id do produto e do carrinho via postman.
func (c cart) InsertTbProductTbcart(codeTbProduct uint64, codeTbCart uint64) (uint64, error) {

	statement, err := c.db.Prepare("insert into tb_cart_tb_product (codetb_product, codetb_cart) values (?,?)")
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	result, err := statement.Exec(codeTbProduct, codeTbCart)
	if err != nil {
		return 0, err
	}

	ID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(ID), err

}

func (c cart) GetCartById(ID uint64) ([]model.Detail, error) {
	row, err := c.db.Query(`select 
	c.idtb_cart,
	c.quantity_of_items
	from 
	tb_product p
	join tb_cart_tb_product cp
	on cp.codetb_product = p.idtb_product
	join tb_cart c
	on c.idtb_cart = cp.codetb_cart
	where c.idtb_cart = ?`, ID)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	var cart []model.Detail
	for row.Next() {
		var crt model.Detail

		if err = row.Scan(&crt.ID_product, &crt.Quantity); err != nil {
			return nil, err
		}
		cart = append(cart, crt)
	}
	return cart, nil
}

// selecionar carrinho final passando o id da tb_cart_tb_produc
func (c cart) CartFinallity(ID uint64) (model.Purchase, error) {

	row, err := c.db.Query(`select 
		p.quantity as QuantityStock,
		c.quantity as QuantityItems,
		p.price
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
		row.Scan(&cart.QuantityStock, &cart.QuantityItems, &cart.PriceFinal)
	}
	Result = cart.QuantityStock - cart.QuantityItems
	ValueFinal = float64(cart.QuantityItems) * cart.PriceFinal
	fmt.Println(Result, ValueFinal)
	c.Purchase(uint64(Result), ID)
	c.UpdatePurchaseAmount(uint64(ValueFinal), ID)

	return cart, nil

}

// vai atualizar a coluna de quantidade na tb_product e o preco na tb_cart
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

func (c cart) UpdatePurchaseAmount(Result, ID uint64) error {

	statement, err := c.db.Prepare(
		`update tb_cart c
		join tb_cart_tb_product cp
		on cp.codetb_cart = c.idtb_cart
		join tb_product p
		on p.idtb_product = cp.codetb_product
		set c.price_final = ?
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
