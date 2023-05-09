package service

import (
	"go-api-meli/model"
	"go-api-meli/repository"
)

type ProductService interface {
	CreateProduct(product model.Product) (uint64, error)
}
type productService struct {
	Repository repository.Repository
}

func NewProductService(repo repository.Repository) productService {
	return productService{
		Repository: repo,
	}
}
func (pr productService) CreateProduct(product model.Product) (uint64, error) {
	return pr.Repository.CreateProduct(product)
}
