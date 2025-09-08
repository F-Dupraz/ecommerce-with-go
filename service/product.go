package service

import (
  "fmt"
  "errors"
  "context"

  "github.com/F-Dupraz/ecommerce-with-go/dto"
  "github.com/F-Dupraz/ecommerce-with-go/model"
  "github.com/F-Dupraz/ecommerce-with-go/repository"

  "github.com/google/uuid"
)

var (
    ErrProductNotFound = errors.New("product not found")
    ErrInvalidID = errors.New("invalid ID format")
    ErrInsufficientStock = errors.New("insufficient stock")
    ErrDuplicateSKU = errors.New("SKU already exists")
    ErrInvalidPrice = errors.New("price must be greater than cost")
    ErrStockBelowReserved = errors.New("cannot reduce stock below reserved amount")
)

type ProductService struct {
  repo repository.ProductRepository
}

func NewProductService(repo ProductRepository) *ProductService {
  return &ProductService{
	repo: repo,
  }
}

func (s *ProductService) CreateProduct(ctx context.Context, prod *dto.CreateProductRequest) (*dto.CreateProductResponse, error) {
  if req.Price <= req.CostPrice {
	return nil, ErrInvalidPrice
  }

  existing, _ := s.repo.GetBySKU(ctx, req.SKU)
  if existing != nil {
	return nil, ErrDuplicateSKU
  }

  productID := uuid.New().String()

  var newProduct model.Product{
	ID: productID,
	SKU: prod.SKU,
	Name: prod.Name,
	Description: prod.Description,
	Price: prod.Price,
	CostPrice: prod.CostPrice,
	Stock: prod.Stock,
	CategoryID: prod.CategoryID,
	BrandID: prod.BrandID,
	Weight: prod.Weight,
	Images: prod.Images,
	Tags: prod.Tags,
  }

  if err := s.repo.CreateProduct(ctx, &newProduct); err != nil {
	  return nil, fmt.Errorf("Failed to create product: %w", err)
  }

  return &dto.CreateProductResponse{
	ID: product.ID,
	Product: s.toProductResponse(product),
	Message: "Product created successfully",
  }, nil
}

func (s *ProductService) ListProducts(ctx context.Context, prods dto.ListProductsRequest) (*dto.ListProductsResponse, error) {
  
}

func (s *ProductService) GetProductByID(ctx context.Context, prodID string) (*dto.ProductResponse, error) {
  if err := uuid.Parse(prodID); err != nil {
	return nil, fmt.Errorf("Invalid ID")
  }

  product, err := s.repo.GetProductByID(ctx, prodID)
  if err != nil {
	return nil, fmt.Errorf("Error getting the product")
  }

  return &dto.ProductResponse{
	ID: product.ID,
	SKU: product.SKU,
	Name: product.Name,
	Description: product.Description,
	Price: product.Price,
	CostPrice: product.CostPrice,
	Stock: product.Stock,
	CategoryID: product.Category,
	BrandID: product.BrandID,
	Weight: product.Weight,
	Images: product.Images,
	Tags: product.Tags,
	CreatedAt: product.CreatedAt,
	UpdatedAt: product.UpdatedAt,
  }, nil
}

func (s *ProductService) GetProductsByCategory(ctx context.Context, req *dto.GetProductsByCategoryRequest) (*dto.ListProductsResponse, error) {
  catID := req.CategoryID
  includeSubcatefories := req.IncludeSubcategories
  limit := req.Limit
  offset := req.Offset

  if err := uuid.Parse(catID); err != nil {
	return nil, fmt.Errorf("Invalid category ID")
  }

  productsByCategories, err := s.repo.GetProductsByCategory(ctx, catID, includeSubcatefories, limit, offset)
  if err != nil {
	return nil, fmt.Errorf("Error getting the products")
  }

  return &dto.ListProductsResponse{
	Products: productsByCategories,
	Limit: limit,
	Offset: offset,
  }, nil
}

// func (s *ProductService) GetRelatedProducts(ctx context.Context, prodID string, limit int) (*dto.ListProductsResponse, error) {
//   TODO: I have no idea how to do this!
// }

func (s *ProductService) SearchProducts(ctx context.Context, req dto.SearchProductsRequest) (*dto.ListProductsResponse, error) {
  query := req.Query
  limit := req.Limit
  offset := req.Offset
  catID := req.CategoryID

  if catID != "" {
	if err := uuid.Parse(catID); err != nil {
	  return nil, fmt.Errorf("Invalid category ID")
	}
  }

  searchProductsResponse, err := s.repo.SearchProducts(ctx, query, limit, offset, catID)
  if err != nil {
	return nil, fmt.Errorf("Error getting the products")
  }

  return &dto.SearchProductsResponse{
	Products: searchProductsResponsem
	Query: query,
	Limit: limit,
	Offset: offset,
  }, nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, prodID string, prod *dto.UpdateProductRequest) (*dto.UpdateProductResponse, error) {
  if err := uuid.Parse(prodID); err != nil {
	return nil, fmt.Errorf("Invalid ID")
  }

  var currentProd model.Product

  currentProd, err := s.repo.GetProductByID(ctx, prodID)
  if err != nil {
	return nil, fmt.Errorf("Invalid ID")
  }

  updates := make(map[string]interface{})
  if prod.Name != nil {
	updates["name"] = *prod.Name
  }
  if prod.Description != nil {
	updates["description"] = *prod.Description
  }
  if prod.Price != nil {
	updates["price"] = *prod.Price
  }
  if prod.CostPrice != nil {
	updates["cost_price"] = *prod.CostPrice
  }
  if prod.CategoryID != nil {
	updates["category_id"] = *prod.CategoryID
  }
  if prod.BrandID != nil {
	updates["brand_id"] = *prod.BrandID
  }
  if prod.Weight != nil {
	updates["weight"] = *prod.Weight
  }
  if prod.Images != nil {
	updates["images"] = *prod.Images
  }
  if prod.Tags != nil {
	updates["tags"] = *prod.Tags
  }
  if prod.Status != nil {
	updates["status"] = *prod.Status
  }

  updatedProduct, err := s.repo.Update(ctx, prodID, updates)
  if err != nil {
	return nil, fmt.Errorf("failed to update user: %w", err)
  }

  var product dto.ProductResponse{
	ID: updatedProduct.ID,
	SKU: updatedProduct.SKU,
	Name: updatedProduct.Name,
	Description: updatedProduct.Description,
	Price: updatedProduct.Price,
	CostPrice: updatedProduct.CostPrice,
	Stock: updatedProduct.Stock,
	CategoryID: updatedProduct.Category,
	BrandID: updatedProduct.BrandID,
	Weight: updatedProduct.Weight,
	Images: updatedProduct.Images,
	Tags: updatedProduct.Tags,
	CreatedAt: updatedProduct.CreatedAt,
	UpdatedAt: updatedProduct.UpdatedAt,
  }

  return &dto.UpdateProductResponse{
	Product: &product,
	Message: "User updated successfully!",
  }, nil
}

func (s *ProductService) UpdateProductStock(ctx context.Context, prodID string, req *dto.UpdateProductStockRequest) (*dto.UpdateProductStockResponse, error) {
    _, err := uuid.Parse(prodID)
    if err != nil {
        return nil, ErrInvalidID
    }
    
    product, err := s.repo.GetByID(ctx, prodID)
    if err != nil {
        return nil, ErrProductNotFound
    }
    
    delta := req.Stock
    if !req.Increment {
        delta = req.Stock - product.Stock
    }
    
    newStock := product.Stock + delta
    if newStock < product.ReservedStock {
        return nil, ErrStockBelowReserved
    }
    
    if newStock < 0 {
        return nil, ErrInsufficientStock
    }
    
    if err := s.repo.UpdateStock(ctx, prodID, delta); err != nil {
        return nil, fmt.Errorf("failed to update stock: %w", err)
    }
    
    return &dto.UpdateProductStockResponse{
        ProductID: prodID,
        NewStock:  newStock,
        Message:   "Stock updated successfully",
    }, nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, prodID string) error {
  if err := uuid.Parse(prodID); err != nil {
	return nil, fmt.Errorf("Invalid ID")
  }

  deletedProduct, err := s.repo.DeleteProduct(ctx, prodID)
  if err != nil {
	return nil, fmt.Errorf("failed to update user: %w", err)
  }

  return &dto.DeleProductResponse{
	ID: deletedProduct.ID,
	Message: "Product deleted successfully!",
	DeletedAt: deletedProduct.DeletedAt,
  }, nil
}

func (s *ProductService) toProductResponse(p *model.Product) dto.ProductResponse {
    return dto.ProductResponse{
        ID:          p.ID,
        SKU:         p.SKU,
        Name:        p.Name,
        Description: p.Description,
        Price:       p.Price,
        Stock:       p.Stock,
        Available:   p.Stock - p.ReservedStock,
        Images:      p.Images,
        Tags:        p.Tags,
        CreatedAt:   p.CreatedAt,
        UpdatedAt:   p.UpdatedAt,
    }
}

func (s *ProductService) toProductResponses(products []*model.Product) []dto.ProductResponse {
    responses := make([]dto.ProductResponse, len(products))
    for i, p := range products {
        responses[i] = s.toProductResponse(p)
    }
    return responses
}

