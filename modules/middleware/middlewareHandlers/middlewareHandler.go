package middlewareHandlers

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/montheankul-k/assessment-tax/config"
	"github.com/montheankul-k/assessment-tax/modules/admin"
	"github.com/montheankul-k/assessment-tax/modules/tax"
	"github.com/montheankul-k/assessment-tax/modules/tax/taxHandlers"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
)

type IMiddlewareHandler interface {
	ValidateCalculateTaxRequest(next echo.HandlerFunc) echo.HandlerFunc
	ValidateSetDeductionRequest(next echo.HandlerFunc) echo.HandlerFunc
	GetDataFromTaxCSV(next echo.HandlerFunc) echo.HandlerFunc
	ChangeStructFormat(next echo.HandlerFunc) echo.HandlerFunc
	ValidateTaxFromCSV(next echo.HandlerFunc) echo.HandlerFunc
}

type middlewareHandler struct {
	config     config.IConfig
	taxHandler taxHandlers.ITaxHandler
}

func MiddlewareHandler(config config.IConfig, taxHandler taxHandlers.ITaxHandler) IMiddlewareHandler {
	return &middlewareHandler{
		config:     config,
		taxHandler: taxHandler,
	}
}

func (m *middlewareHandler) ValidateCalculateTaxRequest(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req = taxHandlers.NewCalculateTaxRequest()
		err := c.Bind(&req)
		if err != nil {
			return taxHandlers.NewResponse(c).ResponseError(http.StatusBadRequest, err.Error())
		}

		if err := m.validateCalculateTaxRequest(req); err != nil {
			return taxHandlers.NewResponse(c).ResponseError(http.StatusBadRequest, err.Error())
		}

		c.Set("request", req)
		return next(c)
	}
}

func (m *middlewareHandler) validateCalculateTaxRequest(req *taxHandlers.CalculateTaxRequest) error {
	if req.TotalIncome <= 0 {
		return errors.New("total income must be gather than zero")
	}

	if req.Wht < 0 || req.Wht > req.TotalIncome {
		return errors.New("wht must be between 0 and total income")
	}

	for _, allowance := range req.Allowances {
		if err := m.validateAllowance(&allowance); err != nil {
			return err
		}
	}

	return nil
}

func (m *middlewareHandler) validateAllowance(allowance *taxHandlers.TaxAllowanceDetails) error {
	minAmount, maxAmount, err := m.findBaselineAmount(allowance.AllowanceType)
	if err != nil {
		return err
	}

	switch allowance.AllowanceType {
	case "donation":
		return m.validateDonationAllowance(allowance.Amount, minAmount, maxAmount)
	case "k-receipt":
		return m.validateKReceiptAllowance(allowance.Amount, minAmount, maxAmount)
	default:
		return errors.New("invalid allowance type")
	}
}

func (m *middlewareHandler) findBaselineAmount(allowanceType string) (float64, float64, error) {
	return m.taxHandler.FindBaseline(allowanceType)
}

func (m *middlewareHandler) validateDonationAllowance(amount, minAmount, maxAmount float64) error {
	if amount < minAmount || amount > maxAmount {
		return errors.New("donation amount must be between 0 and 100000")
	}

	return nil
}

func (m *middlewareHandler) validateKReceiptAllowance(amount, minAmount, maxAmount float64) error {
	if amount < minAmount || amount > maxAmount {
		return fmt.Errorf("k-receipt amount must be between 0 and %.1f", maxAmount)
	}

	return nil
}

func (m *middlewareHandler) ValidateSetDeductionRequest(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req *admin.DeductionAmount
		err := c.Bind(&req)
		if err != nil {
			return taxHandlers.NewResponse(c).ResponseError(http.StatusBadRequest, err.Error())
		}

		if req.Amount > 100000 {
			return taxHandlers.NewResponse(c).ResponseError(http.StatusBadRequest, "req shouldn't be gather than 100000")
		}

		c.Set("request", req)
		return next(c)
	}
}

func (m *middlewareHandler) GetDataFromTaxCSV(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		file, err := c.FormFile("taxes")
		if err != nil {
			return taxHandlers.NewResponse(c).ResponseError(http.StatusBadRequest, err.Error())
		}

		src, err := file.Open()
		if err != nil {
			return taxHandlers.NewResponse(c).ResponseError(http.StatusInternalServerError, "could not open file")
		}
		defer func(src multipart.File) {
			err := src.Close()
			if err != nil {
			}
		}(src)

		reader := csv.NewReader(src)
		var req []tax.TaxFromCSV

		if _, err := reader.Read(); err != nil {
			if err == io.EOF {
				return taxHandlers.NewResponse(c).ResponseError(http.StatusBadRequest, "empty file")
			}
			return taxHandlers.NewResponse(c).ResponseError(http.StatusInternalServerError, "could not read file")
		}

		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return taxHandlers.NewResponse(c).ResponseError(http.StatusInternalServerError, "could not read file")
			}

			totalIncome, err := strconv.ParseFloat(record[0], 64)
			if err != nil {
				return taxHandlers.NewResponse(c).ResponseError(http.StatusInternalServerError, "could not parse total income")
			}

			wht, err := strconv.ParseFloat(record[1], 64)
			if err != nil {
				return taxHandlers.NewResponse(c).ResponseError(http.StatusInternalServerError, "could not parse total wht")
			}

			donation, err := strconv.ParseFloat(record[2], 64)
			if err != nil {
				return taxHandlers.NewResponse(c).ResponseError(http.StatusInternalServerError, "could not parse donation")
			}

			req = append(req, tax.TaxFromCSV{
				TotalIncome: totalIncome,
				Wht:         wht,
				Donation:    donation,
			})
		}

		c.Set("request", req)
		return next(c)
	}
}

func (m *middlewareHandler) ChangeStructFormat(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req, ok := c.Get("request").([]tax.TaxFromCSV)
		if !ok {
			return taxHandlers.NewResponse(c).ResponseError(http.StatusInternalServerError, "failed to get request from context")
		}

		var result []taxHandlers.CalculateTaxRequest
		for _, taxData := range req {
			calculateTaxRequest := taxHandlers.CalculateTaxRequest{
				TotalIncome: taxData.TotalIncome,
				Wht:         taxData.Wht,
				Allowances: []taxHandlers.TaxAllowanceDetails{
					{AllowanceType: "donation", Amount: taxData.Donation},
				},
			}

			result = append(result, calculateTaxRequest)
		}

		c.Set("request", result)
		return next(c)
	}
}

func (m *middlewareHandler) ValidateTaxFromCSV(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req, ok := c.Get("request").([]taxHandlers.CalculateTaxRequest)
		if !ok {
			return taxHandlers.NewResponse(c).ResponseError(http.StatusInternalServerError, "failed to get request from context")
		}

		for _, taxData := range req {
			if err := m.validateCalculateTaxRequest(&taxData); err != nil {
				return taxHandlers.NewResponse(c).ResponseError(http.StatusBadRequest, err.Error())
			}
		}

		c.Set("request", req)
		return next(c)
	}
}
