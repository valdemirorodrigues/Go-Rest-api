package model

type Detail struct {
	IdProduct int64 `json: "id_product"`
	Quantity  int8  `json: "quantity"`
}

type Cart struct {
	Products []Detail
}
