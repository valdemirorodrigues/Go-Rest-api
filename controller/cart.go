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
	Checkout(w http.ResponseWriter, r *http.Request)
	InsertTbProductTbcart(w http.ResponseWriter, r *http.Request)
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
	err = service.CartService.AddProductToCart(cart)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("Checkout realizado com sucesso")

}

func (service cartController) GetCartById(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)
	ID, _ := strconv.ParseUint(parameters["cartID"], 10, 32)
	cart, err := service.CartService.GetCartById(ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cart)
	w.WriteHeader(http.StatusOK)
	return
}

func (service cartController) Checkout(w http.ResponseWriter, r *http.Request) {
	paramters := mux.Vars(r)
	ID, err := strconv.ParseUint(paramters["cartId"], 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = service.CartService.Checkout(ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Checkout realizado com sucesso")

}

func (service cartController) InsertTbProductTbcart(w http.ResponseWriter, r *http.Request) {
	paramters := mux.Vars(r)
	codeTbProduct, err := strconv.ParseUint(paramters["codeTbProduct"], 10, 32)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("error processing product id")))
		return
	}
	codeTbCart, _ := strconv.ParseUint(paramters["codeTbCart"], 10, 32)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("error processing cart id")))
		return
	}

	ID, err := service.CartService.InsertTbProductTbcart(codeTbProduct, codeTbCart)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(fmt.Sprintf("Cart created id %d", ID))
}
