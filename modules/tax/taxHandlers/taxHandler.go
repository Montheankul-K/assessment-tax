package taxHandlers

import (
	"fmt"
	"github.com/Montheankul-K/assessment-tax/config"
	"github.com/Montheankul-K/assessment-tax/modules/tax/taxUsecases"
	"github.com/labstack/echo/v4"
	"math"
	"net/http"
)

type ITaxHandler interface {
	CalculateTax(c echo.Context) error
	CalculateTaxFromCSV(c echo.Context) error
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

	taxPercent, err := h.taxUsecase.FindTaxPercent(income)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate tax: %v", err)
	}

	return income * (taxPercent / 100), nil
}

func (h *taxHandler) decreasePersonalAllowance(totalIncome float64) (float64, error) {
	_, maxAllowanceAmount, err := h.taxUsecase.FindBaseline("personal")
	if err != nil {
		return 0, fmt.Errorf("failed to decrease personal allowance")
	}

	return totalIncome - maxAllowanceAmount, nil
}

func (h *taxHandler) decreaseWHT(tax, wht float64) float64 {
	return tax - wht
}

func (h *taxHandler) decreaseAllowance(tax float64, allowances []TaxAllowanceDetails) float64 {
	result := tax
	for _, allowance := range allowances {
		result -= allowance.Amount
	}

	return result
}

func (h *taxHandler) getTaxLevel() ([]taxUsecases.EachTaxLevel, error) {
	result, err := h.taxUsecase.GetTaxLevel()
	if err != nil {
		return result, fmt.Errorf("failed to get tax level")
	}

	return result, nil
}

type TaxLevelResponse struct {
	Level string  `json:"level"`
	Tax   float64 `json:"tax"`
}

func (h *taxHandler) setValueToTaxLevel(taxLevels []taxUsecases.EachTaxLevel, tax float64) ([]TaxLevelResponse, error) {
	result := make([]TaxLevelResponse, 0, len(taxLevels))

	for _, level := range taxLevels {
		levelMinIncome, levelMaxIncome := level.MinMax[0], level.MinMax[1]
		if tax >= levelMinIncome && tax <= levelMaxIncome || levelMinIncome == levelMaxIncome && tax >= levelMaxIncome {
			result = append(result, TaxLevelResponse{
				Level: level.Level,
				Tax:   h.roundToOneDecimal(tax),
			})
			continue
		}

		result = append(result, TaxLevelResponse{
			Level: level.Level,
			Tax:   h.roundToOneDecimal(level.Tax),
		})
	}

	return result, nil
}

func (h *taxHandler) roundToOneDecimal(num float64) float64 {
	return math.Round(num*10) / 10
}

func (h *taxHandler) getTaxLevelDetails(tax float64) ([]TaxLevelResponse, error) {
	taxLevel, err := h.getTaxLevel()
	if err != nil {
		return nil, fmt.Errorf("failed to get tax level")
	}

	result, err := h.setValueToTaxLevel(taxLevel, tax)
	if err != nil {
		return nil, fmt.Errorf("failed to get tax level")
	}

	return result, nil
}

func (h *taxHandler) calculateTaxWithoutWHT(req *CalculateTaxRequest) (float64, error) {
	result, err := h.decreasePersonalAllowance(req.TotalIncome)
	if err != nil {
		return 0, err
	}
	result = math.Max(0, result)

	result = h.decreaseAllowance(result, req.Allowances)
	result = math.Max(0, result)

	result, err = h.calculateTaxByTaxLevel(result)
	if err != nil {
		return 0, err
	}

	return result, nil
}

func (h *taxHandler) CalculateTax(c echo.Context) error {
	req, ok := c.Get("request").(*CalculateTaxRequest)
	if !ok {
		return NewResponse(c).ResponseError(http.StatusInternalServerError, "failed to get request from context")
	}

	result, err := h.calculateTaxWithoutWHT(req)
	if err != nil {
		return NewResponse(c).ResponseError(http.StatusInternalServerError, err.Error())
	}

	taxLevel, err := h.getTaxLevelDetails(result)
	if err != nil {
		return NewResponse(c).ResponseError(http.StatusInternalServerError, err.Error())
	}

	summaryTax := h.decreaseWHT(result, req.Wht)

	responseData := TaxResponse{
		Tax:      h.roundToOneDecimal(result),
		TaxLevel: taxLevel,
		TotalTax: h.roundToOneDecimal(summaryTax),
	}

	if summaryTax < 0 {
		responseData.TotalTax = 0
		responseDataWithRefund := TaxResponseWithRefund{
			TaxResponse: responseData,
			TaxRefund:   h.roundToOneDecimal(math.Abs(summaryTax)),
		}

		return NewResponse(c).ResponseSuccess(http.StatusOK, responseDataWithRefund)
	}

	return NewResponse(c).ResponseSuccess(http.StatusOK, responseData)
}

func (h *taxHandler) CalculateTaxFromCSV(c echo.Context) error {
	req, ok := c.Get("request").([]CalculateTaxRequest)
	if !ok {
		return NewResponse(c).ResponseError(http.StatusInternalServerError, "failed to get request from context")
	}

	var responseData []interface{}
	for _, taxData := range req {
		result, err := h.calculateTaxWithoutWHT(&taxData)
		if err != nil {
			return NewResponse(c).ResponseError(http.StatusInternalServerError, err.Error())
		}

		result = h.decreaseWHT(result, taxData.Wht)
		if result < 0 {
			taxResponse := TaxCSVResponseWithRefund{
				TotalIncome: h.roundToOneDecimal(taxData.TotalIncome),
				Tax:         0,
				TaxRefund:   h.roundToOneDecimal(math.Abs(result)),
			}

			responseData = append(responseData, taxResponse)
		}

		taxResponse := TaxCSVResponse{
			TotalIncome: h.roundToOneDecimal(taxData.TotalIncome),
			Tax:         result,
		}

		responseData = append(responseData, taxResponse)
	}

	return NewResponse(c).ResponseSuccess(http.StatusOK, responseData)
}
