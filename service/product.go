package service

import (
  "fmt"
  "time"
  "context"
  "errors"

  "github.com/F-Dupraz/ecommerce-with-go/dto"
  "github.com/F-Dupraz/ecommerce-with-go/model"
  "github.com/F-Dupraz/ecommerce-with-go/repository"

  "golang.org/x/crypto/bcrypt"
  "github.com/google/uuid"
)

type ProductRepository interface {}

type ProductService struct {
  repo ProductRepository
}

func NewProductService(repo ProductRepository) *ProductService {
  return &ProductService{
	repo: repo,
  }
}

func (s *ProductService) CreateProduct(ctx context.Context, prod *dto.CreateProductRequest) (*dto.CreateProductResponse, error) {
  // Checkear que el user sea admin

  // Hasheos/IDs

  // Pasar del dto al model

  // Pedirle al repo que cree el producto y pasarle un puntero al modelo

  // Devolver el dto del modelo actualizado
}

func (s *ProductService) ListProducts(ctx context.Context, prods dto.ListProductsRequest) (*dto.ListProductsResponse, error) {
  
}

func (s *ProductService) GetProductByID(ctx context.Context, prodID string) (*dto.ProductResponse, error) {
  
}

func (s *ProductService) GetProductsByCategory(ctx context.Context, categoryID string, limit, offset int) (*dto.ListProductsResponse, error) {
  
}

func (s *ProductService) GetRelatedProducts(ctx context.Context, prodID string, limit int) (*dto.ListProductsResponse, error) {
  
}

func (s *ProductService) SearchProducts(ctx context.Context, query string, filters dto.SearchFilters) (*dto.ListProductsResponse, error) {
  
}

func (s *ProductService) UpdateProduct(ctx context.Context, prodID string, prod *dto.UpdateProductRequest) (*dto.UpdateProductResponse, error) {
  
}

func (s *ProductService) UpdateProductStock(ctx context.Context, prodID string, stock *dto.UpdateProductStockRequest) (*dto.UpdateProductStockResponse, error) {
  
}

func (s *ProductService) DeleteProduct(ctx context.Context, prodID string) error {
  
}



