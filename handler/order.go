package handler

import (
  "fmt"
  "net/http"
  "encoding/json"

  "github.com/go-chi/chi/v5"
  "github.com/go-playground/validator/v10"

  "github.com/F-Dupraz/ecommerce-with-go/dto"
)

type Order struct {}

func (o *Order) Create(w http.ResponseWriter, r *http.Request) {
  fmt.Println("Create an order")
}

func (o *Order) ListAllOrders(w http.ResponseWriter, r *http.Request) {
  fmt.Println("List all orders")
}

func (o *Order) GetByID(w http.ResponseWriter, r *http.Request) {
  fmt.Println("Get and order by ID")
}

func (o *Order) UpdateByID(w http.ResponseWriter, r *http.Request) {
  fmt.Println("Update an order by ID")
}

func (o *Order) DeleteByID(w http.ResponseWriter, r *http.Request) {
  fmt.Println("Delete an order by ID")
}
