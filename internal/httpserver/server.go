package httpserver

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	accounthttp "github.com/ngoctrng/bookz/internal/account/delivery"
	accountrepo "github.com/ngoctrng/bookz/internal/account/repository"
	accountusecases "github.com/ngoctrng/bookz/internal/account/usecases"
	bookhttp "github.com/ngoctrng/bookz/internal/book/delivery"
	bookrepo "github.com/ngoctrng/bookz/internal/book/repository"
	bookusecases "github.com/ngoctrng/bookz/internal/book/usecases"
	"github.com/ngoctrng/bookz/pkg/config"
	"github.com/ngoctrng/bookz/pkg/token"
	"gorm.io/gorm"
)

type Server struct {
	cfg  *config.Config
	echo *echo.Echo
}

func New(cfg *config.Config, db *gorm.DB) *Server {
	e := echo.New()

	publicRoutes := [][]string{
		{"POST", "/api/account/register"},
		{"POST", "/api/account/login"},
		{"GET", "/api/books/:id"},
		{"GET", "/api/books"},
		{"GET", "/health"},
	}

	accountHandlers := initAccount(db, cfg)
	bookHandlers := initBook(db)

	configCORS(e, cfg)
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Use(middleware.RequestID())
	e.Use(middleware.Gzip())
	e.Use(sentryecho.New(sentryecho.Options{Repanic: true}))
	e.Use(authMiddleware(cfg, publicRoutes))

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	api := e.Group("/api")

	// public routes
	public := api.Group("")
	public.POST("/account/register", accountHandlers.Register)
	public.POST("/account/login", accountHandlers.Login)
	public.GET("/books/:id", bookHandlers.Get)
	public.GET("/books", bookHandlers.List)

	// private routes
	api.POST("/books", bookHandlers.Create)
	api.PUT("/books/:id", bookHandlers.Update)
	api.DELETE("/books/:id", bookHandlers.Delete)

	return &Server{echo: e, cfg: cfg}
}

func configCORS(e *echo.Echo, cfg *config.Config) {
	if cfg.AllowOrigins != "" {
		aos := strings.Split(cfg.AllowOrigins, ",")
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: aos,
		}))
	}
}

func (s *Server) Start() error {
	slog.Info("server started!", "port", s.cfg.Port)
	return s.echo.Start(fmt.Sprintf(":%d", s.cfg.Port))
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.echo.ServeHTTP(w, r)
}

func initAccount(db *gorm.DB, cfg *config.Config) *accounthttp.Handler {
	aRepo := accountrepo.New(db)
	uc := accountusecases.NewService(aRepo)
	aHandler := accounthttp.NewHandler(cfg, uc)
	return aHandler
}

func initBook(db *gorm.DB) *bookhttp.Handler {
	bRepo := bookrepo.New(db)
	uc := bookusecases.NewService(bRepo)
	bHandler := bookhttp.NewHandler(uc)
	return bHandler
}

func authMiddleware(cfg *config.Config, publicRoutes [][]string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			method := c.Request().Method
			path := c.Request().URL.Path

			// check if the route is public
			for _, route := range publicRoutes {
				if method == route[0] && path == route[1] {
					return next(c)
				}
			}

			// check jwt for private routes
			tk := c.Request().Header.Get("Authorization")
			if tk == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing or invalid tk")
			}

			tk = strings.TrimPrefix(tk, "Bearer ")
			claims, err := token.Verify(tk, cfg.TokenSecret)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid jwt")
			}

			c.Set("user_id", claims.UserID)

			return next(c)
		}
	}
}
