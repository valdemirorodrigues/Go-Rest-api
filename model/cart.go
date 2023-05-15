package model

type Detail struct {
	ID_product int64 `json: "id_product"`
	Quantity   int8  `json: "quantity"`
}

type Cart struct {
	Products []Detail
}
type CartFinallity struct {
	Item            string `json: "item"`
	QuantityInStock int64  `json: "QuantityInStock"`
	QuantityInItems int64  `json: "QuantityInItems"`
	DateOfPurchase  string `json: "DateOfPurchase"`
}
type Purchase struct {
	QuantityStock int64 `json: "QuantityInStock"`
	QuantityItems int64 `json: "QuantityInItems"`
}
