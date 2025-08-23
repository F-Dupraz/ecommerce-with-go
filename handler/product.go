package handler

import (
  "fmt"
  "net/http"
  "encoding/json"

  "github.com/go-chi/chi/v5"
  "github.com/go-playground/validator/v10"

  "github.com/F-Dupraz/ecommerce-with-go/dto"
)

type ProductService interface {
  CreatedProduct(ctx context.Context, prod dto.CreateProductResponse) (*dto.CreateProductResponse, error)
  ListProducts(ctx context.Context, prods dto.ListProductsRequest) (*dto.ListProductsResponse, error)
  GetProductById(ctx context.Context, prod_id dto.GetProductByIDRequest) (*dto.ProductResponse, error)
  GetProductsByCategory(ctx context.Context, cat_id dto.GetProductsByCategoryRequest) (*dto.ListProductsResponse, error)
  GetRelatedProducts(ctx context.Context, prod dto.GetRelatedProductsRequest) (*dto.ListProductsResponse, error)
  SearchProducts(ctx context.Context, query dto.SearchProductsRequest) (*dto.ListProductsResponse, error)
  UpdateProduct(ctx context.Context, prod_id string, prod dto.UpdateProductRequest) (*dto.UpdateProductResponse, error)
  UpdateProductStock(ctx context.Context, prod_id string, stock dto.UpdateProductStockRequest) (*dto.UpdateProductStockResponse, error)
  DeleteProduct(ctx context.Context, prod_id dto.DeleteProductRequest) (*dto.DeleteProductResponse, error)
}

func (p *Product) Create(w http.ResponseWriter, r *http.Request) {
  fmt.Println("Create a product")
}

func (p *Product) ListAllOrders(w http.ResponseWriter, r *http.Request) {
  fmt.Println("List all products")
}

func (p *Product) GetByID(w http.ResponseWriter, r *http.Request) {
  fmt.Println("Get a product by ID")
}

func (p *Product) UpdateByID(w http.ResponseWriter, r *http.Request) {
  fmt.Println("Update a product by ID")
}

func (p *Product) DeleteByID(w http.ResponseWriter, r *http.Request) {
  fmt.Println("Delete a product by ID")
}
