package model

import (
  "time"
  "database/sql/driver"
)

type ProductStatus string

const (
  ProductStatusActive      ProductStatus = "active"
  ProductStatusInactive    ProductStatus = "inactive"
  ProductStatusOutOfStock  ProductStatus = "out_of_stock"
  ProductStatusDiscontinued ProductStatus = "discontinued"
)

func (ps *ProductStatus) Scan(value interface{}) error {
  *ps = ProductStatus(value.(string))
  return nil
}

func (ps ProductStatus) Value() (driver.Value, error) {
  return string(ps), nil
}

type Product struct {
  ID          string         `db:"id"`
  SKU         string         `db:"sku"`
  Name        string         `db:"name"`
  Description string         `db:"description"`
  Price       float64        `db:"price"`
  CostPrice   float64        `db:"cost_price"`
  Stock       int            `db:"stock"`
  ReservedStock int          `db:"reserved_stock"`
  CategoryID  string         `db:"category_id"`
  BrandID     *string        `db:"brand_id"`
  Weight      float64        `db:"weight"`
  Status      ProductStatus  `db:"status"`
  Images      []string       `db:"images"`
  Tags        []string       `db:"tags"`
  CreatedAt   time.Time      `db:"created_at"`
  UpdatedAt   time.Time      `db:"updated_at"`
  DeletedAt   *time.Time     `db:"deleted_at"`
}

type Category struct {
  ID          string     `db:"id"`
  Name        string     `db:"name"`
  Slug        string     `db:"slug"`
  Description string     `db:"description"`
  ParentID    *string    `db:"parent_id"`
  ImageURL    string     `db:"image_url"`
  SortOrder   int        `db:"sort_order"`
  IsActive    bool       `db:"is_active"`
  CreatedAt   time.Time  `db:"created_at"`
  UpdatedAt   time.Time  `db:"updated_at"`
  DeletedAt   *time.Time `db:"deleted_at"`
}

type Brand struct {
  ID          string     `db:"id"`
  Name        string     `db:"name"`
  Slug        string     `db:"slug"`
  Description string     `db:"description"`
  LogoURL     string     `db:"logo_url"`
  WebsiteURL  string     `db:"website_url"`
  CreatedAt   time.Time  `db:"created_at"`
  UpdatedAt   time.Time  `db:"updated_at"`
  DeletedAt   *time.Time `db:"deleted_at"`
}

type ProductVariant struct {
  ID         string     `db:"id"`
  ProductID  string     `db:"product_id"`
  SKU        string     `db:"sku"`
  Name       string     `db:"name"`
  Price      float64    `db:"price"`
  Stock      int        `db:"stock"`
  Attributes map[string]string `db:"attributes"`
  CreatedAt  time.Time  `db:"created_at"`
  UpdatedAt  time.Time  `db:"updated_at"`
}
