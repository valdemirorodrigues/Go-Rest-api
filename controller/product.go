package controller

import (
	"encoding/json"
	"fmt"
	"go-api-meli/database"
	"go-api-meli/model"
	"go-api-meli/repository"
	"go-api-meli/service"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ProductController interface {
	CreateProduct(w http.ResponseWriter, r *http.Request)
}
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
func GetProducts(w http.ResponseWriter, r *http.Request) {
	db, err := database.Connection()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer db.Close()

	repository := repository.RepositoryProduct(db)
	product, err := repository.GetAll()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(product)
	w.WriteHeader(http.StatusOK)

}

func GetProductById(w http.ResponseWriter, r *http.Request) {
	paramters := mux.Vars(r)
	ID, err := strconv.ParseUint(paramters["productID"], 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db, err := database.Connection()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer db.Close()

	repository := repository.RepositoryProduct(db)
	product, err := repository.GetById(ID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(product)
	w.WriteHeader(http.StatusOK)
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
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

	db, err := database.Connection()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer db.Close()

	repository := repository.RepositoryProduct(db)
	if err = repository.UpdateProduct(uint64(ID), product); err != nil {
		w.Write([]byte(fmt.Sprintf("Error when updating product.")))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Product updated successfully.")))
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	paramters := mux.Vars(r)
	ID, err := strconv.ParseInt(paramters["productID"], 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db, err := database.Connection()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	repository := repository.RepositoryProduct(db)
	if err = repository.Delete(uint64(ID)); err != nil {
		w.Write([]byte(fmt.Sprintf("Error deleting productt.")))
		return
	}

	defer db.Close()

	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte(fmt.Sprintf("Product deleted successfully.")))

}
