package httpserver

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/ngoctrng/bookz/api" // Import the generated docs
	accounthttp "github.com/ngoctrng/bookz/internal/account/delivery"
	accountrepo "github.com/ngoctrng/bookz/internal/account/repository"
	accountusecases "github.com/ngoctrng/bookz/internal/account/usecases"
	bookhttp "github.com/ngoctrng/bookz/internal/book/delivery"
	bookrepo "github.com/ngoctrng/bookz/internal/book/repository"
	bookusecases "github.com/ngoctrng/bookz/internal/book/usecases"
	exchangehttp "github.com/ngoctrng/bookz/internal/exchange/delivery"
	exchangebus "github.com/ngoctrng/bookz/internal/exchange/messagebus"
	exchangerepo "github.com/ngoctrng/bookz/internal/exchange/repository"
	exchangeusecases "github.com/ngoctrng/bookz/internal/exchange/usecases"
	"github.com/ngoctrng/bookz/pkg/config"
	"github.com/ngoctrng/bookz/pkg/token"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gorm.io/gorm"
)

type Server struct {
	cfg  *config.Config
	echo *echo.Echo
}

// @title           Bookz API
// @version         1.0
// @description     The FCC Book Trading Club is a backend system designed to manage a community-driven book trading platform.
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8088
// @BasePath  /api

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func New(cfg *config.Config, db *gorm.DB, client *asynq.Client) *Server {
	e := echo.New()

	publicRoutes := [][]string{
		{http.MethodPost, "/api/account/register"},
		{http.MethodPost, "/api/account/login"},
		{http.MethodGet, "/api/books/:id"},
		{http.MethodGet, "/api/books"},
		{http.MethodGet, "/health"},
		{http.MethodGet, "/swagger/*"},
		{http.MethodGet, "/swagger/index.html"},
		{http.MethodGet, "/swagger/doc.json"},
	}

	accountHandlers := initAccount(db, cfg)
	bookHandlers := initBook(db)
	exchangeHandlers := initExchange(db, client)

	configCORS(e, cfg)
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Use(middleware.RequestID())
	e.Use(middleware.Gzip())
	e.Use(sentryecho.New(sentryecho.Options{Repanic: true}))
	e.Use(authMiddleware(cfg, publicRoutes))

	e.GET("/swagger/*", echoSwagger.WrapHandler)
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
	api.POST("/exchange/proposals", exchangeHandlers.CreateProposal)
	api.GET("/exchange/proposals/:id", exchangeHandlers.GetProposalByID)
	api.GET("/exchange/proposals", exchangeHandlers.GetAllProposals)
	api.PUT("/exchange/proposals/:id/accept", exchangeHandlers.AcceptProposal)

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
	r := accountrepo.New(db)
	uc := accountusecases.NewService(r)
	aHandler := accounthttp.NewHandler(cfg, uc)
	return aHandler
}

func initBook(db *gorm.DB) *bookhttp.Handler {
	r := bookrepo.New(db)
	uc := bookusecases.NewService(r)
	bHandler := bookhttp.NewHandler(uc)
	return bHandler
}

func initExchange(db *gorm.DB, client *asynq.Client) *exchangehttp.Handler {
	r := exchangerepo.New(db)
	bus := exchangebus.New(client)
	uc := exchangeusecases.NewProposalService(r, bus)
	return exchangehttp.NewHandler(uc)
}

func authMiddleware(cfg *config.Config, publicRoutes [][]string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			method := c.Request().Method
			path := c.Path()

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
