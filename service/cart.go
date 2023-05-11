package service

import (
	"database/sql"
	"go-api-meli/model"
	"go-api-meli/repository"
)

type products struct {
	db *sql.DB
}

type CartService interface {
	AddProductToCart(cart model.Cart) (uint64, error)
	GetCartById(ID uint64) (model.Detail, error)
	CartFinallity(ID uint64) (model.Detail, error)
}
type cartService struct {
	Repository repository.CartRepository
}

func NewCartService(repo repository.CartRepository) cartService {
	return cartService{
		Repository: repo,
	}
}
func (cs cartService) AddProductToCart(cart model.Cart) (uint64, error) {
	return cs.Repository.AddProductToCart(cart)
}
func (cs cartService) GetCartById(ID uint64) (model.Detail, error) {
	return cs.Repository.GetCartById(ID)
}
func (cs cartService) CartFinallity(ID uint64) (model.Detail, error) {
	return cs.Repository.CartFinallity(ID)

}
