package model

type Details struct {
	ID         int64   `json:id`
	Id_Product uint64  `json: "id_product"`
	Quantity   int8    `json: "quantity"`
	PriceFinal float64 `json: "PriceFinal"`
}

type Cart struct {
	Products []Details `json: "products"`
}
type CartFinallity struct {
	IDProduct       string `json: "idtb_product"`
	IDCat           string `json: "idtb_cart"`
	Item            string `json: "item"`
	QuantityInItems int64  `json: "QuantityInItems"`
	DateOfPurchase  string `json: "DateOfPurchase"`
}
type Purchase struct {
	ID            int64   `json: "id"`
	QuantityStock int64   `json: "QuantityInStock"`
	QuantityItems int64   `json: "QuantityInItems"`
	PriceFinal    float64 `json: "PriceFinal"`
}

type CartResponse struct {
	IdCart          uint64  `json: "idtb_cart"`
	Title           string  `json: "title"`
	QuantityOfItems uint64  `json: "QuantityOfItems"`
	PriceFinal      float64 `json: "PriceFinal"`
}
