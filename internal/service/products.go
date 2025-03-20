package service

import(
	"context"
	"berezovskyivalerii/csv-rest-app/internal/domain"
)

type ProductRepository interface {
	Create(ctx context.Context, product domain.Product) (int64, error)
	GetAll(ctx context.Context) ([]domain.Product, error)
	Update(ctx context.Context, id int64, product domain.ProductUpdate) (int64, error)
	Delete(ctx context.Context, id int64) (int64, error)
}

type Products struct{
	productRepository ProductRepository
}

func NewProducts(productRepository ProductRepository) *Products{
	return &Products{
		productRepository: productRepository,
	}
}

func (s* Products) Create(ctx context.Context, product domain.Product) (int64, error){
	return s.productRepository.Create(ctx, product)
}

func (s *Products) GetAll(ctx context.Context) ([]domain.Product, error){
	return s.productRepository.GetAll(ctx)
}

func (s* Products) Update(ctx context.Context, id int64, product domain.ProductUpdate) (int64, error){
	return s.productRepository.Update(ctx, id, product)
}

func (s* Products) Delete(ctx context.Context, id int64) (int64, error){
	return s.productRepository.Delete(ctx, id)
}
