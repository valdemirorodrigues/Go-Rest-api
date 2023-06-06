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
}

type cartController struct {
	CartService    service.CartService
	ProductService service.ProductService
}

func NewCartController(service service.CartService, cartService service.ProductService) cartController {
	return cartController{
		CartService:    service,
		ProductService: cartService,
	}
}

func (controller cartController) AddProductToCart(w http.ResponseWriter, r *http.Request) {
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

	for _, product := range cart.Products {
		result, _ := controller.ProductService.ProductValidate(product.IDProduct)
		if result.ID != product.IDProduct {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf("One of the cart products was not found.")))
			return
		}
		if result.QuantityInStock < int64(product.Quantity) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte(fmt.Sprintf("One of the cart products does not have sufficient stock.")))
			return
		}
	}

	result, err := controller.CartService.AddProductToCart(cart)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("There was an error when trying to create a shopping cart.")))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)

}

func (service cartController) GetCartById(w http.ResponseWriter, r *http.Request) {
	parameters := mux.Vars(r)
	ID, _ := strconv.ParseUint(parameters["cartID"], 10, 32)

	if err := service.CartService.CartValidate(ID); err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("Shopping cart with ID %d,%v", ID, " was not found.")))
		return

	}

	cart, err := service.CartService.GetCartById(ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("There was an error when trying to get the shopping cart.")))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cart)
	w.WriteHeader(http.StatusOK)
	return
}

func (service cartController) Checkout(w http.ResponseWriter, r *http.Request) {
	paramters := mux.Vars(r)
	IDCart, _ := strconv.ParseUint(paramters["cartId"], 10, 32)

	if err := service.CartService.CartValidate(IDCart); err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("Shopping cart with ID %d,%v", IDCart, " was not found.")))
		return

	}

	items, _ := service.CartService.CheckoutValidate(IDCart)

	for _, item := range items {
		result, _ := service.ProductService.ProductValidate(uint64(item.IDProduct))
		if result.QuantityInStock < int64(item.QuantityOfItems) {
			fmt.Println(result.QuantityInStock)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("A cart product does not have enough stock. This cart is invalid.")))
			return
		}

	}

	response, err := service.CartService.Checkout(IDCart)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("There was an error when trying to checkout.")))

		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}
