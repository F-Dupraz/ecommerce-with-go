package handler

import (
  "fmt"
  "net/http"
)

type Product struct {}

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
