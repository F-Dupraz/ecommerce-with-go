package dto

import (
  "time"

  "github.com/F-Dupraz/ecommerce-with-go/model"
)

// Requests

type CreateProductRequest struct {
  SKU         string   `json:"sku" validate:"required,min=3,max=50"`
  Name        string   `json:"name" validate:"required,min=3,max=200"`
  Description string   `json:"description" validate:"required,min=10,max=2000"`
  Price       float64  `json:"price" validate:"required,gt=0"`
  CostPrice   float64  `json:"cost_price" validate:"required,gte=0"`
  Stock       int      `json:"stock" validate:"required,gte=0"`
  CategoryID  string   `json:"category_id" validate:"required,uuid"`
  BrandID     *string  `json:"brand_id,omitempty" validate:"omitempty,uuid"`
  Weight      float64  `json:"weight" validate:"required,gt=0"` // grams
  Images      []string `json:"images" validate:"required,min=1,max=10,dive,url"`
  Tags        []string `json:"tags,omitempty" validate:"omitempty,max=20,dive,min=2,max=30"`
}

type UpdateProductRequest struct {
  Name        *string   `json:"name,omitempty" validate:"omitempty,min=3,max=200"`
  Description *string   `json:"description,omitempty" validate:"omitempty,min=10,max=2000"`
  Price       *float64  `json:"price,omitempty" validate:"omitempty,gt=0"`
  CostPrice   *float64  `json:"cost_price,omitempty" validate:"omitempty,gte=0"`
  CategoryID  *string   `json:"category_id,omitempty" validate:"omitempty,uuid"`
  BrandID     *string   `json:"brand_id,omitempty" validate:"omitempty,uuid"`
  Weight      *float64  `json:"weight,omitempty" validate:"omitempty,gt=0"`
  Images      []string  `json:"images,omitempty" validate:"omitempty,min=1,max=10,dive,url"`
  Tags        []string  `json:"tags,omitempty" validate:"omitempty,max=20,dive,min=2,max=30"`
  Status      *model.ProductStatus `json:"status,omitempty" validate:"omitempty,oneof=active inactive out_of_stock discontinued"`
}

type UpdateProductStockRequest struct {
  Stock     int  `json:"stock" validate:"required,gte=0"`
  Increment bool `json:"increment"`
}

type ReserveStockRequest struct {
  ProductID string `json:"product_id" validate:"required,uuid"`
  Quantity  int    `json:"quantity" validate:"required,gt=0"`
  OrderID   string `json:"order_id" validate:"required,uuid"`
  Duration  int    `json:"duration_minutes" validate:"required,min=5,max=60"`
}

type ListProductsRequest struct {
  Limit  int `query:"limit" validate:"omitempty,min=1,max=100"`
  Offset int `query:"offset" validate:"omitempty,gte=0"`
  CategoryID *string  `query:"category_id" validate:"omitempty,uuid"`
  BrandID    *string  `query:"brand_id" validate:"omitempty,uuid"`
  MinPrice   *float64 `query:"min_price" validate:"omitempty,gte=0"`
  MaxPrice   *float64 `query:"max_price" validate:"omitempty,gt=0"`
  InStock    *bool    `query:"in_stock"`
  Status     *model.ProductStatus `query:"status" validate:"omitempty,oneof=active inactive out_of_stock discontinued"`
  Tags       []string `query:"tags" validate:"omitempty,dive,min=2,max=30"`
  SortBy    string `query:"sort_by" validate:"omitempty,oneof=price name created_at stock popularity"`
  SortOrder string `query:"sort_order" validate:"omitempty,oneof=asc desc"`
}

type SearchProductsRequest struct {
  Query     string `query:"q" validate:"required,min=2,max=100"`
  Limit     int    `query:"limit" validate:"omitempty,min=1,max=50"`
  Offset    int    `query:"offset" validate:"omitempty,gte=0"`
  CategoryID *string `query:"category_id" validate:"omitempty,uuid"`
}

type GetProductByIDRequest struct {
  ID string `param:"id" validate:"required,uuid"`
}

type GetProductsByCategoryRequest struct {
  CategoryID       string `param:"category_id" validate:"required,uuid"`
  IncludeSubcategories bool `query:"include_subcategories"`
  Limit            int    `query:"limit" validate:"omitempty,min=1,max=100"`
  Offset           int    `query:"offset" validate:"omitempty,gte=0"`
}

type GetRelatedProductsRequest struct {
  ProductID string `param:"product_id" validate:"required,uuid"`
  Limit     int    `query:"limit" validate:"omitempty,min=1,max=20"`
}

type DeleteProductRequest struct {
  ID string `param:"id" validate:"required,uuid"`
}

// Responses

type ProductResponse struct {
  ID          string                `json:"id"`
  SKU         string                `json:"sku"`
  Name        string                `json:"name"`
  Description string                `json:"description"`
  Price       float64               `json:"price"`
  Stock       int                   `json:"stock"`
  Available   int                   `json:"available"`
  CategoryID  string                `json:"category_id"`
  Category    *CategoryResponse     `json:"category,omitempty"`
  BrandID     *string               `json:"brand_id,omitempty"`
  Brand       *BrandResponse        `json:"brand,omitempty"`
  Weight      float64               `json:"weight"`
  Status      model.ProductStatus   `json:"status"`
  Images      []string              `json:"images"`
  Tags        []string              `json:"tags"`
  CreatedAt   time.Time             `json:"created_at"`
  UpdatedAt   time.Time             `json:"updated_at"`
}

type CategoryResponse struct {
  ID          string    `json:"id"`
  Name        string    `json:"name"`
  Slug        string    `json:"slug"`
  Description string    `json:"description"`
  ParentID    *string   `json:"parent_id,omitempty"`
  ImageURL    string    `json:"image_url"`
}

type BrandResponse struct {
  ID         string `json:"id"`
  Name       string `json:"name"`
  Slug       string `json:"slug"`
  LogoURL    string `json:"logo_url"`
}

type CreateProductResponse struct {
  ID      string           `json:"id"`
  Product *ProductResponse `json:"product"`
  Message string           `json:"message"`
}

type UpdateProductResponse struct {
  Product *ProductResponse `json:"product"`
  Message string           `json:"message"`
}

type UpdateProductStockResponse struct {
  ProductID    string `json:"product_id"`
  NewStock     int    `json:"new_stock"`
  OldStock     int    `json:"old_stock"`
  Message      string `json:"message"`
}

type ReserveStockResponse struct {
  ReservationID string    `json:"reservation_id"`
  ProductID     string    `json:"product_id"`
  Quantity      int       `json:"quantity"`
  ExpiresAt     time.Time `json:"expires_at"`
  Message       string    `json:"message"`
}

type ListProductsResponse struct {
  Products   []ProductResponse `json:"products"`
  Total      int               `json:"total"`
  Limit      int               `json:"limit"`
  Offset     int               `json:"offset"`
  HasMore    bool              `json:"has_more"`
}

type SearchProductsResponse struct {
  Products   []ProductResponse `json:"products"`
  Query      string            `json:"query"`
  Total      int               `json:"total"`
  Limit      int               `json:"limit"`
  Offset     int               `json:"offset"`
}

type DeleteProductResponse struct {
  ID        string    `json:"id"`
  Message   string    `json:"message"`
  DeletedAt time.Time `json:"deleted_at"`
}
