package application

import (
  "net/http"

  "github.com/go-chi/chi/v5"
  "github.com/go-chi/chi/v5/middleware"

  "github.com/F-Dupraz/ecommerce-with-go/handler"
)

func loadRoutes() *chi.Mux {
  router := chi.NewRouter()

  router.Use(middleware.Logger)
  
  router.Route("/orders", loadOrderRoutes)
  router.Route("/users", loadUsersRoutes)

  return router
}

func loadOrderRoutes(router chi.Router) {
  orderHandler := &handler.Order{}

  router.Get("/", orderHandler.ListAllOrders)
  router.Get("/{id}", orderHandler.GetByID)
  router.Post("/", orderHandler.Create)
  router.Put("/{id}", orderHandler.UpdateByID)
  router.Delete("/{id}", orderHandler.DeleteByID)
}

func loadUsersRoutes(router chi.Router) {
  userHandler := &handler.User{}

  router.Get("/", userHandler.ListAllUsers)
  router.Get("/{id}", userHandler.GetByID)
  router.Post("/", userHandler.Create)
  router.Put("/{id}", userHandler.UpdateByID)
  router.Delete("/{id}", userHandler.DeleteByID)
}
