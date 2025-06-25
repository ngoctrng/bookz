package delivery

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/ngoctrng/bookz/internal/book"
	"github.com/ngoctrng/bookz/internal/book/usecases"
)

type Handler struct {
	uc usecases.Usecase
}

func NewHandler(uc usecases.Usecase) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) Create(c echo.Context) error {
	ownerID := c.Get("user_id").(string)

	var req BookRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request format")
	}

	if err := req.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	b := &book.Book{
		OwnerID:     ownerID,
		ISBN:        req.ISBN,
		Title:       req.Title,
		Description: req.Description,
		BriefReview: req.BriefReview,
		Author:      req.Author,
		Year:        req.Year,
	}
	if err := h.uc.Create(b); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())

	}
	return c.JSON(http.StatusCreated, req)
}

func (h *Handler) Get(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	b, err := h.uc.Get(id)
	if err != nil || b == nil {
		return echo.NewHTTPError(http.StatusNotFound, "book not found")
	}

	resp := BookResponse{
		ISBN:        b.ISBN,
		Title:       b.Title,
		Description: b.Description,
		BriefReview: b.BriefReview,
		Author:      b.Author,
		Year:        b.Year,
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) Update(c echo.Context) error {
	ownerID := c.Get("user_id").(string)
	isbn := c.Param("isbn")

	var req BookRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request format")
	}

	if err := req.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	b := &book.Book{
		OwnerID:     ownerID,
		ISBN:        isbn,
		Title:       req.Title,
		Description: req.Description,
		BriefReview: req.BriefReview,
		Author:      req.Author,
		Year:        req.Year,
	}

	if err := h.uc.Update(b, ownerID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, req)
}

func (h *Handler) Delete(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	ownerID := c.Get("user_id").(string)

	if err := h.uc.Delete(id, ownerID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) List(c echo.Context) error {
	books, err := h.uc.List()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, books)
}
