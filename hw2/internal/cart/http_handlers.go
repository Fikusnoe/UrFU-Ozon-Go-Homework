package cart

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
)

type HTTPHandler struct {
	CartService *CartService
}

type addRequestBody struct {
	Count int32 `json:"count"`
}

type cartResponse struct {
	Items      []CartItem `json:"items"`
	TotalPrice uint32     `json:"total_price"`
}

func NewHTTPHandler(cartService *CartService) *HTTPHandler {
	return &HTTPHandler{CartService: cartService}
}

func (h *HTTPHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(r.PathValue("user_id"), 10, 64)
	if err != nil {
		http.Error(w,
			"Идентификатор пользователя должен быть натуральным числом (больше нуля)",
			http.StatusBadRequest,
		)
		return
	}

	sku, err := strconv.ParseInt(r.PathValue("sku"), 10, 64)
	if err != nil || sku <= 0 {
		http.Error(
			w,
			"SKU должен быть натуральным числом (больше нуля)",
			http.StatusBadRequest,
		)
		return
	}

	body := addRequestBody{}
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Неверный формат JSON в теле запроса", http.StatusBadRequest)
		return
	}

	if body.Count <= 0 {
		http.Error(
			w,
			"Количество должно быть натуральным числом (больше нуля)",
			http.StatusBadRequest,
		)
		return
	}

	err = h.CartService.AddItem(r.Context(), UserID(userID), sku, uint32(body.Count))
	if err != nil {
		log.Printf("Ошибка добавления предмета в корзину: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, "ok")
}

func (h *HTTPHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(r.PathValue("user_id"), 10, 64)
	if err != nil {
		http.Error(w,
			"Идентификатор пользователя должен быть натуральным числом (больше нуля)",
			http.StatusBadRequest,
		)
		return
	}

	items, totalPrice, err := h.CartService.GetCart(r.Context(), UserID(userID))
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	response := cartResponse{
		Items:      items,
		TotalPrice: totalPrice,
	}

	writeJSON(w, http.StatusOK, response)
}

func writeJSON(w http.ResponseWriter, code int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, "Неверный формат JSON в теле запроса", http.StatusBadRequest)
		return
	}
}

func (h *HTTPHandler) handleServiceError(w http.ResponseWriter, err error) {
	var validationErr *ValidationError

	switch {
	case errors.As(err, &validationErr):
		http.Error(w, validationErr.Message, http.StatusBadRequest)

	case errors.Is(err, ErrProductNotFound):
		http.Error(
			w,
			"SKU должен существовать в сервисе product-service",
			http.StatusPreconditionFailed,
		)

	case errors.Is(err, ErrCartNotFound):
		http.Error(w, "Корзина не найдена", http.StatusNotFound)

	default:
		log.Printf("internal error: %v", err)
		http.Error(
			w,
			"Проблемы из-за неисправностей в системе",
			http.StatusInternalServerError,
		)
	}
}
