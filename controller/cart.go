package controller

import (
	"encoding/json"
	"fmt"
	"go-api-meli/model"
	"go-api-meli/service"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type CartController interface {
	AddProductToCart(w http.ResponseWriter, r *http.Request)
	GetCartById(w http.ResponseWriter, r *http.Request)
}

type cartController struct {
	CartService service.CartService
}

func NewCartController(service service.CartService) cartController {
	return cartController{
		CartService: service,
	}
}

func (service cartController) AddProductToCart(w http.ResponseWriter, r *http.Request) {
	request, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	var cart model.Cart
	if err = json.Unmarshal(request, &cart); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if ID, err := service.CartService.AddProductToCart(cart); err != nil {
		if err != nil {

			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(cart)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf("Product id %d", ID)))

		return
	}
}

func (service cartController) GetCartById(w http.ResponseWriter, r *http.Request) {
	paramters := mux.Vars(r)
	ID, _ := strconv.ParseUint(paramters["id"], 10, 32)
	cart, err := service.CartService.GetCartById(ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(cart)
	w.WriteHeader(http.StatusOK)
	return
}

/*
func CartFinallity(w http.ResponseWriter, r *http.Request) {

	paramters := mux.Vars(r)
	ID, _ := strconv.ParseUint(paramters["id"], 10, 32)

	row, err := db.Query("select * from tb_product p join tb_cart c on p.idtb_product = c.idtb_product where c.idtb_cart = ?", ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
*/
