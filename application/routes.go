package application

import (
  // "database/sql"

  "github.com/go-chi/chi/v5"
  "github.com/go-chi/chi/v5/middleware"
  "github.com/go-playground/validator/v10"

  "github.com/F-Dupraz/ecommerce-with-go/handler"
  "github.com/F-Dupraz/ecommerce-with-go/dto"
)

func loadRoutes() *chi.Mux {
  router := chi.NewRouter()

  router.Use(middleware.Logger)
  router.Use(middleware.Recoverer) // Agregá este, es importante
  router.Use(middleware.RequestID)  // Para trackear requests

  validator, err := dto.NewValidator()
  if err != nil {
    panic("Failed to initialize validator: " + err.Error())
  }

  router.Route("/api/v1", func(r chi.Router) {
    r.Route("/users", func(r chi.Router) {
      loadUsersRoutes(r, validator)
    })
  })

  return router
}

func loadUsersRoutes(router chi.Router, validator *validator.Validate) {
  // mockUserService := &MockUserService{}

  userHandler := handler.NewUserHandler(mockUserService, validator)

  userHandler.RegisterRoutes(router)
}

func loadOrderRoutes(router chi.Router) {
  orderHandler := &handler.Order{}

  router.Get("/", orderHandler.ListAllOrders)
  router.Get("/{id}", orderHandler.GetByID)
  router.Post("/", orderHandler.Create)
  router.Put("/{id}", orderHandler.UpdateByID)
  router.Delete("/{id}", orderHandler.DeleteByID)
}

func loadProductRoutes(router chi.Router) {
  productHandler := &handler.Product{}

  router.Get("/", productHandler.ListAllOrders)
  router.Get("/{id}", productHandler.GetByID)
  router.Post("/", productHandler.Create)
  router.Put("/{id}", productHandler.UpdateByID)
  router.Delete("/{id}", productHandler.DeleteByID)
}

