package handler

import (
  "fmt"
  "net/http"
  "encoding/json"

  "github.com/go-chi/chi/v5"
  "github.com/go-playground/validator/v10"

  "github.com/F-Dupraz/ecommerce-with-go/dto"
)

type OrderService interface {}

type OrderHandler struct {
  BaseHandler
  orderService OrderService
}

func NewOrderHandler(orderService OrderService, validator *validator.Validate) *OrderHandler {
  return &OrderHandler{
	orderService: orderService,
	BaseHandler: BaseHandler{validator: validator}
  }
}

func (o *OrderHandler) RegisterRoutes(router chi.Router) {
  router.Route("/orders", func (r chi.Router) {
	router.Get("/", o.GetOrders)
	router.Get("/{id}", o.GetOrderByID)
	router.Post("/", o.CreateOrder)
	router.Put("/{id}", o.UpdateOrderStatus)
	router.Delete("/{id}", o.DeleteOrder)
  })
}

func (o *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
  fmt.Println("Create an order")
}

func (o *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
  fmt.Println("List all orders")
}

func (o *OrderHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
  fmt.Println("Get and order by ID")
}

func (o *OrderHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
  fmt.Println("Update an order by ID")
}

func (o *OrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
  fmt.Println("Delete an order by ID")
}
