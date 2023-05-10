package main

import (
	"fmt"
	"go-api-meli/controller"
	"go-api-meli/database"
	"go-api-meli/repository"
	"go-api-meli/service"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	build := buildControllers()
	routers(build)

}

func routers(controllers *Controllers) {

	router := mux.NewRouter()
	router.HandleFunc("/products", controllers.productController.CreateProduct).Methods(http.MethodPost)
	router.HandleFunc("/products", controllers.productController.GetProducts).Methods(http.MethodGet)
	router.HandleFunc("/products/{productID}", controllers.productController.GetProductById).Methods(http.MethodGet)
	router.HandleFunc("/products/{productID}", controllers.productController.DeleteProduct).Methods(http.MethodDelete)
	router.HandleFunc("/products/{productID}", controllers.productController.UpdateProduct).Methods(http.MethodPut)
	/*
		router.HandleFunc("/cart", controller.AddProductToCart).Methods(http.MethodPost)
		router.HandleFunc("/cart/{id}", controller.GetCartById).Methods(http.MethodGet)
	*/

	fmt.Println("servidor rodando")
	log.Fatal(http.ListenAndServe(":5000", router))
}

type Controllers struct {
	productController controller.ProductController
}

func buildControllers() *Controllers {
	db, _ := database.Connection()
	productRepository := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepository)
	return &Controllers{
		productController: controller.NewProductController(productService),
	}
}
