package rest

import (
	"encoding/json"
	"net/http"

	"github.com/Peqchji/go-inbound-adapter-benchmark/internal/domain/wallet"
)

type WalletHandler struct {
	service *wallet.WalletService
}

func NewWalletHandler(service *wallet.WalletService) *WalletHandler {
	return &WalletHandler{
		service: service,
	}
}

func (h *WalletHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /wallets/{id}", h.GetWallet)
	mux.HandleFunc("POST /wallets", h.CreateWallet)
}

type WalletResponse struct {
	ID      string        `json:"id"`
	Balance uint64        `json:"balance"`
	Owner   OwnerResponse `json:"owner"`
}

type OwnerResponse struct {
	ID        string `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

func toWalletResponse(w wallet.Wallet) WalletResponse {
	o := w.Owner()
	return WalletResponse{
		ID:      w.ID(),
		Balance: w.Balance(),
		Owner: OwnerResponse{
			ID:        o.ID(),
			Firstname: o.Firstname(),
			Lastname:  o.Lastname(),
		},
	}
}

func (h *WalletHandler) GetWallet(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	result := h.service.GetWallet(id)
	if result.Err != nil {
		http.Error(w, result.Err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(toWalletResponse(result.Res))
}

type CreateWalletRequest struct {
	ID        string `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

func (h *WalletHandler) CreateWallet(w http.ResponseWriter, r *http.Request) {
	var req CreateWalletRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := h.service.CreateWallet(req.ID, req.Firstname, req.Lastname)
	if result.Err != nil {
		http.Error(w, result.Err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(toWalletResponse(result.Res))
}
