package delivery_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/ngoctrng/bookz/internal/exchange"
	"github.com/ngoctrng/bookz/internal/exchange/delivery"
	"github.com/ngoctrng/bookz/internal/exchange/mocks"
	"github.com/ngoctrng/bookz/internal/exchange/usecases"
	"github.com/stretchr/testify/assert"
)

func TestCreateProposalHandler(t *testing.T) {
	e := echo.New()
	mockUC := new(mocks.MockUsecase)
	h := delivery.NewHandler(mockUC)

	uid := uuid.New()
	reqBody := delivery.CreateProposalRequest{
		Requested:     2,
		ForExchangeID: 3,
		Message:       "Let's trade!",
	}
	body, _ := json.Marshal(reqBody)

	input := usecases.CreateProposalInput{
		RequestBy:     uid,
		Requested:     2,
		ForExchangeID: 3,
		Message:       "Let's trade!",
	}

	mockUC.EXPECT().CreateProposal(input).Return(nil).Once()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/exchange/proposals", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := e.NewContext(req, rec)
	c.Set("user_id", uid.String())

	err := h.CreateProposal(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	mockUC.AssertExpectations(t)
}

func TestCreateProposalHandler_InvalidRequest(t *testing.T) {
	e := echo.New()
	mockUC := new(mocks.MockUsecase)
	h := delivery.NewHandler(mockUC)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/exchange/proposals", bytes.NewReader([]byte("not-json")))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := e.NewContext(req, rec)
	c.Set("user_id", uuid.New().String())

	err := h.CreateProposal(c)
	httpErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, httpErr.Code)
}

func TestCreateProposalHandler_UsecaseError(t *testing.T) {
	e := echo.New()
	mockUC := new(mocks.MockUsecase)
	h := delivery.NewHandler(mockUC)

	uid := uuid.New()
	reqBody := delivery.CreateProposalRequest{
		Requested:     2,
		ForExchangeID: 3,
		Message:       "Let's trade!",
	}
	body, _ := json.Marshal(reqBody)

	input := usecases.CreateProposalInput{
		RequestBy:     uid,
		Requested:     2,
		ForExchangeID: 3,
		Message:       "Let's trade!",
	}

	mockUC.EXPECT().CreateProposal(input).Return(errors.New("fail")).Once()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/exchange/proposals", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c := e.NewContext(req, rec)
	c.Set("user_id", uid.String())

	err := h.CreateProposal(c)
	httpErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusInternalServerError, httpErr.Code)
	mockUC.AssertExpectations(t)
}

func TestGetProposalByIDHandler(t *testing.T) {
	e := echo.New()
	mockUC := new(mocks.MockUsecase)
	h := delivery.NewHandler(mockUC)

	proposal := &exchange.Proposal{ID: 1, RequestBy: uuid.New()}
	mockUC.EXPECT().GetProposalByID(1).Return(proposal, nil).Once()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/exchange/proposals/1", nil)
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")
	c.Set("user_id", proposal.RequestBy.String())

	err := h.GetProposalByID(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUC.AssertExpectations(t)
}

func TestGetProposalByIDHandler_NotFound(t *testing.T) {
	e := echo.New()
	mockUC := new(mocks.MockUsecase)
	h := delivery.NewHandler(mockUC)

	mockUC.EXPECT().GetProposalByID(999).Return(nil, errors.New("not found")).Once()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/exchange/proposals/999", nil)
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("999")
	c.Set("user_id", uuid.New().String())

	err := h.GetProposalByID(c)
	httpErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusNotFound, httpErr.Code)
	mockUC.AssertExpectations(t)
}

func TestGetAllProposalsHandler(t *testing.T) {
	e := echo.New()
	mockUC := new(mocks.MockUsecase)
	h := delivery.NewHandler(mockUC)

	uid := uuid.New()
	proposals := []*exchange.Proposal{
		{ID: 1}, {ID: 2},
	}
	mockUC.EXPECT().GetAllProposals(uid).Return(proposals, nil).Once()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/exchange/proposals", nil)
	c := e.NewContext(req, rec)
	c.Set("user_id", uid.String())

	err := h.GetAllProposals(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUC.AssertExpectations(t)
}

func TestGetAllProposalsHandler_UsecaseError(t *testing.T) {
	e := echo.New()
	mockUC := new(mocks.MockUsecase)
	h := delivery.NewHandler(mockUC)

	uid := uuid.New()
	mockUC.EXPECT().GetAllProposals(uid).Return(nil, errors.New("fail")).Once()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/exchange/proposals", nil)
	c := e.NewContext(req, rec)
	c.Set("user_id", uid.String())

	err := h.GetAllProposals(c)
	httpErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusInternalServerError, httpErr.Code)
	mockUC.AssertExpectations(t)
}
