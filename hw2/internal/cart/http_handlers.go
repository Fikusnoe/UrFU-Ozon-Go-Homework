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

type addRequestBody struct {
	Count uint32 `json:"count"`
}

type cartResponse struct {
	Items      []CartItem `json:"items"`
	TotalPrice uint32     `json:"total_price"`
}

func NewHTTPHandler(cartService *CartService) *HTTPHandler {
	return &HTTPHandler{CartService: cartService}
}

func (h *HTTPHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.ParseInt(r.PathValue("user_id"), 10, 64)

	sku, _ := strconv.ParseInt(r.PathValue("sku"), 10, 64)

	body := addRequestBody{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil{
		http.Error(w, "Неверный формат JSON в теле запроса", http.StatusBadRequest)
		return
	}

	err = h.CartService.AddItem(r.Context(), UserID(userID), sku,  body.Count)
	if err != nil {
		fmt.Printf("Ошибка добавления предмета в корзину: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, "ok")
}

func (h *HTTPHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.ParseInt(r.PathValue("user_id"), 10, 64)

	items, totalPrice, _ := h.CartService.GetCart(UserID(userID))

	response := cartResponse{
		Items:      items,
		TotalPrice: totalPrice,
	}

	writeJSON(w, http.StatusOK, response)
}

func writeJSON(w http.ResponseWriter, code int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(data)
}
