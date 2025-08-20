package dto

import (
	"time"
	"github.com/F-Dupraz/ecommerce-with-go/model"
)

// Requests

type CreateOrderRequest struct {
  Items          []OrderItemInput      `json:"items" validate:"required,min=1,dive"`
  ShippingMethod  model.ShippingMethod `json:"shipping_method" validate:"required,oneof=standard express overnight pickup"`
  ShippingAddress string               `json:"shipping_address" validate:"required,min=10,max=200"`
  ShippingCity    string               `json:"shipping_city" validate:"required,min=2,max=100"`
  ShippingZip     string               `json:"shipping_zip" validate:"required,min=3,max=20"`
  ShippingCountry string               `json:"shipping_country" validate:"required,iso3166_1_alpha2"`
  PaymentMethod   model.PaymentMethod  `json:"payment_method" validate:"required,oneof=card paypal transfer cash_on_delivery crypto"`
}

type OrderItemInput struct {
  ProductID  string  `json:"product_id" validate:"required,uuid"`
  VariantID  *string `json:"variant_id,omitempty" validate:"omitempty,uuid"`
  Quantity   int     `json:"quantity" validate:"required,min=1,max=100"`
}

type UpdateOrderStatusRequest struct {
  Status        model.OrderStatus `json:"status" validate:"required,oneof=pending paid processing shipped delivered cancelled refunded failed"`
  InternalNotes string            `json:"internal_notes,omitempty" validate:"omitempty,max=500"`
  NotifyUser    bool              `json:"notify_user"`
}

type AddTrackingInfoRequest struct {
  TrackingNumber    string    `json:"tracking_number" validate:"required,min=5,max=100"`
  TrackingURL       string    `json:"tracking_url" validate:"required,url"`
  EstimatedDelivery time.Time `json:"estimated_delivery" validate:"required"`
  NotifyUser        bool      `json:"notify_user"`
}

type CancelOrderRequest struct {
  Reason        string `json:"reason" validate:"required,min=10,max=500"`
  RefundPayment bool   `json:"refund_payment"`
  RestockItems  bool   `json:"restock_items"`
}

type ListOrdersRequest struct {
  Limit  int `query:"limit" validate:"omitempty,min=1,max=100"`
  Offset int `query:"offset" validate:"omitempty,gte=0"`
  UserID      *string           `query:"user_id" validate:"omitempty,uuid"`
  Status      *model.OrderStatus `query:"status" validate:"omitempty,oneof=pending paid processing shipped delivered cancelled refunded failed"`
  DateFrom    *time.Time        `query:"date_from"`
  DateTo      *time.Time        `query:"date_to"`
  MinAmount   *int64            `query:"min_amount" validate:"omitempty,gte=0"`
  MaxAmount   *int64            `query:"max_amount" validate:"omitempty,gt=0"`
  OrderNumber *string           `query:"order_number" validate:"omitempty,min=1,max=50"`
  SortBy    string `query:"sort_by" validate:"omitempty,oneof=created_at total_amount status"`
  SortOrder string `query:"sort_order" validate:"omitempty,oneof=asc desc"`
}

type ListUserOrdersRequest struct {
  UserID    string             `param:"user_id" validate:"required,uuid"`
  Status    *model.OrderStatus `query:"status" validate:"omitempty,oneof=pending paid processing shipped delivered cancelled refunded failed"`
  Limit     int                `query:"limit" validate:"omitempty,min=1,max=50"`
  Offset    int                `query:"offset" validate:"omitempty,gte=0"`
}

type GetOrderByIDRequest struct {
  ID            string `param:"id" validate:"required,uuid"`
  IncludeItems  bool   `query:"include_items"`
}

type GetOrderByNumberRequest struct {
  OrderNumber   string `param:"order_number" validate:"required,min=1,max=50"`
  IncludeItems  bool   `query:"include_items"`
}

type GetOrderItemsRequest struct {
  OrderID string `param:"order_id" validate:"required,uuid"`
}

type ProcessPaymentRequest struct {
  OrderID       string              `json:"order_id" validate:"required,uuid"`
  PaymentMethod model.PaymentMethod `json:"payment_method" validate:"required,oneof=card paypal transfer crypto"`
  PaymentToken  string              `json:"payment_token" validate:"required"`
}

type AddToCartRequest struct {
  ProductID string  `json:"product_id" validate:"required,uuid"`
  VariantID *string `json:"variant_id,omitempty" validate:"omitempty,uuid"`
  Quantity  int     `json:"quantity" validate:"required,min=1,max=100"`
}

type UpdateCartItemRequest struct {
  Quantity int `json:"quantity" validate:"required,min=1,max=100"`
}

type RemoveFromCartRequest struct {
  ItemID string `param:"item_id" validate:"required,uuid"`
}

// Responses

type OrderResponse struct {
  ID              string               `json:"id"`
  OrderNumber     string               `json:"order_number"`
  UserID          string               `json:"user_id"`
  Status          model.OrderStatus    `json:"status"`
  SubtotalAmount  float64              `json:"subtotal_amount"`
  TaxAmount       float64              `json:"tax_amount"`
  ShippingAmount  float64              `json:"shipping_amount"`
  TotalAmount     float64              `json:"total_amount"`
  PaymentMethod   model.PaymentMethod  `json:"payment_method"`
  PaymentID       *string              `json:"payment_id,omitempty"`
  PaidAt          *time.Time           `json:"paid_at,omitempty"`
  ShippingMethod  model.ShippingMethod `json:"shipping_method"`
  ShippingAddress string               `json:"shipping_address"`
  ShippingCity    string               `json:"shipping_city"`
  ShippingCountry string               `json:"shipping_country"`
  TrackingNumber    *string            `json:"tracking_number,omitempty"`
  TrackingURL       *string            `json:"tracking_url,omitempty"`
  EstimatedDelivery *time.Time         `json:"estimated_delivery,omitempty"`
  DeliveredAt       *time.Time         `json:"delivered_at,omitempty"`
  Notes           string               `json:"notes,omitempty"`
  CouponCode      *string              `json:"coupon_code,omitempty"`
  Items           []OrderItemResponse  `json:"items,omitempty"`
  CreatedAt       time.Time            `json:"created_at"`
  UpdatedAt       time.Time            `json:"updated_at"`
  CancelledAt     *time.Time           `json:"cancelled_at,omitempty"`
}

type OrderItemResponse struct {
  ID             string  `json:"id"`
  OrderID        string  `json:"order_id"`
  ProductID      string  `json:"product_id"`
  VariantID      *string `json:"variant_id,omitempty"`
  ProductSKU     string  `json:"product_sku"`
  ProductName    string  `json:"product_name"`
  ProductImage   string  `json:"product_image"`
  UnitPrice      float64 `json:"unit_price"`
  Quantity       int     `json:"quantity"`
  SubtotalAmount float64 `json:"subtotal_amount"`
  TotalAmount    float64 `json:"total_amount"`
}

type CreateOrderResponse struct {
  ID          string         `json:"id"`
  OrderNumber string         `json:"order_number"`
  Order       *OrderResponse `json:"order"`
  Message     string         `json:"message"`
}

type UpdateOrderStatusResponse struct {
  Order      *OrderResponse    `json:"order"`
  OldStatus  model.OrderStatus `json:"old_status"`
  NewStatus  model.OrderStatus `json:"new_status"`
  Message    string            `json:"message"`
}

type CancelOrderResponse struct {
  OrderID      string    `json:"order_id"`
  RefundStatus string    `json:"refund_status,omitempty"`
  Message      string    `json:"message"`
  CancelledAt  time.Time `json:"cancelled_at"`
}

type ListOrdersResponse struct {
  Orders  []OrderResponse `json:"orders"`
  Total   int             `json:"total"`
  Limit   int             `json:"limit"`
  Offset  int             `json:"offset"`
  HasMore bool            `json:"has_more"`
}

type OrderItemsResponse struct {
  OrderID string              `json:"order_id"`
  Items   []OrderItemResponse `json:"items"`
  Total   int                 `json:"total"`
}

type ProcessPaymentResponse struct {
  OrderID       string            `json:"order_id"`
  PaymentID     string            `json:"payment_id"`
  Status        model.OrderStatus `json:"status"`
  Message       string            `json:"message"`
  ProcessedAt   time.Time         `json:"processed_at"`
}

// Cart responses
type CartResponse struct {
  ID         string             `json:"id"`
  UserID     string             `json:"user_id"`
  Items      []CartItemResponse `json:"items"`
  TotalItems int                `json:"total_items"`
  TotalAmount float64           `json:"total_amount"`
  ExpiresAt  time.Time          `json:"expires_at"`
  UpdatedAt  time.Time          `json:"updated_at"`
}

type CartItemResponse struct {
  ID           string  `json:"id"`
  ProductID    string  `json:"product_id"`
  ProductName  string  `json:"product_name"`
  ProductImage string  `json:"product_image"`
  UnitPrice    float64 `json:"unit_price"`
  Quantity     int     `json:"quantity"`
  Subtotal     float64 `json:"subtotal"`
  Available    bool    `json:"available"`
}

type AddToCartResponse struct {
  Cart    *CartResponse `json:"cart"`
  Message string        `json:"message"`
}

type OrderSummaryResponse struct {
  TotalOrders    int     `json:"total_orders"`
  TotalRevenue   float64 `json:"total_revenue"`
  AverageOrder   float64 `json:"average_order"`
  PendingOrders  int     `json:"pending_orders"`
  ProcessingOrders int   `json:"processing_orders"`
}
