package delivery

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/ngoctrng/bookz/internal/account/usecases"
	"github.com/ngoctrng/bookz/pkg/config"
	"github.com/ngoctrng/bookz/pkg/hasher"
	"github.com/ngoctrng/bookz/pkg/token"
)

type Handler struct {
	cfg *config.Config
	uc  usecases.Usecase
}

func NewHandler(cfg *config.Config, uc usecases.Usecase) *Handler {
	return &Handler{cfg: cfg, uc: uc}
}

func (h *Handler) Register(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request format")
	}

	if err := req.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	hashed, err := hasher.HashPassword(req.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to hash password")
	}

	uid := uuid.New()
	err = h.uc.Register(uid, req.Username, req.Email, hashed)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	signed, err := token.Sign(uid.String(), h.cfg.TokenSecret, 24*time.Hour)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, TokenResponse{
		Token: signed,
	})
}

func (h *Handler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request format")
	}

	if err := req.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := h.uc.Login(req.Email, req.Password)
	if err != nil {
		slog.Info("login failed", "email", req.Email, "error", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid email or password")
	}

	signed, err := token.Sign(user.ID.String(), h.cfg.TokenSecret, 24*time.Hour)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to generate token")
	}

	return c.JSON(http.StatusOK, TokenResponse{
		Token: signed,
	})
}
