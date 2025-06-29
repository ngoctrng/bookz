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

// Register godoc
// @Summary      Register a new user
// @Description  Create a new user account
// @Tags         account
// @Accept       json
// @Produce      json
// @Param        body  body      RegisterRequest  true  "User registration info"
// @Success      201   {object}  TokenResponse
// @Failure      400   {object}  echo.HTTPError
// @Failure      500   {object}  echo.HTTPError
// @Router       /account/register [post]
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

// Login godoc
// @Summary      Login
// @Description  Authenticate user and return JWT token
// @Tags         account
// @Accept       json
// @Produce      json
// @Param        body  body      LoginRequest  true  "User login info"
// @Success      200   {object}  TokenResponse
// @Failure      400   {object}  echo.HTTPError
// @Failure      500   {object}  echo.HTTPError
// @Router       /account/login [post]
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
