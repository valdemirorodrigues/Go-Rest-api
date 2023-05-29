package repository

import (
	"database/sql"
	"go-api-meli/model"
	"time"
)

type cart struct {
	db *sql.DB
}

type CartRepository interface {
	AddProductToCart(products []model.Details) error
	GetCartById(ID uint64) ([]model.CartResponse, error)
	Checkout(ID uint64) (model.Purchase, error)
	UpdateInventoryColumn(quantity int64, idCart uint64) error
	InsertTbProductTbcart(codeTbProduct uint64, codeTbCart uint64) (uint64, error)
	SubtractOfItems(ID uint64) (int64, error)
}

func NewCartRepository(db *sql.DB) *cart {
	return &cart{db}
}

// Adiciona o carrinho com produtos
func (c cart) AddProductToCart(lp []model.Details) error {
	var IdCart, err = c.CreateCart()
	if err != nil {
		return err
	}
	for _, product := range lp {

		statement, err := c.db.Prepare("insert into tb_cart_tb_product (codetb_product, codetb_cart, quantity_of_items,  price_final) values (?, ?,?,?)")

		if err != nil {
			return err
		}

		defer statement.Close()

		statement.Exec(product.Id_Product, IdCart, product.Quantity, product.PriceFinal)
	}

	return nil
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

func (c cart) GetCartById(ID uint64) ([]model.CartResponse, error) {
	row, err := c.db.Query(`select
	c.idtb_cart,
	cp.quantity_of_items, 
	cp.price_final,
	p.title
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

	var cart []model.CartResponse
	for row.Next() {
		var crt model.CartResponse

		if err = row.Scan(&crt.IdCart, &crt.QuantityOfItems, &crt.PriceFinal, &crt.Title); err != nil {
			return nil, err
		}
		cart = append(cart, crt)
	}
	return cart, nil
}

// Realiza o calculo e retorna o valor total do carrinnho
func (c cart) Checkout(idCart uint64) (model.Purchase, error) {

	result, err := c.SubtractOfItems(idCart)
	if err != nil {
		return model.Purchase{}, err
	}

	return model.Purchase{QuantityStock: result}, nil

}

// Realiza a subtracao dos itens em estoque com os do carrinho
func (c cart) SubtractOfItems(idCart uint64) (int64, error) {

	rows, err := c.db.Query(`select
	p.quantity_in_stock, cp.quantity_of_items
	from
	tb_product p
	join tb_cart_tb_product cp
	on cp.codetb_product = p.idtb_product
	join tb_cart c
	on c.idtb_cart = cp.codetb_cart
	where c.idtb_cart = ?`, idCart)
	if err != nil {
		return 0, err

	}
	defer rows.Close()

	var sumQuantityFinal int64
	var product model.Purchase
	for rows.Next() {
		if err = rows.Scan(&product.QuantityStock, &product.QuantityItems); err != nil {
			return 0, err
		}
		sumQuantityFinal += product.QuantityItems
	}
	sumQuantityFinal = product.QuantityStock - sumQuantityFinal
	c.UpdateInventoryColumn(sumQuantityFinal, idCart)

	return sumQuantityFinal, nil

}

// vai atualizar a coluna de quantidade realizar um for
func (c cart) UpdateInventoryColumn(quantity int64, idCart uint64) error {

	statement, err := c.db.Prepare(
		`update tb_cart c
			join tb_cart_tb_product cp
			on cp.codetb_cart = c.idtb_cart
			join tb_product p
			on p.idtb_product = cp.codetb_product
			set p.quantity_in_stock = ?
			where c.idtb_cart = ? `)
	if err != nil {
		return err
	}
	defer statement.Close()
	statement.Exec(quantity, idCart)

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
