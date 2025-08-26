package handler

import (
  "fmt"
  "time"
  "errors"
  "strconv"
  "context"
  "net/http"
  "encoding/json"

  "github.com/go-chi/chi/v5"
  "github.com/go-playground/validator/v10"

  "github.com/F-Dupraz/ecommerce-with-go/dto"
  "github.com/F-Dupraz/ecommerce-with-go/model"
)

type OrderService interface {
    CreateOrder(ctx context.Context, req *dto.CreateOrderRequest) (*dto.CreateOrderResponse, error)
    ListOrders(ctx context.Context, req dto.ListOrdersRequest) (*dto.ListOrdersResponse, error)
    GetOrderByID(ctx context.Context, orderID string, includeItems bool) (*dto.OrderResponse, error)
    UpdateOrderStatus(ctx context.Context, orderID string, req *dto.UpdateOrderStatusRequest) (*dto.UpdateOrderStatusResponse, error)
    CancelOrder(ctx context.Context, orderID string, req *dto.CancelOrderRequest) (*dto.CancelOrderResponse, error)
}

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
    router.Route("/orders", func(r chi.Router) {
        r.Get("/", o.GetOrders)
        r.Get("/{id}", o.GetOrderByID)
        r.Post("/", o.CreateOrder)
        r.Put("/{id}", o.UpdateOrderStatus)
        r.Delete("/{id}", o.DeleteOrder)
    })
}

func (o *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
    user, ok := r.Context().Value("user").(AuthUser)
    if !ok {
        o.respondWithError(w, http.StatusUnauthorized, "Authentication required", nil)
        return
    }
    
    var req dto.CreateOrderRequest
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        o.respondWithError(w, http.StatusBadRequest, "Invalid JSON format: " + err.Error(), nil)
        return
    }
    
    if err := o.validator.Struct(req); err != nil {
        validationErrors := dto.FormatValidationErrors(err)
        o.respondWithError(w, http.StatusUnprocessableEntity, "Validation failed", validationErrors)
        return
    }
    
    response, err := o.orderService.CreateOrder(r.Context(), &req)
    if err != nil {
        switch {
        case errors.Is(err, ErrInsufficientStock):
            o.respondWithError(w, http.StatusConflict, "Insufficient stock for one or more items", nil)
        case errors.Is(err, ErrProductNotFound):
            o.respondWithError(w, http.StatusNotFound, "One or more products not found", nil)
        case errors.Is(err, ErrInvalidShippingAddress):
            o.respondWithError(w, http.StatusBadRequest, "Invalid shipping address", nil)
        default:
            o.respondWithError(w, http.StatusInternalServerError, "Failed to create order", nil)
        }
        return
    }
    
    o.respondWithSuccess(w, http.StatusCreated, response)
}

func (o *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
    user, ok := r.Context().Value("user").(AuthUser)
    if !ok {
        o.respondWithError(w, http.StatusUnauthorized, "Authentication required", nil)
        return
    }

    req := dto.ListOrdersRequest{
        Limit:     10,
        Offset:    0,
        SortBy:    "created_at",
        SortOrder: "desc",
    }

    query := r.URL.Query()
    
    if limitStr := query.Get("limit"); limitStr != "" {
        limit, err := strconv.Atoi(limitStr)
        if err != nil {
            o.respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid limit: %s", limitStr), nil)
            return
        }
        req.Limit = limit
    }

    if offsetStr := query.Get("offset"); offsetStr != "" {
        offset, err := strconv.Atoi(offsetStr)
        if err != nil {
            o.respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid offset: %s", offsetStr), nil)
            return
        }
        req.Offset = offset
    }

    if statusStr := query.Get("status"); statusStr != "" {
        status := model.OrderStatus(statusStr)
        req.Status = &status
    }

    if dateFromStr := query.Get("date_from"); dateFromStr != "" {
        dateFrom, err := time.Parse(time.RFC3339, dateFromStr)
        if err != nil {
            o.respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid date_from format (use RFC3339): %s", dateFromStr), nil)
            return
        }
        req.DateFrom = &dateFrom
    }

    if dateToStr := query.Get("date_to"); dateToStr != "" {
        dateTo, err := time.Parse(time.RFC3339, dateToStr)
        if err != nil {
            o.respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid date_to format (use RFC3339): %s", dateToStr), nil)
            return
        }
        req.DateTo = &dateTo
    }

    if minAmountStr := query.Get("min_amount"); minAmountStr != "" {
        minAmount, err := strconv.ParseInt(minAmountStr, 10, 64)
        if err != nil {
            o.respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid min_amount: %s", minAmountStr), nil)
            return
        }
        req.MinAmount = &minAmount
    }

    if maxAmountStr := query.Get("max_amount"); maxAmountStr != "" {
        maxAmount, err := strconv.ParseInt(maxAmountStr, 10, 64)
        if err != nil {
            o.respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid max_amount: %s", maxAmountStr), nil)
            return
        }
        req.MaxAmount = &maxAmount
    }

    if user.Role == "admin" && query.Get("user_id") != "" {
        userID := query.Get("user_id")
        req.UserID = &userID
    }

    if orderNumber := query.Get("order_number"); orderNumber != "" {
        req.OrderNumber = &orderNumber
    }

    if sortBy := query.Get("sort_by"); sortBy != "" {
        req.SortBy = sortBy
    }

    if sortOrder := query.Get("sort_order"); sortOrder != "" {
        req.SortOrder = sortOrder
    }

    if err := o.validator.Struct(req); err != nil {
        validationErrors := dto.FormatValidationErrors(err)
        o.respondWithError(w, http.StatusUnprocessableEntity, "Validation failed", validationErrors)
        return
    }

    var orders *dto.ListOrdersResponse
    var err error

    orders, err = o.orderService.ListOrders(r.Context(), req)

    if err != nil {
        o.respondWithError(w, http.StatusInternalServerError, "Failed to get orders", nil)
        return
    }

    o.respondWithSuccess(w, http.StatusOK, orders)
}

func (o *OrderHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
    user, ok := r.Context().Value("user").(AuthUser)
    if !ok {
        o.respondWithError(w, http.StatusUnauthorized, "Authentication required", nil)
        return
    }

    orderID := chi.URLParam(r, "id")
    if orderID == "" {
        o.respondWithError(w, http.StatusBadRequest, "Order ID is required", nil)
        return
    }

    includeItems := false
    if includeItemsStr := r.URL.Query().Get("include_items"); includeItemsStr != "" {
        parsedBool, err := strconv.ParseBool(includeItemsStr)
        if err != nil {
            o.respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid include_items value: %s", includeItemsStr), nil)
            return
        }
        includeItems = parsedBool
    }

    response, err := o.orderService.GetOrderByID(r.Context(), orderID, includeItems)
    if err != nil {
        switch {
        case errors.Is(err, ErrOrderNotFound):
            o.respondWithError(w, http.StatusNotFound, "Order not found", nil)
        case errors.Is(err, ErrForbidden):
            o.respondWithError(w, http.StatusForbidden, "You don't have permission to view this order", nil)
        default:
            o.respondWithError(w, http.StatusInternalServerError, "Failed to get order", nil)
        }
        return
    }

    o.respondWithSuccess(w, http.StatusOK, response)
}

func (o *OrderHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(AuthUser)
    if !ok {
        o.respondWithError(w, http.StatusUnauthorized, "Authentication required", nil)
        return
    }
   
	if user.Role != "admin" {
		o.respondWithError(w, http.StatusForbidden, "Only admins can update the status of orders", nil)
		return
	}

	orderID := chi.URLParam(r, "id")
	if orderID == "" {
		o.respondWithError(w, http.StatusBadRequest, "Order ID is required", nil)
		return
	}

    var req dto.UpdateOrderStatusRequest
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        o.respondWithError(w, http.StatusBadRequest, "Invalid JSON format: " + err.Error(), nil)
        return
    }
    
    if err := o.validator.Struct(req); err != nil {
        validationErrors := dto.FormatValidationErrors(err)
        o.respondWithError(w, http.StatusUnprocessableEntity, "Validation failed", validationErrors)
        return
    }
    
    response, err := o.orderService.UpdateOrderStatus(r.Context(), orderID, &req)
    if err != nil {
	  o.respondWithError(w, http.StatusInternalServerError, "Failed to update the order status!", nil)
	  return
    }

    o.respondWithSuccess(w, http.StatusOK, response)
}

func (o *OrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(AuthUser)
    if !ok {
        o.respondWithError(w, http.StatusUnauthorized, "Authentication required", nil)
        return
    }

	order_id := chi.URLParam(r, "id")
	if order_id == "" {
        o.respondWithError(w, http.StatusBadRequest, "Order ID is required", nil)
        return
    }

	var req dto.CanclerOrderRequest

	if err := json.NewDecoder(r.Body).Decode(&cancelReq); err != nil {
		o.respondWithError(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	if err := o.validator.Struct(cancelReq); err != nil {
		validationErrors := dto.FormatValidationErrors(err)
		o.respondWithError(w, http.StatusUnprocessableEntity, "Validation failed", validationErrors)
		return
	}

	response, err := o.orderService.CancelOrder(r.Context(), orderID, &cancelReq)
	if err != nil {
		switch {
		case errors.Is(err, ErrOrderNotFound):
			o.respondWithError(w, http.StatusNotFound, "Order not found", nil)
		case errors.Is(err, ErrForbidden):
			o.respondWithError(w, http.StatusForbidden, "You don't have permission to cancel this order", nil)
		case errors.Is(err, ErrOrderNotCancellable):
			o.respondWithError(w, http.StatusConflict, err.Error(), nil)
		case errors.Is(err, ErrMissingReason):
			o.respondWithError(w, http.StatusBadRequest, "Cancellation reason is required", nil)
		default:
			o.respondWithError(w, http.StatusInternalServerError, "Failed to cancel order", nil)
		}
		return
	}

	o.respondWithSuccess(w, http.StatusOK, response)
}
