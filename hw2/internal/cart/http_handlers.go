package cart

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type HTTPHandler struct {
	CartService *CartService
}

func NewHTTPHandler(cartService *CartService) *HTTPHandler {
	return &HTTPHandler{ CartService: cartService }
}

func (h *HTTPHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	userID, _ := strconv.ParseInt(query.Get("user_id"), 10, 64)

	sku, _ := strconv.ParseInt(query.Get("sku"), 10, 64)

	count, _ := strconv.ParseUint("count", 10, 32)

	err:= h.CartService.AddItem(r.Context(), UserID(userID), sku, uint32(count))
	if err != nil {
		fmt.Println("Ошибка добавления предмета в корзину")
	}

	writeJSON(w, http.StatusOK, "ok")
}

func (h *HTTPHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.ParseInt(r.URL.Query().Get("user_id"), 10, 64)

	userCart, _ := h.CartService.GetCart(UserID(userID))

	writeJSON(w, http.StatusOK, userCart)
}

func writeJSON(w http.ResponseWriter, code int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(data)
}

