package model

type Detail struct {
	ID         int64 `json:id`
	ID_product int64 `json: "id_product"`
	Quantity   int8  `json: "quantity"`
}

type Cart struct {
	Products []Detail
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
	PriceFinal    float64 `json: "PriceFina"`
}
