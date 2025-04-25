package handler

import (
	"net/http"

	"github.com/alfaa19/service-account-test/internal/models/dto"
	"github.com/alfaa19/service-account-test/internal/service"
	logger "github.com/alfaa19/service-account-test/pkg/logrus"
	"github.com/labstack/echo/v4"
)

type accountHandler struct {
	service service.Service
	log     *logger.CustomLogger
}

type AccountHandler interface {
	CreateAccount(ctx echo.Context) error
	GetSaldo(ctx echo.Context) error
	Withdraw(ctx echo.Context) error
	Deposit(ctx echo.Context) error
}

func NewAccountHandler(service service.Service, log *logger.CustomLogger) *accountHandler {
	return &accountHandler{
		service: service,
		log:     log,
	}
}

func (h *accountHandler) CreateAccount(c echo.Context) error {
	reqAccount := &dto.AccountRegistration{}
	if err := c.Bind(reqAccount); err != nil {
		h.log.Error("Failed to bind account: ", err)
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Remark: err.Error()})
	}

	createdAccount, err := h.service.CreateAccount(c.Request().Context(), reqAccount)

	if err != nil {
		h.log.Error("Failed to create account: ", err)
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Remark: err.Error()})
	}

	return c.JSON(http.StatusCreated, dto.AccountResponse{NoRekening: createdAccount.AccountNumber})
}

func (h *accountHandler) GetSaldo(c echo.Context) error {
	noRekening := c.Param("noRekening")
	account, err := h.service.GetAccountByNoRekening(c.Request().Context(), noRekening)
	if err != nil {
		h.log.Error("Failed to get saldo: ", err)
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Remark: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.BalanceResponse{
		Saldo: account.Balance})
}

func (h *accountHandler) Withdraw(c echo.Context) error {
	req := &dto.WithdrawDepositRequest{}
	if err := c.Bind(&req); err != nil {
		h.log.Error("Failed to bind withdraw request: ", err)
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{})
	}

	if err := h.service.UpdateBalanceWithdraw(c.Request().Context(), req.NoRekening, req.Amount); err != nil {
		h.log.Error("Failed to withdraw: ", err)
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Remark: err.Error(),
		})
	}

	account, err := h.service.GetAccountByNoRekening(c.Request().Context(), req.NoRekening)
	if err != nil {
		h.log.Error("Failed to get saldo: ", err)
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Remark: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.BalanceResponse{
		Saldo: account.Balance})
}

func (h *accountHandler) Deposit(c echo.Context) error {
	req := &dto.WithdrawDepositRequest{}
	if err := c.Bind(&req); err != nil {
		h.log.Error("Failed to bind deposit request: ", err)
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{})
	}

	if err := h.service.UpdateBalanceDeposit(c.Request().Context(), req.NoRekening, req.Amount); err != nil {
		h.log.Error("Failed to deposit: ", err)
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Remark: err.Error(),
		})
	}

	account, err := h.service.GetAccountByNoRekening(c.Request().Context(), req.NoRekening)
	if err != nil {
		h.log.Error("Failed to get saldo: ", err)
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Remark: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.BalanceResponse{
		Saldo: account.Balance})
}
