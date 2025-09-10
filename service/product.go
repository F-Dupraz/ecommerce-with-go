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
	ErrCategoryNotFound = errors.New("category not found")
    ErrInvalidParams = errors.New("invalid parameters")
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
  if prod.Price <= prod.CostPrice {
	return nil, ErrInvalidPrice
  }

  existing, _ := s.repo.GetBySKU(ctx, prod.SKU)
  if existing != nil {
	return nil, ErrDuplicateSKU
  }

  productID := uuid.New().String()

  newProduct := model.Product{
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
	ID: newProduct.ID,
	Product: s.toProductResponse(product),
	Message: "Product created successfully",
  }, nil
}

func (s *ProductService) ListProducts(ctx context.Context, prods dto.ListProductsRequest) (*dto.ListProductsResponse, error) {
  limit := prods.Limit
  offset := prods.Offset
  sort := prods.SortBy
  if sort == "" {
	sort = "created_at"
  }
  sort_order := prods.SortOrder
  if sort_order == "" {
	sort_order = "asc"
  }

  params := make(map[string]interface{})
  if prods.CategoryID != nil {
	if err := uuid.Parse(prods.CategoryID); err != nil {
	  return nil, ErrInvalidID
	}
	  params["category_id"] = *prods.CategoryID
  }
  if prods.BrandID != nil {
	if err := uuid.Parse(prods.BrandID); err != nil {
	  return nil, ErrInvalidID
	}
	  params["brand_id"] = *prods.BrandID
  }
  if prods.MinPrice != nil {
	params["min_price"] = *prods.MinPrice
  }
  if prods.MaxPrice != nil {
	params["max_price"] = *prods.MaxPrice
  }
  if prods.InStock != nil {
	params["in_stock"] = *prods.InStock
  }
  if prods.Status != nil {
	params["status"] = *prods.Status
  }
  
  products, err := s.repo.ListProducts(ctx, limit, offset, sort, sort_order, params)
  if err != nil {
    return nil, fmt.Errorf("failed to list products: %w", err)
  }

  return &dto.ListProductsResponse{
	Products: s.toProductResponses(products),
	Limit: limit,
	Offset: offset,
  }, nil
}

func (s *ProductService) GetProductByID(ctx context.Context, prodID string) (*dto.ProductResponse, error) {
  if err := uuid.Parse(prodID); err != nil {
	return nil, ErrInvalidID
  }

  product, err := s.repo.GetProductByID(ctx, prodID)
  if err != nil {
    return nil, fmt.Errorf("failed to get product: %w", err)
  }

  response := s.toProductResponse(product)
  return &response, nil
}

func (s *ProductService) GetProductsByCategory(ctx context.Context, prod*dto.GetProductsByCategoryRequest) (*dto.ListProductsResponse, error) {
  catID := prodCategoryID
  includeSubcatefories := prodIncludeSubcategories
  limit := prodLimit
  offset := prodOffset

  if err := uuid.Parse(catID); err != nil {
	return nil, ErrInvalidID
  }

  products, err := s.repo.GetProductsByCategory(ctx, catID, includeSubcatefories, limit, offset)
  if err != nil {
    return nil, fmt.Errorf("failed to get products by category: %w", err)
  }

  return &dto.ListProductsResponse{
	Products: s.toProductResponses(products),
	Limit: limit,
	Offset: offset,
  }, nil
}

// func (s *ProductService) GetRelatedProducts(ctx context.Context, prodID string, limit int) (*dto.ListProductsResponse, error) {
//   TODO: I have no idea how to do this!
// }

func (s *ProductService) SearchProducts(ctx context.Context, proddto.SearchProductsRequest) (*dto.SearchProductsResponse, error) {
  query := prodQuery
  limit := prodLimit
  offset := prodOffset
  catID := prodCategoryID

  if catID != "" {
	if err := uuid.Parse(catID); err != nil {
	  return nil, ErrInvalidID
	}
  }

  products, err := s.repo.SearchProducts(ctx, query, limit, offset, catID)
  if err != nil {
    return nil, fmt.Errorf("failed to search products: %w", err)
  }

  return &dto.SearchProductsResponse{
	Products: s.toProductResponses(products),
	Query: query,
	Limit: limit,
	Offset: offset,
  }, nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, prodID string, prod *dto.UpdateProductRequest) (*dto.UpdateProductResponse, error) {
  if err := uuid.Parse(prodID); err != nil {
	return nil, ErrInvalidID
  }

  var currentProd model.Product

  currentProd, err := s.repo.GetProductByID(ctx, prodID)
  if err != nil {
	return nil, ErrProductNotFound
  }

  if req.Price != nil && req.CostPrice != nil {
    if *req.Price <= *req.CostPrice {
      return nil, ErrInvalidPrice
    }
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
    return nil, fmt.Errorf("failed to update product: %w", err)
  }

  response := s.toProductResponse(updatedProduct)
  return &dto.UpdateProductResponse{
      Product: &response,
      Message: "Product updated successfully",
  }, nil
}

func (s *ProductService) UpdateProductStock(ctx context.Context, prodID string, prod*dto.UpdateProductStockRequest) (*dto.UpdateProductStockResponse, error) {
    if _, err := uuid.Parse(prodID); err != nil {
        return nil, ErrInvalidID
    }
    
    product, err := s.repo.GetByID(ctx, prodID)
    if err != nil {
        if errors.Is(err, repository.ErrNotFound) {
            return nil, ErrProductNotFound
        }
        return nil, fmt.Errorf("failed to get product: %w", err)
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
    if _, err := uuid.Parse(prodID); err != nil {
        return ErrInvalidID
    }

    err := s.repo.DeleteProduct(ctx, prodID)
    if err != nil {
        if errors.Is(err, repository.ErrNotFound) {
            return ErrProductNotFound
        }
        return fmt.Errorf("failed to delete product: %w", err)
    }

    return nil
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

