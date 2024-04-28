package taxHandlers

import (
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

func (h *taxHandler) CalculateTax(c echo.Context) error {
	req, ok := c.Get("request").(*taxUsecases.CalculateTaxRequest)
	if !ok {
		return taxUsecases.NewResponse(c).ResponseError(http.StatusInternalServerError, "failed to get request from context")
	}

	result, err := h.taxUsecase.CalculateTaxWithoutWHT(req)
	if err != nil {
		return taxUsecases.NewResponse(c).ResponseError(http.StatusInternalServerError, err.Error())
	}

	taxLevel, err := h.taxUsecase.GetTaxLevelDetails(result)
	if err != nil {
		return taxUsecases.NewResponse(c).ResponseError(http.StatusInternalServerError, err.Error())
	}

	summaryTax := h.taxUsecase.DecreaseWHT(result, req.Wht)

	responseData := taxUsecases.TaxResponse{
		Tax:      result,
		TaxLevel: taxLevel,
		TotalTax: summaryTax,
	}

	if summaryTax < 0 {
		responseData.TotalTax = 0
		responseDataWithRefund := taxUsecases.TaxResponseWithRefund{
			TaxResponse: responseData,
			TaxRefund:   math.Abs(summaryTax),
		}

		return taxUsecases.NewResponse(c).ResponseSuccess(http.StatusOK, responseDataWithRefund)
	}

	return taxUsecases.NewResponse(c).ResponseSuccess(http.StatusOK, responseData)
}

func (h *taxHandler) CalculateTaxFromCSV(c echo.Context) error {
	req, ok := c.Get("request").([]taxUsecases.CalculateTaxRequest)
	if !ok {
		return taxUsecases.NewResponse(c).ResponseError(http.StatusInternalServerError, "failed to get request from context")
	}

	var responseData []interface{}
	for _, taxData := range req {
		result, err := h.taxUsecase.CalculateTaxWithoutWHT(&taxData)
		if err != nil {
			return taxUsecases.NewResponse(c).ResponseError(http.StatusInternalServerError, err.Error())
		}

		result = h.taxUsecase.DecreaseWHT(result, taxData.Wht)
		if result < 0 {
			taxResponse := taxUsecases.TaxCSVResponseWithRefund{
				TotalIncome: taxData.TotalIncome,
				Tax:         0,
				TaxRefund:   math.Abs(result),
			}

			responseData = append(responseData, taxResponse)
		}

		taxResponse := taxUsecases.TaxCSVResponse{
			TotalIncome: taxData.TotalIncome,
			Tax:         result,
		}

		responseData = append(responseData, taxResponse)
	}

	return taxUsecases.NewResponse(c).ResponseSuccess(http.StatusOK, responseData)
}
