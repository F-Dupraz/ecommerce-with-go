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

type OrderRepository interface {}

type OrderService struct {
  repo OrderRepository
}

func NewOrderService(repo OrderRepository) *OrderService {
  return &OrderService{
	repo: repo,
  }
}

func (s *OrderService) CreateOrder(ctx context.Context, req *dto.CreateOrderRequest) (*dto.CreateOrderResponse, error) {
  
}

func (s *OrderService) ListOrders(ctx context.Context, req dto.ListOrdersRequest) (*dto.ListOrdersResponse, error) {
  
}

func (s *OrderService) GetOrderByID(ctx context.Context, orderID string, includeItems bool) (*dto.OrderResponse, error) {
  
}

func (s *OrderService) UpdateOrderStatus(ctx context.Context, orderID string, req *dto.UpdateOrderStatusRequest) (*dto.UpdateOrderStatusResponse, error) {
  
}

func (s *OrderService) CancelOrder(ctx context.Context, orderID string, req *dto.CancelOrderRequest) (*dto.CancelOrderResponse, error) {
  
}



