package repository

import (
	"database/sql"
	"fmt"
	"go-api-meli/model"
	"time"
)

var ValueFinal float64
var Stock int64

type cart struct {
	db *sql.DB
}

type CartRepository interface {
	AddProductToCart(cart model.Cart) (model.Cart, error)
	GetCartById(ID uint64) ([]model.Detail, error)
	Checkout(ID uint64) (model.Purchase, error)
	UpdateInventoryColumn(quantityStock int64, ID int64) error
	InsertTbProductTbcart(codeTbProduct uint64, codeTbCart uint64) (uint64, error)
	SubtractOfItems(ID uint64) ([]model.Purchase, error)
}

func NewCartRepository(db *sql.DB) *cart {
	return &cart{db}
}

// Adiciona o carrinho com produtos
func (c cart) AddProductToCart(cart model.Cart) (model.Cart, error) {
	var ID_cart, _ = c.CreateCart()
	detail := []model.Detail{}

	for _, products := range cart.Products {

		statement, err := c.db.Prepare("insert into tb_cart_tb_product (codetb_product, codetb_cart, quantity_of_items) values (?,?,?)")

		if err != nil {
			return model.Cart{}, err
		}

		defer statement.Close()

		statement.Exec(products.ID_product, ID_cart, products.Quantity)
		detail = append(detail, model.Detail{ID: ID_cart, ID_product: products.ID_product, Quantity: products.Quantity})

	}
	cartResponse := model.Cart{
		Products: detail,
	}
	return cartResponse, nil
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

// Realiza o calculo e retorna o valor total do carrinnho
func (c cart) Checkout(ID uint64) (model.Purchase, error) {

	row, err := c.db.Query(`select
	sum(cp.quantity_of_items * p.price)
	from
	tb_product p
	join tb_cart_tb_product cp
	on cp.codetb_product = p.idtb_product
	join tb_cart c
	on c.idtb_cart = cp.codetb_cart
	where c.idtb_cart = ?`, ID)

	if err != nil {
		return model.Purchase{}, err
	}
	defer row.Close()

	var cart model.Purchase
	if row.Next() {
		row.Scan(&cart.PriceFinal)
	}

	ValueFinal = cart.PriceFinal
	fmt.Println(cart.PriceFinal)

	c.UpdatePurchaseAmount(uint64(ValueFinal), ID)
	c.SubtractOfItems(ID)
	//fmt.Println(ID)

	return cart, nil

}

// Realiza a subtracao dos itens em estoque com os do carrinho
func (c cart) SubtractOfItems(ID uint64) ([]model.Purchase, error) {

	rows, err := c.db.Query(`select
	p.idtb_product,
	p.quantity_in_stock - cp.quantity_of_items
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
	defer rows.Close()

	var products []model.Purchase

	for rows.Next() {
		var product model.Purchase

		if err = rows.Scan(&product.ID, &product.QuantityStock); err != nil {
			return nil, err
		}

		fmt.Println(product.ID, product.QuantityStock)

		products = append(products, product)
		c.UpdateInventoryColumn(product.QuantityStock, product.ID)

	}

	return products, nil

}

// vai atualizar a coluna de quantidade realizar um for
func (c cart) UpdateInventoryColumn(quantityStock int64, ID int64) error {

	statement, err := c.db.Prepare(
		`update tb_cart c
			join tb_cart_tb_product cp
			on cp.codetb_cart = c.idtb_cart
			join tb_product p
			on p.idtb_product = cp.codetb_product
			set p.quantity_in_stock = ?
			where p.idtb_product = ? `)
	if err != nil {
		return err
	}
	//	fmt.Println(param)
	defer statement.Close()
	statement.Exec(quantityStock, ID)

	//}

	return nil

}

func (c cart) UpdatePurchaseAmount(ValueFinal, ID uint64) error {
	statement, err := c.db.Prepare(
		`update tb_cart_tb_product cp
		join tb_cart c
		on cp.codetb_cart = c.idtb_cart
		set cp.price_final = ?
		where c.idtb_cart = ? `)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(ValueFinal, ID)
	if err != nil {
		return err

	}

	return nil

}

func (c cart) CreateCart() (int64, error) {
	date := time.Now()
	statement, err := c.db.Prepare(
		"insert into tb_cart (date) values (?)",
	)
	if err != nil {
		return 0, err
	}

	defer statement.Close()

	result, err := statement.Exec(date)
	if err != nil {
		return 0, err
	}

	ID, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	return ID, err

}
