package handlers_test

// import (
// 	"context"
// 	"errors"
// 	"loopit/internal/api/handlers"
// 	"loopit/internal/constants"
// 	"loopit/internal/enums"
// 	"loopit/internal/enums/order_status"
// 	"loopit/internal/initializer"
// 	"loopit/internal/models"
// 	"loopit/pkg/logger"
// 	"net/http"
// 	"net/http/httptest"
// 	"strings"
// 	"testing"
// )

// // FakeOrderService is a test double for OrderServiceInterface
// type FakeOrderService struct {
// 	UpdateOrderStatusFn   func(orderID int, newStatus order_status.Status) error
// 	GetOrderHistoryFn     func(userCtx *models.UserContext, filterStatus []order_status.Status) ([]*models.Order, error)
// 	GetApprovedAwaitingFn func(userCtx *models.UserContext) ([]*models.Order, error)
// 	MarkOrderAsReturnedFn func(orderID int, userCtx *models.UserContext) error
// 	GetLenderOrdersFn     func(userCtx *models.UserContext) ([]*models.Order, error)
// }

// // UpdateOrderStatus mock
// func (f *FakeOrderService) UpdateOrderStatus(orderID int, newStatus order_status.Status) error {
// 	if f.UpdateOrderStatusFn != nil {
// 		return f.UpdateOrderStatusFn(orderID, newStatus)
// 	}
// 	return nil
// }

// // GetOrderHistory mock
// func (f *FakeOrderService) GetOrderHistory(userCtx *models.UserContext, filterStatus []order_status.Status) ([]*models.Order, error) {
// 	if f.GetOrderHistoryFn != nil {
// 		return f.GetOrderHistoryFn(userCtx, filterStatus)
// 	}
// 	return nil, nil
// }

// // GetAllApprovedAwaitingOrders mock
// func (f *FakeOrderService) GetAllApprovedAwaitingOrders(userCtx *models.UserContext) ([]*models.Order, error) {
// 	if f.GetApprovedAwaitingFn != nil {
// 		return f.GetApprovedAwaitingFn(userCtx)
// 	}
// 	return nil, nil
// }

// // MarkOrderAsReturned mock
// func (f *FakeOrderService) MarkOrderAsReturned(orderID int, userCtx *models.UserContext) error {
// 	if f.MarkOrderAsReturnedFn != nil {
// 		return f.MarkOrderAsReturnedFn(orderID, userCtx)
// 	}
// 	return nil
// }

// // GetLenderOrders mock
// func (f *FakeOrderService) GetLenderOrders(userCtx *models.UserContext) ([]*models.Order, error) {
// 	if f.GetLenderOrdersFn != nil {
// 		return f.GetLenderOrdersFn(userCtx)
// 	}
// 	return nil, nil
// }

// // ----------------- Test GetOrderHistory -----------------
// func TestGetOrderHistory(t *testing.T) {
// 	log := logger.NewFakeLogger()

// 	tests := []struct {
// 		name       string
// 		userCtx    any
// 		query      string
// 		serviceFn  func(user *models.UserContext, status []order_status.Status) ([]*models.Order, error)
// 		wantStatus int
// 		wantBody   string
// 	}{
// 		{"no user context", nil, "", nil, http.StatusUnauthorized, "unauthorized"},
// 		{"invalid status", &models.UserContext{ID: 1}, "?status=WRONG", nil, http.StatusBadRequest, "invalid status"},
// 		{"service error", &models.UserContext{ID: 1}, "?status=Returned",
// 			func(u *models.UserContext, s []order_status.Status) ([]*models.Order, error) {
// 				return nil, errors.New("db error")
// 			}, http.StatusInternalServerError, "failed to fetch order history"},
// 		{"success", &models.UserContext{ID: 2}, "?status=Returned",
// 			func(u *models.UserContext, s []order_status.Status) ([]*models.Order, error) {
// 				return []*models.Order{{ID: 1}}, nil
// 			}, http.StatusOK, `"status":true`},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			initializer.OrderService = &FakeOrderService{GetOrderHistoryFn: tt.serviceFn}
// 			h := handlers.NewOrderHandler(initializer.OrderService, log)

// 			req := httptest.NewRequest(http.MethodGet, "/orders/history"+tt.query, nil)
// 			if tt.userCtx != nil {
// 				req = req.WithContext(context.WithValue(req.Context(), constants.UserCtxKey, tt.userCtx))
// 			}
// 			w := httptest.NewRecorder()
// 			h.GetOrderHistory(w, req)

// 			if w.Result().StatusCode != tt.wantStatus {
// 				t.Errorf("expected %d got %d", tt.wantStatus, w.Result().StatusCode)
// 			}
// 			if !strings.Contains(w.Body.String(), tt.wantBody) {
// 				t.Errorf("expected body to contain %q got %s", tt.wantBody, w.Body.String())
// 			}
// 		})
// 	}
// }

// // ----------------- Test MarkOrderAsReturned -----------------
// func TestMarkOrderAsReturned(t *testing.T) {
// 	log := logger.NewFakeLogger()

// 	tests := []struct {
// 		name       string
// 		userCtx    any
// 		orderID    string
// 		serviceFn  func(orderID int, u *models.UserContext) error
// 		wantStatus int
// 		wantBody   string
// 	}{
// 		{"no user", nil, "1", nil, http.StatusUnauthorized, "unauthorized"},
// 		{"invalid id", &models.UserContext{ID: 1}, "bad", nil, http.StatusBadRequest, "invalid order id"},
// 		{"service error", &models.UserContext{ID: 1}, "1",
// 			func(id int, u *models.UserContext) error { return errors.New("fail") },
// 			http.StatusBadRequest, "failed to update"},
// 		{"success", &models.UserContext{ID: 2}, "2",
// 			func(id int, u *models.UserContext) error { return nil },
// 			http.StatusOK, `"status":true`},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			initializer.OrderService = &FakeOrderService{MarkOrderAsReturnedFn: tt.serviceFn}
// 			h := handlers.NewOrderHandler(initializer.OrderService, log)

// 			req := httptest.NewRequest(http.MethodPatch, "/orders/"+tt.orderID+"/return", nil)
// 			if tt.userCtx != nil {
// 				req = req.WithContext(context.WithValue(req.Context(), constants.UserCtxKey, tt.userCtx))
// 			}
// 			req.SetPathValue("orderId", tt.orderID)
// 			w := httptest.NewRecorder()

// 			h.MarkOrderAsReturned(w, req)

// 			if w.Result().StatusCode != tt.wantStatus {
// 				t.Errorf("expected %d got %d", tt.wantStatus, w.Result().StatusCode)
// 			}
// 			if !strings.Contains(w.Body.String(), tt.wantBody) {
// 				t.Errorf("expected body contain %q got %s", tt.wantBody, w.Body.String())
// 			}
// 		})
// 	}
// }

// // ----------------- Test GetAllApprovedAwaitingOrders -----------------
// func TestGetAllApprovedAwaitingOrders(t *testing.T) {
// 	log := logger.NewFakeLogger()

// 	tests := []struct {
// 		name       string
// 		userCtx    *models.UserContext
// 		serviceFn  func(userCtx *models.UserContext) ([]*models.Order, error)
// 		wantStatus int
// 		wantBody   string
// 	}{
// 		{"unauthorized", nil, nil, http.StatusUnauthorized, "unauthorized"},
// 		{"forbidden - not lender", &models.UserContext{ID: 1, Role: enums.RoleUser}, nil, http.StatusForbidden, "forbidden"},
// 		{"service error", &models.UserContext{ID: 2, Role: enums.RoleLender},
// 			func(u *models.UserContext) ([]*models.Order, error) { return nil, errors.New("fail") },
// 			http.StatusInternalServerError, "failed to fetch approved awaiting"},
// 		{"success", &models.UserContext{ID: 3, Role: enums.RoleLender},
// 			func(u *models.UserContext) ([]*models.Order, error) { return []*models.Order{{ID: 5}}, nil },
// 			http.StatusOK, `"status":true`},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			initializer.OrderService = &FakeOrderService{GetApprovedAwaitingFn: tt.serviceFn}
// 			h := handlers.NewOrderHandler(initializer.OrderService, log)

// 			req := httptest.NewRequest(http.MethodGet, "/orders/approved-awaiting", nil)
// 			if tt.userCtx != nil {
// 				req = req.WithContext(context.WithValue(req.Context(), constants.UserCtxKey, tt.userCtx))
// 			}
// 			w := httptest.NewRecorder()
// 			h.GetAllApprovedAwaitingOrders(w, req)

// 			if w.Result().StatusCode != tt.wantStatus {
// 				t.Errorf("expected %d got %d", tt.wantStatus, w.Result().StatusCode)
// 			}
// 			if !strings.Contains(w.Body.String(), tt.wantBody) {
// 				t.Errorf("expected body contain %q got %s", tt.wantBody, w.Body.String())
// 			}
// 		})
// 	}
// }

// // ----------------- Test GetLenderOrders -----------------
// func TestGetLenderOrders(t *testing.T) {
// 	log := logger.NewFakeLogger()

// 	tests := []struct {
// 		name       string
// 		userCtx    *models.UserContext
// 		serviceFn  func(u *models.UserContext) ([]*models.Order, error)
// 		wantStatus int
// 		wantBody   string
// 	}{
// 		{"unauthorized", nil, nil, http.StatusUnauthorized, "unauthorized"},
// 		{"forbidden - not lender", &models.UserContext{ID: 1, Role: enums.RoleUser}, nil, http.StatusForbidden, "forbidden"},
// 		{"service error", &models.UserContext{ID: 2, Role: enums.RoleLender},
// 			func(u *models.UserContext) ([]*models.Order, error) { return nil, errors.New("fail") },
// 			http.StatusInternalServerError, "failed to fetch lender orders"},
// 		{"success", &models.UserContext{ID: 3, Role: enums.RoleLender},
// 			func(u *models.UserContext) ([]*models.Order, error) { return []*models.Order{{ID: 7}}, nil },
// 			http.StatusOK, `"status":true`},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			initializer.OrderService = &FakeOrderService{GetLenderOrdersFn: tt.serviceFn}
// 			h := handlers.NewOrderHandler(initializer.OrderService, log)

// 			req := httptest.NewRequest(http.MethodGet, "/orders/lender", nil)
// 			if tt.userCtx != nil {
// 				req = req.WithContext(context.WithValue(req.Context(), constants.UserCtxKey, tt.userCtx))
// 			}
// 			w := httptest.NewRecorder()
// 			h.GetLenderOrders(w, req)

// 			if w.Result().StatusCode != tt.wantStatus {
// 				t.Errorf("expected %d got %d", tt.wantStatus, w.Result().StatusCode)
// 			}
// 			if !strings.Contains(w.Body.String(), tt.wantBody) {
// 				t.Errorf("expected body contain %q got %s", tt.wantBody, w.Body.String())
// 			}
// 		})
// 	}
// }
