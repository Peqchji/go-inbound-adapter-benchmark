package rest

import (
	"net/http"

	"github.com/Peqchji/go-inbound-adapter-benchmark/internal/domain/wallet"
	"github.com/labstack/echo/v4"
)

type WalletHandler struct {
	service *wallet.WalletService
}

func NewWalletHandler(service *wallet.WalletService) *WalletHandler {
	return &WalletHandler{
		service: service,
	}
}

func (h *WalletHandler) RegisterRoutes(e *echo.Echo) {
	e.GET("/wallets/:id", h.GetWallet)
	e.POST("/wallets", h.CreateWallet)
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

func (h *WalletHandler) GetWallet(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.String(http.StatusBadRequest, "missing id")
	}

	result := h.service.GetWallet(id)
	if result.Err != nil {
		return c.String(http.StatusNotFound, result.Err.Error())
	}

	return c.JSON(http.StatusOK, toWalletResponse(result.Res))
}

type CreateWalletRequest struct {
	ID        string `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

func (h *WalletHandler) CreateWallet(c echo.Context) error {
	var req CreateWalletRequest
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	result := h.service.CreateWallet(req.ID, req.Firstname, req.Lastname)
	if result.Err != nil {
		return c.String(http.StatusInternalServerError, result.Err.Error())
	}

	return c.JSON(http.StatusCreated, toWalletResponse(result.Res))
}
