package handler

import (
  "fmt"
  "net/http"
  "encoding/json"

  "github.com/go-chi/chi/v5"
  "github.com/go-playground/validator/v10"

  "github.com/F-Dupraz/ecommerce-with-go/dto"
)

type ProductService interface {
  CreateProduct(ctx context.Context, prod dto.CreateProductResponse) (*dto.CreateProductResponse, error)ListProducts(ctx context.Context, prods dto.ListProductsRequest) (*dto.ListProductsResponse, error)
  GetProducts(ctx context.Context, prods dto.ListProductsRequest) (*dto.ListProductsResponse, error)
  GetProductByID(ctx context.Context, prods dto.GetProductByIDRequest) (*dto.ProductResponse, error)
  GetProductsByCategory(ctx context.Context, cat_id dto.GetProductsByCategoryRequest) (*dto.ListProductsResponse, error)
  GetRelatedProducts(ctx context.Context, prod dto.GetRelatedProductsRequest) (*dto.ListProductsResponse, error)
  SearchProducts(ctx context.Context, query dto.SearchProductsRequest) (*dto.ListProductsResponse, error)
  UpdateProduct(ctx context.Context, prod_id string, prod dto.UpdateProductRequest) (*dto.UpdateProductResponse, error)
  UpdateProductStock(ctx context.Context, prod_id string, stock dto.UpdateProductStockRequest) (*dto.UpdateProductStockResponse, error)
  DeleteProduct(ctx context.Context, prod_id dto.DeleteProductRequest) (*dto.DeleteProductResponse, error)
}

type ProductHandler struct {
  BaseHandler
  productService ProductService
}

func NewProductHandler(productService ProductService, validator *validator.Validate) *ProductHandler {
  return &ProductHandler{
	productService: productService,
	BaseHandler: BaseHandler{validator: validator},
  }
}

func (p *ProductHandler) RegisterRoutes(router chi.Router) {
  router.Route("/products", func (r chi.Router) {
	router.Get("/", p.GetProducts)
	router.Get("/{id}", p.GetProductByID)
	router.Get("/category/{id}", p.GetProductByCategory)
	router.Post("/", p.CreateProduct)
    router.Put("/{id}", p.UpdateProduct)
    router.Put("/{id}/stock", p.UpdateProductStock)
	router.Delete("/{id}", p.DeleteProduct)
  })
}

func (p *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
  query := r.URL.Query()

  req := &dto.ListProductsRequest{}

  if limitStr := query.Get("limit"); limitStr != "" {
	  limit, err := strconv.Atoi(limitStr)
	  if err != nil {
		h.respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid limit value '%s': must be an integer", offsetStr), nil)
		return
	  }
	  req.Limit = limit
	} else {
	req.Limit = 20
  }

  if offsetStr := query.Get("offset"); offsetStr != "" {
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
	  h.respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid offset value '%s': must be an integer", offsetStr), nil)
	  return
	}
	req.Offset = offset
  }

  if categoryID := query.Get("category_id"); categoryID != "" {
	req.CategoryID = &categoryID
  }

  if minPriceStr := query.Get("min_price"); minPriceStr != "" {
	minPrice, err := strconv.ParseFloat(minPriceStr, 64)
	if err != nil {
	  h.respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid min_price value '%s': must be a number", offsetStr), nil)
	  return
	}
	req.MinPrice = &minPrice
  }

  if maxPriceStr := query.Get("max_price"); maxPriceStr != "" {
	maxPrice, err := strconv.ParseFloat(maxPriceStr, 64)
	if err != nil {
	  h.respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid max_price value '%s': must be a number", offsetStr), nil)
	  return
	}
	req.MaxPrice = &maxPrice
  }

  if inStockStr := query.Get("in_stock"); inStockStr != "" {
	inStock, err := strconv.ParseBool(inStockStr)
	if err != nil {
	  h.respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid in_stock value '%s': must be a boolean", offsetStr), nil)
	  return
	}
	req.InStock = &inStock
  }

  if tags := query["tags"]; len(tags) > 0 {
	req.Tags = tags
  }

  if err := p.validator.Struct(req); err != nil {
	validationErrors := dto.FormatValidationErrors(err)
	u.respondWithError(w, http.StatusUnprocessableEntity, "Validation failed!", validationErrors)
	return
  }

  response, err := p.productService.GetProducts(r.Context(), req)
  if err != nil {
	u.respondWithError(w, http.StatusInternalServerError, "Failed to get user by email", nil)
	return
  }

  u.respondWithSuccess(w, http.StatusOK, response)
}

func (p *ProductHandler) GetProductByCategory(w http.ResponseWriter, r *http.Request) {
  cat_id := chi.URLParam(r, "category")

  limit := r.URL.Query().Get("limit")
  offset := r.URL.Query().Get("offset")

  req := dto.GetProductsByCategoryRequest{
	ID: cat_id,
	Limit: limit,
	Offset: offset,
  }

  if err := p.validator.Struct(req); err != nil {
	validationErrors := dto.FormatValidationErrors(err)
	u.respondWithError(w, http.StatusUnprocessableEntity, "Validation failed!", validationErrors)
	return
  }

  response, err := p.productService.GetProductsByCategory(r.Context(), req)
  if err != nil {
	u.respondWithError(w, http.StatusInternalServerError, "Failed to get user by ID!", nil)
	return
  }

  u.respondWithSucces(w. http.StatusOK, response)
}

func (p *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
  prod_id := chi.URLParam(r, "id")

  req := dto.GetProductByIDRequest{
	ID: prod_id,
  }

  if err := p.validator.Struct(req); err != nil {
	validationErrors := dto.FormatValidationErrors(err)
	u.respondWithError(w, http.StatusUnprocessableEntity, "Validation failed!", validationErrors)
	return
  }

  response, err := p.productService.GetByID(r.Context(), req)
  if err != nil {
	u.respondWithError(w, http.StatusInternalServerError, "Failed to get user by ID!", nil)
	return
  }

  u.respondWithSucces(w. http.StatusOK, response)
}

func (p *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
  var req dto.CreateProductRequest

  if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
	u.respondWithError(w, http/StatusBadRequest, "Cannot parse JSON: " + err.Error(), nil)
	return
  }

  if err := p.validator.Struct(req); err != nil {
	validatorErrors := dto.FormatValidationErrors(err)
	u.respondWithError(w, http.StatusUnprocessableEntity, "Validation failed!", validationErrors)
	return
  }

  response, err := p.productService.CreateProduct(r.Context(), req)
  if err != nil {
	u.respondWithError(w, http.StatusInternalServerError, "Failed to create user", nil)
	return
  }

  u.respondWithSuccess(w, http.StatusCreated, response)
}

func (p *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
  prod_id := chi.URLParam(r, "id")

  var req dto.UpdateProductRequest

  if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
	u.respondWithError(w, http.StatusBadRequest, "Cannot parse JSON: " + err.Error(), nil)
	return
  }

  if err := p.validator.Struct(req); err != nil {
	validatorErrors := dto.FormatValidationErrors(err)
	u.respondWithError(w, http.StatusUnprocessableEntity, "Validation failed!", validationErrors)
	return
  }

  response, err := p.productService.UpdateProduct(r.Context(), prod_id, req)
  if err != nil {
	u.respondWithError(w, http.StatusInternalServerError, "Failed to create user", nil)
	return
  }

  u.respondWithSuccess(w, http.StatusOK, response)
}

func (p *ProductHandler) UpdateProductStock(w http.ResponseWriter, r *http.Request) {
  prod_id := chi.URLParam(r, "id")

  var req dto.UpdateProductStockRequest

  if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
	u.respondWithError(w, http.StatusBadRequest, "Cannot parse JSON: " + err.Error(), nil)
	return
  }

  if err := p.validator.Struct(req); err != nil {
	validatorErrors := dto.FormatValidatonErrors(err)
	u.respondWithErrors(w, http.StatusUnprocesableEntity, "Validation failed!", validationErrors)
	return
  }

  response, err := p.productService.UpdateProductStock(r.Context(), prod_id, req)
  if err != nil {
	u.respondWithError(w, http.StatusInternalServerError, "Failed to update the stock of the product", nil)
    return
  }

  u.respondWithSucces(w, http.StatusOK, response)
}

func (p *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	prod_id := chi.URLParam(r, "id")

	var req dto.DeleteProductRequest{
	  ID: prod_id,
	}

  if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
	u.respondWithError(w, http/StatusBadRequest, "Cannot parse JSON: " + err.Error(), nil)
	return
  }

  if err := p.validator.Struct(req); err != nil {
	validatorErrors := dto.FormatValidationErrors(err)
	u.respondWithError(w, http.StatusUnprocessableEntity, "Validation failed!", validationErrors)
	return
  }

  response, err := p.productService.DeleteProduct(r.Context(), req)
  if err != nil {
	u.respondWithError(w, http.StatusInternalServerError, "Failed to create user", nil)
	return
  }

  u.respondWithSuccess(w, http.StatusOK, response)
}
