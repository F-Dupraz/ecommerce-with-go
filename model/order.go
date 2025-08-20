package model

import (
  "time"
  "database/sql/driver"
)

type OrderStatus string

const (
  OrderStatusPending     OrderStatus = "pending"
  OrderStatusPaid        OrderStatus = "paid"
  OrderStatusProcessing  OrderStatus = "processing"
  OrderStatusShipped     OrderStatus = "shipped"
  OrderStatusDelivered   OrderStatus = "delivered"
  OrderStatusCancelled   OrderStatus = "cancelled"
  OrderStatusRefunded    OrderStatus = "refunded"
  OrderStatusFailed      OrderStatus = "failed"
)

type PaymentMethod string

const (
  PaymentMethodCard        PaymentMethod = "card"
  PaymentMethodPayPal      PaymentMethod = "paypal"
  PaymentMethodTransfer    PaymentMethod = "transfer"
  PaymentMethodCash        PaymentMethod = "cash_on_delivery"
  PaymentMethodCrypto      PaymentMethod = "crypto"
)

type ShippingMethod string

const (
  ShippingMethodStandard  ShippingMethod = "standard"
  ShippingMethodExpress   ShippingMethod = "express"
  ShippingMethodOvernight ShippingMethod = "overnight"
  ShippingMethodPickup    ShippingMethod = "pickup"
)

func (os *OrderStatus) Scan(value interface{}) error {
  *os = OrderStatus(value.(string))
  return nil
}

func (os OrderStatus) Value() (driver.Value, error) {
  return string(os), nil
}

func (pm *PaymentMethod) Scan(value interface{}) error {
  *pm = PaymentMethod(value.(string))
  return nil
}

func (pm PaymentMethod) Value() (driver.Value, error) {
  return string(pm), nil
}

func (sm *ShippingMethod) Scan(value interface{}) error {
  *sm = ShippingMethod(value.(string))
  return nil
}

func (sm ShippingMethod) Value() (driver.Value, error) {
  return string(sm), nil
}

type Order struct {
  ID              string         `db:"id"`
  OrderNumber     string         `db:"order_number"`
  UserID          string         `db:"user_id"`
  Status          OrderStatus    `db:"status"`
  SubtotalAmount  int64          `db:"subtotal_amount"`
  TaxAmount       int64          `db:"tax_amount"`
  ShippingAmount  int64          `db:"shipping_amount"`
  TotalAmount     int64          `db:"total_amount"`
  PaymentMethod   PaymentMethod  `db:"payment_method"`
  PaymentID       *string        `db:"payment_id"`
  PaidAt          *time.Time     `db:"paid_at"`
  ShippingMethod  ShippingMethod `db:"shipping_method"`
  ShippingAddress string         `db:"shipping_address"`
  ShippingCity    string         `db:"shipping_city"`
  ShippingCountry string         `db:"shipping_country"`
  TrackingNumber  *string        `db:"tracking_number"`
  TrackingURL     *string        `db:"tracking_url"`
  EstimatedDelivery *time.Time   `db:"estimated_delivery"`
  DeliveredAt     *time.Time     `db:"delivered_at"`
  CreatedAt       time.Time      `db:"created_at"`
  UpdatedAt       time.Time      `db:"updated_at"`
  CancelledAt     *time.Time     `db:"cancelled_at"`
}

type OrderItem struct {
  ID              string         `db:"id"`
  OrderID         string         `db:"order_id"`
  ProductID       string         `db:"product_id"`
  VariantID       *string        `db:"variant_id"`
  ProductSKU      string         `db:"product_sku"`
  ProductName     string         `db:"product_name"`
  ProductImage    string         `db:"product_image"`
  UnitPrice       int64          `db:"unit_price"`
  Quantity        int            `db:"quantity"`
  SubtotalAmount  int64          `db:"subtotal_amount"`
  TotalAmount     int64          `db:"total_amount"`
  CreatedAt       time.Time      `db:"created_at"`
  UpdatedAt       time.Time      `db:"updated_at"`
}

type ShoppingCart struct {
  ID              string         `db:"id"`
  UserID          string         `db:"user_id"`
  ExpiresAt       time.Time      `db:"expires_at"`
  CreatedAt       time.Time      `db:"created_at"`
  UpdatedAt       time.Time      `db:"updated_at"`
}

type CartItem struct {
  ID              string         `db:"id"`
  CartID          string         `db:"cart_id"`
  ProductID       string         `db:"product_id"`
  VariantID       *string        `db:"variant_id"`
  Quantity        int            `db:"quantity"`
  ReservationID   *string        `db:"reservation_id"`
  AddedAt         time.Time      `db:"added_at"`
  UpdatedAt       time.Time      `db:"updated_at"`
}

