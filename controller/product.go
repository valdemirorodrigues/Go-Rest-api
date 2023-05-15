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

type ProductController interface {
	CreateProduct(w http.ResponseWriter, r *http.Request)
	GetProducts(w http.ResponseWriter, r *http.Request)
	GetProductById(w http.ResponseWriter, r *http.Request)
	DeleteProduct(w http.ResponseWriter, r *http.Request)
	UpdateProduct(w http.ResponseWriter, r *http.Request)
}

// injetando o service no controller
type productController struct {
	ProductService service.ProductService
}

func NewProductController(service service.ProductService) productController {
	return productController{
		ProductService: service,
	}
}
func (service productController) CreateProduct(w http.ResponseWriter, r *http.Request) {
	request, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	var product model.Product

	if err = json.Unmarshal(request, &product); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if ID, err := service.ProductService.CreateProduct(product); err != nil {
		if err != nil {

			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(product)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf("Product id %d", ID)))

	}
}
func (service productController) GetProducts(w http.ResponseWriter, r *http.Request) {

	product, err := service.ProductService.GetAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(product)
	w.WriteHeader(http.StatusOK)

}

func (service productController) GetProductById(w http.ResponseWriter, r *http.Request) {
	paramters := mux.Vars(r)
	ID, err := strconv.ParseUint(paramters["productID"], 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	product, err := service.ProductService.GetById(ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(product)
	w.WriteHeader(http.StatusOK)
}

func (service productController) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	paramters := mux.Vars(r)
	ID, err := strconv.ParseInt(paramters["productID"], 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	request, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	var product model.Product

	if err = json.Unmarshal(request, &product); err != nil {
		w.Write([]byte(fmt.Sprintf("Error converting product object to struct.")))
		return
	}

	if err = service.ProductService.UpdateProduct(uint64(ID), product); err != nil {
		w.Write([]byte(fmt.Sprintf("Error when updating product.")))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Product updated successfully.")))
}

func (service productController) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	paramters := mux.Vars(r)
	ID, err := strconv.ParseUint(paramters["productID"], 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err = service.ProductService.DeleteProduct(uint64(ID)); err != nil {
		w.Write([]byte(fmt.Sprintf("Error deleting product.")))
		return
	}

	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte(fmt.Sprintf("Product deleted successfully.")))

}
