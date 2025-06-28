package delivery

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/ngoctrng/bookz/internal/exchange/usecases"
)

type Handler struct {
	uc usecases.Usecase
}

func NewHandler(uc usecases.Usecase) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) CreateProposal(c echo.Context) error {
	var req CreateProposalRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request format")
	}
	if err := req.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	userID := c.Get("user_id").(string)
	uid, err := uuid.Parse(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid user_id")
	}

	input := usecases.CreateProposalInput{
		RequestBy:     uid,
		Requested:     req.Requested,
		ForExchangeID: req.ForExchangeID,
		Message:       req.Message,
	}

	if err := h.uc.CreateProposal(input); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusCreated)
}

func (h *Handler) GetProposalByID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	uid := c.Get("user_id").(string)

	proposal, err := h.uc.GetProposalByID(id)
	if err != nil || proposal == nil {
		return echo.NewHTTPError(http.StatusNotFound, "proposal not found")
	}

	if proposal.RequestBy.String() != uid && proposal.RequestTo.String() != uid {
		return echo.NewHTTPError(http.StatusForbidden, "you are not allowed to view this proposal")
	}

	return c.JSON(http.StatusOK, proposal)
}

func (h *Handler) GetAllProposals(c echo.Context) error {
	userID := c.Get("user_id").(string)

	uid, err := uuid.Parse(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid user_id")
	}

	proposals, err := h.uc.GetAllProposals(uid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, proposals)
}

func (h *Handler) AcceptProposal(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	userID := c.Get("user_id").(string)

	uid, err := uuid.Parse(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid user_id")
	}

	if err := h.uc.AcceptProposal(id, uid); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}
