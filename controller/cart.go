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
	CartFinallity(w http.ResponseWriter, r *http.Request)
	InsertTbcartTbProduct(w http.ResponseWriter, r *http.Request)
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
	ID, err := service.CartService.AddProductToCart(cart)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	/*
		for _, products := range cart.Products {
			fmt.Println("-------", products.ID_product)
		}
	*/

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Cart id %d", ID)))
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

func (service cartController) CartFinallity(w http.ResponseWriter, r *http.Request) {
	paramters := mux.Vars(r)
	ID, _ := strconv.ParseUint(paramters["id"], 10, 32)

	row, err := service.CartService.CartFinallity(uint64(ID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

	}
	json.NewEncoder(w).Encode(row)
	w.WriteHeader(http.StatusOK)
	return
}
func (service cartController) InsertTbcartTbProduct(w http.ResponseWriter, r *http.Request) {
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

	ID, err := service.CartService.InsertTbcartTbProduct(codeTbProduct, codeTbCart)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Cart created id %d", ID)))

}
