package delivery_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/ngoctrng/bookz/internal/book"
	"github.com/ngoctrng/bookz/internal/book/delivery"
	"github.com/ngoctrng/bookz/internal/book/mocks"
	"github.com/stretchr/testify/assert"
)

func TestCreateBookHandler(t *testing.T) {
	reqBody := delivery.BookRequest{ISBN: "123", Title: "Go", Author: "Alan", Year: 2024}
	body, _ := json.Marshal(reqBody)

	uc := new(mocks.MockUsecase)
	h := delivery.NewHandler(uc)

	uc.EXPECT().
		Create(book.New("owner1", "123", "Go", "Alan", 2024)).
		Return(nil).Once()

	c, rec := setupEchoConext("/books", bytes.NewReader(body))
	c.Set("user_id", "owner1")

	err := h.Create(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	uc.AssertExpectations(t)
}

func TestGetBookHandler(t *testing.T) {
	b := &book.BookInfo{ID: 1, ISBN: "123", Title: "Go"}
	uc := new(mocks.MockUsecase)
	uc.EXPECT().Get(1).Return(b, nil).Once()

	h := delivery.NewHandler(uc)
	c, rec := setupEchoConext("/books/1", nil)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := h.Get(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	uc.AssertExpectations(t)
}

func TestUpdateBookHandler(t *testing.T) {
	reqBody := delivery.BookRequest{ISBN: "123", Title: "Go Updated", Author: "Alan", Year: 2025}
	body, _ := json.Marshal(reqBody)

	uc := new(mocks.MockUsecase)
	uc.EXPECT().
		Update(book.New("owner1", "123", "Go Updated", "Alan", 2025), "owner1").
		Return(nil).Once()

	h := delivery.NewHandler(uc)
	c, rec := setupEchoConext("/books/123", bytes.NewReader(body))
	c.SetParamNames("isbn")
	c.SetParamValues("123")
	c.Set("user_id", "owner1")

	err := h.Update(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	uc.AssertExpectations(t)
}

func TestDeleteBookHandler(t *testing.T) {
	uc := new(mocks.MockUsecase)
	uc.EXPECT().Delete(1, "owner1").Return(nil).Once()

	h := delivery.NewHandler(uc)
	c, rec := setupEchoConext("/books/1", nil)
	c.SetParamNames("id")
	c.SetParamValues("1")
	c.Set("user_id", "owner1")

	err := h.Delete(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	uc.AssertExpectations(t)
}

func TestListBooksHandler(t *testing.T) {
	books := []*book.BookInfo{
		{ISBN: "1", Title: "Book 1", Author: "Author 1", Year: 2021, Owner: book.BookOwner{ID: "owner1"}},
		{ISBN: "2", Title: "Book 2", Author: "Author 2", Year: 2022, Owner: book.BookOwner{ID: "owner2"}},
	}
	uc := new(mocks.MockUsecase)
	uc.EXPECT().List().Return(books, nil).Once()

	h := delivery.NewHandler(uc)
	c, rec := setupEchoConext("/books", nil)

	err := h.List(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	uc.AssertExpectations(t)
}

func TestBookHandlerErrors(t *testing.T) {
	uc := new(mocks.MockUsecase)
	uc.EXPECT().Get(222).Return(nil, errors.New("not found")).Once()

	h := delivery.NewHandler(uc)

	c, _ := setupEchoConext("/books/222", nil)
	c.SetParamNames("id")
	c.SetParamValues("222")

	err := h.Get(c)
	assert.Error(t, err)
	uc.AssertExpectations(t)
}

func setupEchoConext(path string, body io.Reader) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, path, body)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := e.NewContext(req, rec)
	return c, rec
}
