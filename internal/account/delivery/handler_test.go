package delivery_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/ngoctrng/bookz/internal/account"
	"github.com/ngoctrng/bookz/internal/account/delivery"
	"github.com/ngoctrng/bookz/internal/account/mocks"
	"github.com/ngoctrng/bookz/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterHandler(t *testing.T) {
	e := echo.New()

	cfg := &config.Config{
		TokenSecret:     "testsecret",
		TokenExpiration: int((time.Minute * 15).Seconds()),
	}

	validReq := delivery.RegisterRequest{
		Username:        "testuser",
		Email:           "test@example.com",
		Password:        "password123",
		PasswordConfirm: "password123",
	}
	body, _ := json.Marshal(validReq)

	t.Run("should register user and return token", func(t *testing.T) {
		mockUC := new(mocks.MockUsecase)
		mockUC.EXPECT().
			Register(
				mock.AnythingOfType("uuid.UUID"),
				validReq.Username,
				validReq.Email,
				mock.AnythingOfType("string"),
			).
			Return(nil).Once()
		h := delivery.NewHandler(cfg, mockUC)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c := e.NewContext(req, rec)

		err := h.Register(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)

		var resp delivery.TokenResponse
		_ = json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NotEmpty(t, resp.Token)
		mockUC.AssertExpectations(t)
	})

	t.Run("should return 400 for invalid request", func(t *testing.T) {
		mockUC := new(mocks.MockUsecase)
		h := delivery.NewHandler(cfg, mockUC)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader([]byte("somedata")))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c := e.NewContext(req, rec)

		err := h.Register(c)

		httpErr, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, httpErr.Code)
	})

	t.Run("should return 500 if usecase returns error", func(t *testing.T) {
		mockUC := new(mocks.MockUsecase)
		mockUC.EXPECT().
			Register(
				mock.AnythingOfType("uuid.UUID"),
				validReq.Username,
				validReq.Email,
				mock.AnythingOfType("string"),
			).
			Return(errors.New("failure")).Once()
		h := delivery.NewHandler(cfg, mockUC)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c := e.NewContext(req, rec)

		err := h.Register(c)
		httpErr, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, httpErr.Code)
		mockUC.AssertExpectations(t)
	})
}

func TestLoginSuccess(t *testing.T) {
	e := echo.New()
	cfg := &config.Config{
		TokenSecret:     "testsecret",
		TokenExpiration: int((time.Minute * 15).Seconds()),
	}
	mockUC := new(mocks.MockUsecase)
	h := delivery.NewHandler(cfg, mockUC)

	user := &account.User{
		ID:       uuid.New(),
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	reqBody := delivery.LoginRequest{
		Email:    user.Email,
		Password: user.Password,
	}
	body, _ := json.Marshal(reqBody)

	t.Run("should login successfully and return token", func(t *testing.T) {
		mockUC.EXPECT().Login(user.Email, user.Password).Return(user, nil).Once()

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c := e.NewContext(req, rec)

		err := h.Login(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp delivery.TokenResponse
		_ = json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NotEmpty(t, resp.Token)
		mockUC.AssertExpectations(t)
	})
}

func TestLoginInvalidRequest(t *testing.T) {
	e := echo.New()
	cfg := &config.Config{
		TokenSecret:     "testsecret",
		TokenExpiration: int((time.Minute * 15).Seconds()),
	}
	mockUC := new(mocks.MockUsecase)
	h := delivery.NewHandler(cfg, mockUC)

	t.Run("should return 400 for invalid request", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader([]byte("not-json")))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c := e.NewContext(req, rec)

		err := h.Login(c)
		httpErr, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, httpErr.Code)
	})
}

func TestLoginWithInvalidCredentials(t *testing.T) {
	e := echo.New()
	cfg := &config.Config{
		TokenSecret:     "testsecret",
		TokenExpiration: int((time.Minute * 15).Seconds()),
	}
	mockUC := new(mocks.MockUsecase)
	h := delivery.NewHandler(cfg, mockUC)

	reqBody := delivery.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}
	body, _ := json.Marshal(reqBody)

	t.Run("should return 400 if invalid credentials", func(t *testing.T) {
		mockUC.EXPECT().Login(reqBody.Email, reqBody.Password).Return(nil, errors.New("invalid credentials")).Once()

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c := e.NewContext(req, rec)

		err := h.Login(c)
		httpErr, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, httpErr.Code)
		mockUC.AssertExpectations(t)
	})
}
