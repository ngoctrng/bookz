package httpserver

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	accounthttp "github.com/ngoctrng/bookz/internal/account/delivery"
	"github.com/ngoctrng/bookz/internal/account/repository"
	"github.com/ngoctrng/bookz/internal/account/usecases"
	"github.com/ngoctrng/bookz/pkg/config"
	"gorm.io/gorm"
)

type Server struct {
	cfg  *config.Config
	echo *echo.Echo
}

func New(cfg *config.Config, db *gorm.DB) *Server {
	e := echo.New()

	accountHandlers := initAccount(db, cfg)

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	api := e.Group("/api")
	api.POST("/account/register", accountHandlers.Register)
	api.POST("/account/login", accountHandlers.Login)

	return &Server{echo: e, cfg: cfg}
}

func (s *Server) Start() error {
	slog.Info("server started!", "port", s.cfg.Port)
	return s.echo.Start(fmt.Sprintf(":%d", s.cfg.Port))
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.echo.ServeHTTP(w, r)
}

func initAccount(db *gorm.DB, cfg *config.Config) *accounthttp.Handler {
	aRepo := repository.New(db)
	uc := usecases.NewService(aRepo)
	aHandler := accounthttp.NewHandler(cfg, uc)
	return aHandler
}
