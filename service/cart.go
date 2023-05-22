package service

import (
	"go-api-meli/model"
	"go-api-meli/repository"
)

type CartService interface {
	AddProductToCart(cart model.Cart) (model.Cart, error)
	GetCartById(ID uint64) ([]model.Detail, error)
	CartFinallity(ID uint64) (model.Purchase, error)
	InsertTbProductTbcart(codeTbProduct uint64, codeTbCart uint64) (uint64, error)
}
type cartService struct {
	Repository repository.CartRepository
}

func NewCartService(repo repository.CartRepository) cartService {
	return cartService{
		Repository: repo,
	}
}
func (cs cartService) AddProductToCart(cart model.Cart) (model.Cart, error) {
	return cs.Repository.AddProductToCart(cart)
}
func (cs cartService) GetCartById(ID uint64) ([]model.Detail, error) {
	return cs.Repository.GetCartById(ID)
}
func (cs cartService) CartFinallity(ID uint64) (model.Purchase, error) {
	return cs.Repository.CartFinallity(ID)
}
func (cs cartService) InsertTbProductTbcart(codeTbProduct uint64, codeTbCart uint64) (uint64, error) {
	return cs.Repository.InsertTbProductTbcart(codeTbProduct, codeTbCart)
}
