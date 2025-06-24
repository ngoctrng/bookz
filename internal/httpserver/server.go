package httpserver

import (
	"fmt"
	"log/slog"
	"net/http"
	"slices"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	accounthttp "github.com/ngoctrng/bookz/internal/account/delivery"
	accountrepo "github.com/ngoctrng/bookz/internal/account/repository"
	accountusecases "github.com/ngoctrng/bookz/internal/account/usecases"
	bookhttp "github.com/ngoctrng/bookz/internal/book/delivery"
	bookrepo "github.com/ngoctrng/bookz/internal/book/repository"
	bookusecases "github.com/ngoctrng/bookz/internal/book/usecases"
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
	bookHandlers := initBook(db)

	e.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(cfg.TokenSecret),
		Skipper: func(c echo.Context) bool {
			publicRoutes := []string{
				"/api/account/register",
				"/api/account/login",
				"/api/books/:isbn",
				"/api/books",
				"/health",
			}
			return slices.Contains(publicRoutes, c.Path())
		},
	}))

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	api := e.Group("/api")

	// Public routes
	api.POST("/account/register", accountHandlers.Register)
	api.POST("/account/login", accountHandlers.Login)
	api.GET("/books/:isbn", bookHandlers.Get)
	api.GET("/books", bookHandlers.List)

	// Authenticated routes
	auth := api.Group("")
	auth.Use(authMiddleware(cfg))

	auth.POST("/books", bookHandlers.Create)
	auth.PUT("/books/:isbn", bookHandlers.Update)
	auth.DELETE("/books/:isbn", bookHandlers.Delete)

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

// Example authentication middleware
func authMiddleware(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing or invalid authentication")
			}
			c.Set("user_id", authHeader)
			return next(c)
		}
	}
}
