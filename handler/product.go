package handler

import (
  "context"
  "net/http"
  "encoding/json"

  "github.com/go-chi/chi/v5"
  "github.com/go-playground/validator/v10"

  "github.com/F-Dupraz/ecommerce-with-go/dto"
)

type ProductService interface {
    CreateProduct(ctx context.Context, prod *dto.CreateProductRequest) (*dto.CreateProductResponse, error)
    ListProducts(ctx context.Context, prods dto.ListProductsRequest) (*dto.ListProductsResponse, error)
    GetProductByID(ctx context.Context, prodID string) (*dto.ProductResponse, error)
    GetProductsByCategory(ctx context.Context, categoryID string, limit, offset int) (*dto.ListProductsResponse, error)
    GetRelatedProducts(ctx context.Context, prodID string, limit int) (*dto.ListProductsResponse, error)
    SearchProducts(ctx context.Context, query string, filters dto.SearchFilters) (*dto.ListProductsResponse, error)
    UpdateProduct(ctx context.Context, prodID string, prod *dto.UpdateProductRequest) (*dto.UpdateProductResponse, error)
    UpdateProductStock(ctx context.Context, prodID string, stock *dto.UpdateProductStockRequest) (*dto.UpdateProductStockResponse, error)
    DeleteProduct(ctx context.Context, prodID string) error
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
    router.Route("/products", func(r chi.Router) {
		r.Use(authMiddleware.Authenticate)

		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.RequireAuth)

			r.Get("/", p.GetProducts)
			r.Get("/{id}", p.GetProductByID)
			r.Get("/category/{id}", p.GetProductByCategory)
		})
	
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.RequireAuth)
			r.Use(authMiddleware.RequireAdmin)

			r.Post("/", p.CreateProduct)
			r.Put("/{id}", p.UpdateProduct)
			r.Put("/{id}/stock", p.UpdateProductStock)
			r.Delete("/{id}", p.DeleteProduct)
		})		
    })
}

func (p *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
    query := r.URL.Query()
    
    req := dto.ListProductsRequest{
        Limit:  20,
        Offset: 0,
    }
    
    if limitStr := query.Get("limit"); limitStr != "" {
        limit, err := strconv.Atoi(limitStr)
        if err != nil {
            p.respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid limit value '%s': must be an integer", limitStr), nil)
            return
        }
        req.Limit = limit
    }
    
    if offsetStr := query.Get("offset"); offsetStr != "" {
        offset, err := strconv.Atoi(offsetStr)
        if err != nil {
            p.respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid offset value '%s': must be an integer", offsetStr), nil)
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
            p.respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid min_price value '%s': must be a number", minPriceStr), nil)
            return
        }
        req.MinPrice = &minPrice
    }
    
    if maxPriceStr := query.Get("max_price"); maxPriceStr != "" {
        maxPrice, err := strconv.ParseFloat(maxPriceStr, 64)
        if err != nil {
            p.respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid max_price value '%s': must be a number", maxPriceStr), nil)
            return
        }
        req.MaxPrice = &maxPrice
    }
    
    if inStockStr := query.Get("in_stock"); inStockStr != "" {
        inStock, err := strconv.ParseBool(inStockStr)
        if err != nil {
            p.respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid in_stock value '%s': must be a boolean", inStockStr), nil)
            return
        }
        req.InStock = &inStock
    }
    
    if tags := query["tags"]; len(tags) > 0 {
        req.Tags = tags
    }
    
    if err := p.validator.Struct(req); err != nil {
        validationErrors := dto.FormatValidationErrors(err)
        p.respondWithError(w, http.StatusUnprocessableEntity, "Validation failed", validationErrors)
        return
    }
    
    response, err := p.productService.ListProducts(r.Context(), req)
    if err != nil {
        p.respondWithError(w, http.StatusInternalServerError, "Failed to get products", nil)
        return
    }
    
    p.respondWithSuccess(w, http.StatusOK, response)
}

func (p *ProductHandler) GetProductByCategory(w http.ResponseWriter, r *http.Request) {
    categoryID := chi.URLParam(r, "id")
    if categoryID == "" {
        p.respondWithError(w, http.StatusBadRequest, "Category ID is required", nil)
        return
    }
    
    limit := 20
    offset := 0
    
    if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
        parsedLimit, err := strconv.Atoi(limitStr)
        if err != nil {
            p.respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid limit value: %s", limitStr), nil)
            return
        }
        limit = parsedLimit
    }
    
    if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
        parsedOffset, err := strconv.Atoi(offsetStr)
        if err != nil {
            p.respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid offset value: %s", offsetStr), nil)
            return
        }
        offset = parsedOffset
    }
    
    response, err := p.productService.GetProductsByCategory(r.Context(), categoryID, limit, offset)
    if err != nil {
        p.respondWithError(w, http.StatusInternalServerError, "Failed to get products by category", nil)
        return
    }
    
    p.respondWithSuccess(w, http.StatusOK, response)
}

func (p *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
    prodID := chi.URLParam(r, "id")
    
    response, err := p.productService.GetProductByID(r.Context(), prodID)
    if err != nil {
        switch {
        case errors.Is(err, service.ErrInvalidID):
            p.respondWithError(w, http.StatusBadRequest, "Invalid product ID format", nil)
        case errors.Is(err, service.ErrProductNotFound):
            p.respondWithError(w, http.StatusNotFound, "Product not found", nil)
        default:
            p.respondWithError(w, http.StatusInternalServerError, "Internal server error", nil)
        }
        return
    }
    
	p.respondWithSuccess(w, http.StatusOK, response)
}

func (p *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
    user, ok := r.Context().Value("user").(AuthUser)
    if !ok || user.Role != "admin" {
        p.respondWithError(w, http.StatusForbidden, "Admin access required", nil)
        return
    }
    
    var req dto.CreateProductRequest
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        p.respondWithError(w, http.StatusBadRequest, "Cannot parse JSON: " + err.Error(), nil)
        return
    }
    
    if err := p.validator.Struct(req); err != nil {
        validationErrors := dto.FormatValidationErrors(err)
        p.respondWithError(w, http.StatusUnprocessableEntity, "Validation failed", validationErrors)
        return
    }
    
    response, err := p.productService.CreateProduct(r.Context(), &req)
    if err != nil {
        switch {
        case errors.Is(err, service.ErrInvalidPrice):
            p.respondWithError(w, http.StatusBadRequest, "Price must be greater than cost price", nil)
        case errors.Is(err, service.ErrDuplicateSKU):
            p.respondWithError(w, http.StatusConflict, "Product with this SKU already exists", nil)
        default:
            p.respondWithError(w, http.StatusInternalServerError, "Failed to create product", nil)
        }
        return
    }

    p.respondWithSuccess(w, http.StatusCreated, response)
}

func (p *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
  prodID := chi.URLParam(r, "id")

  var req dto.UpdateProductRequest

  if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
	p.respondWithError(w, http.StatusBadRequest, "Cannot parse JSON: " + err.Error(), nil)
	return
  }

  if err := p.validator.Struct(req); err != nil {
	validationErrors := dto.FormatValidationErrors(err)
	p.respondWithError(w, http.StatusUnprocessableEntity, "Validation failed!", validationErrors)
	return
  }

  response, err := p.productService.UpdateProductStock(r.Context(), prodID, &req)
  if err != nil {
    switch {
      case errors.Is(err, service.ErrInvalidID):
        p.respondWithError(w, http.StatusBadRequest, "Invalid product ID", nil)
      case errors.Is(err, service.ErrProductNotFound):
        p.respondWithError(w, http.StatusNotFound, "Product not found", nil)
      case errors.Is(err, service.ErrInsufficientStock):
        p.respondWithError(w, http.StatusBadRequest, "Insufficient stock available", nil)
      case errors.Is(err, service.ErrStockBelowReserved):
        p.respondWithError(w, http.StatusConflict, "Cannot reduce stock below reserved amount", nil)
      default:
        p.respondWithError(w, http.StatusInternalServerError, "Failed to update stock", nil)
      }
    return
  }

  p.respondWithSuccess(w, http.StatusOK, response)
}

func (p *ProductHandler) UpdateProductStock(w http.ResponseWriter, r *http.Request) {
    user, ok := r.Context().Value("user").(AuthUser)
    if !ok || user.Role != "admin" {
        p.respondWithError(w, http.StatusForbidden, "Admin access required", nil)
        return
    }
    
    prodID := chi.URLParam(r, "id")
    if prodID == "" {
        p.respondWithError(w, http.StatusBadRequest, "Product ID is required", nil)
        return
    }
    
    var req dto.UpdateProductStockRequest
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        p.respondWithError(w, http.StatusBadRequest, "Cannot parse JSON: " + err.Error(), nil)
        return
    }
    
    if err := p.validator.Struct(req); err != nil {
        validationErrors := dto.FormatValidationErrors(err)
        p.respondWithError(w, http.StatusUnprocessableEntity, "Validation failed", validationErrors)
        return
    }
    
    response, err := p.productService.UpdateProductStock(r.Context(), prodID, &req)
    if err != nil {
        p.respondWithError(w, http.StatusInternalServerError, "Failed to update product stock", nil)
        return
    }
    
    p.respondWithSuccess(w, http.StatusOK, response)
}

func (p *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
    user, ok := r.Context().Value("user").(AuthUser)
    if !ok || user.Role != "admin" {
        p.respondWithError(w, http.StatusForbidden, "Admin access required", nil)
        return
    }
    
    prodID := chi.URLParam(r, "id")
    if prodID == "" {
        p.respondWithError(w, http.StatusBadRequest, "Product ID is required", nil)
        return
    }
    
    err := p.productService.DeleteProduct(r.Context(), prodID)
    if err != nil {
        switch {
        case errors.Is(err, service.ErrInvalidID):
            p.respondWithError(w, http.StatusBadRequest, "Invalid product ID", nil)
        case errors.Is(err, service.ErrProductNotFound):
            p.respondWithError(w, http.StatusNotFound, "Product not found", nil)
        default:
            p.respondWithError(w, http.StatusInternalServerError, "Failed to delete product", nil)
        }
        return
    }

    p.respondWithSuccess(w, http.StatusNoContent, nil)
}
