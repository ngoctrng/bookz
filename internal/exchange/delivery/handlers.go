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

// CreateProposal godoc
// @Summary      Create a new exchange proposal
// @Description  Propose a book exchange with another user
// @Tags         exchange
// @Accept       json
// @Produce      json
// @Param        body  body      CreateProposalRequest  true  "Proposal info"
// @Success      201
// @Failure      400  {object}  echo.HTTPError
// @Failure      401  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Security     ApiKeyAuth
// @Router       /exchange/proposals [post]
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

// GetProposalByID godoc
// @Summary      Get proposal by ID
// @Description  Get details of a proposal by its ID
// @Tags         exchange
// @Produce      json
// @Param        id   path      int  true  "Proposal ID"
// @Success      200  {object}  exchange.Proposal
// @Failure      403  {object}  echo.HTTPError
// @Failure      404  {object}  echo.HTTPError
// @Router       /exchange/proposals/{id} [get]
// @Security     ApiKeyAuth
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

// GetAllProposals godoc
// @Summary      List all proposals for the user
// @Description  Get a list of all proposals related to the current user
// @Tags         exchange
// @Produce      json
// @Success      200  {array}   exchange.Proposal
// @Failure      401  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /exchange/proposals [get]
// @Security     ApiKeyAuth
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

// AcceptProposal godoc
// @Summary      Accept a proposal
// @Description  Accept a book exchange proposal
// @Tags         exchange
// @Produce      json
// @Param        id   path      int  true  "Proposal ID"
// @Success      200
// @Failure      401  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /exchange/proposals/{id}/accept [post]
// @Security     ApiKeyAuth
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
