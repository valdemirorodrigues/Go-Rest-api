package service

import (
	"go-api-meli/model"
	"go-api-meli/repository"
)

type CartService interface {
	AddProductToCart(cart model.Cart) error
	GetCartById(ID uint64) ([]model.CartResponse, error)
	Checkout(ID uint64) (model.Purchase, error)
	InsertTbProductTbcart(codeTbProduct uint64, codeTbCart uint64) (uint64, error)
}
type cartService struct {
	CartRepository    repository.CartRepository
	ProductRepository repository.ProductRepository
}

func NewCartService(cr repository.CartRepository, pr repository.ProductRepository) cartService {
	return cartService{
		CartRepository:    cr,
		ProductRepository: pr,
	}
}
func (cs cartService) AddProductToCart(cart model.Cart) error {
	var listProducts []model.Details
	for _, product := range cart.Products {
		pr, err := cs.ProductRepository.GetById(product.Id_Product)
		if err != nil {
			return err
		}
		listProducts = append(listProducts, model.Details{ID: product.ID, Id_Product: product.Id_Product, Quantity: product.Quantity, PriceFinal: pr.Price * float64(product.Quantity)})
	}
	err := cs.CartRepository.AddProductToCart(listProducts)
	if err != nil {
		return err
	}
	return nil
}

func (cs cartService) GetCartById(ID uint64) ([]model.CartResponse, error) {
	return cs.CartRepository.GetCartById(ID)
}
func (cs cartService) Checkout(ID uint64) (model.Purchase, error) {
	return cs.CartRepository.Checkout(ID)
}
func (cs cartService) InsertTbProductTbcart(codeTbProduct uint64, codeTbCart uint64) (uint64, error) {
	return cs.CartRepository.InsertTbProductTbcart(codeTbProduct, codeTbCart)
}
