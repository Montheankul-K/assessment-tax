package taxHandlers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/montheankul-k/assessment-tax/config"
	"github.com/montheankul-k/assessment-tax/modules/tax"
	"github.com/montheankul-k/assessment-tax/modules/tax/taxUsecases"
	"net/http"
)

type ITaxHandler interface {
	FindBaseline(allowanceType string) (float64, float64, error)
	FindTaxPercent(totalIncome float64) (float64, error)
	CalculateTax(c echo.Context) error
}

type taxHandler struct {
	config     config.IConfig
	taxUsecase taxUsecases.ITaxUsecase
}

func TaxHandler(config config.IConfig, taxUsecase taxUsecases.ITaxUsecase) ITaxHandler {
	return &taxHandler{
		config:     config,
		taxUsecase: taxUsecase,
	}
}

func (h *taxHandler) FindBaseline(allowanceType string) (float64, float64, error) {
	req := tax.AllowanceFilter{
		AllowanceType: allowanceType,
	}

	minAllowanceAmount, maxAllowanceAmount, err := h.taxUsecase.FindBaselineAllowance(&req)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to find baseline allowance: %v", err)
	}

	return minAllowanceAmount, maxAllowanceAmount, nil
}

func (h *taxHandler) FindTaxPercent(totalIncome float64) (float64, error) {
	req := tax.TaxLevelFilter{
		Income: totalIncome,
	}

	taxPercent, err := h.taxUsecase.FindTaxPercent(&req)
	if err != nil {
		return 0, fmt.Errorf("failed to find tax percent: %v", err)
	}

	return taxPercent, nil
}

func (h *taxHandler) findMaxIncomeAndPercent() (float64, float64, error) {
	maxIncome, maxPercent, err := h.taxUsecase.FindMaxIncomeAndPercent()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to find max income and percent: %v", err)
	}

	return maxIncome, maxPercent, nil
}

func (h *taxHandler) calculateTaxByTaxLevel(income float64) (float64, error) {
	maxIncome, maxPercent, err := h.findMaxIncomeAndPercent()
	if err != nil {
		return 0, fmt.Errorf("failed to calculate tax: %v", err)
	}

	if income > maxIncome {
		return income * (maxPercent / 100), nil
	}

	taxPercent, err := h.FindTaxPercent(income)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate tax: %v", err)
	}

	return income * (taxPercent / 100), nil
}

func (h *taxHandler) decreasePersonalAllowance(totalIncome float64) (float64, error) {
	minAllowanceAmount, _, err := h.FindBaseline("personal")
	if err != nil {
		return 0, fmt.Errorf("failed to decrease personal allowance")
	}

	return totalIncome - minAllowanceAmount, nil
}

func (h *taxHandler) CalculateTax(c echo.Context) error {
	req, ok := c.Get("request").(*CalculateTaxRequest)
	if !ok {
		return NewResponse(c).ResponseError(http.StatusInternalServerError, "failed to get request from context")
	}

	incomeAfterDecrease, err := h.decreasePersonalAllowance(req.TotalIncome)
	if err != nil {
		return NewResponse(c).ResponseError(http.StatusInternalServerError, err.Error())
	}

	taxAmount, err := h.calculateTaxByTaxLevel(incomeAfterDecrease)
	if err != nil {
		return NewResponse(c).ResponseError(http.StatusInternalServerError, err.Error())
	}

	responseData := map[string]float64{"tax": taxAmount}
	return NewResponse(c).ResponseSuccess(http.StatusOK, responseData)
}
