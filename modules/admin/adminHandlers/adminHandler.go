package adminHandlers

import (
	"github.com/labstack/echo/v4"
	"github.com/montheankul-k/assessment-tax/config"
	"github.com/montheankul-k/assessment-tax/modules/admin"
	"github.com/montheankul-k/assessment-tax/modules/tax"
	"github.com/montheankul-k/assessment-tax/modules/tax/taxHandlers"
	"github.com/montheankul-k/assessment-tax/modules/tax/taxUsecases"
	"net/http"
)

type IAdminHandler interface {
	SetPersonalDeduction(c echo.Context) error
	SetKReceiptDeduction(c echo.Context) error
}

type adminHandler struct {
	config     config.IConfig
	taxUsecase taxUsecases.ITaxUsecase
}

func AdminHandler(config config.IConfig, taxUsecase taxUsecases.ITaxUsecase) IAdminHandler {
	return &adminHandler{
		config:     config,
		taxUsecase: taxUsecase,
	}
}

func (h *adminHandler) setDeduction(amount float64, allowanceType string) (float64, error) {
	newDeduction := tax.SetNewDeductionAmount{
		AllowanceFilter: tax.AllowanceFilter{
			AllowanceType: allowanceType,
		},
		NewDeductionAmount: amount,
	}

	result, err := h.taxUsecase.SetDeduction(&newDeduction)
	if err != nil {
		return 0, err
	}

	return result, nil
}

func (h *adminHandler) SetPersonalDeduction(c echo.Context) error {
	req, ok := c.Get("request").(*admin.DeductionAmount)
	if !ok {
		return taxHandlers.NewResponse(c).ResponseError(http.StatusInternalServerError, "failed to get request from context")
	}

	result, err := h.setDeduction(req.Amount, "personal")
	if err != nil {
		return taxHandlers.NewResponse(c).ResponseError(http.StatusInternalServerError, err.Error())
	}

	responseData := admin.DeductionAmount{
		Amount: result,
	}
	return taxHandlers.NewResponse(c).ResponseSuccess(http.StatusOK, responseData)
}

func (h *adminHandler) SetKReceiptDeduction(c echo.Context) error {
	req, ok := c.Get("request").(*admin.DeductionAmount)
	if !ok {
		return taxHandlers.NewResponse(c).ResponseError(http.StatusInternalServerError, "failed to get request from context")
	}

	result, err := h.setDeduction(req.Amount, "k-receipt")
	if err != nil {
		return taxHandlers.NewResponse(c).ResponseError(http.StatusInternalServerError, err.Error())
	}

	responseData := admin.DeductionAmount{
		Amount: result,
	}
	return taxHandlers.NewResponse(c).ResponseSuccess(http.StatusOK, responseData)
}
