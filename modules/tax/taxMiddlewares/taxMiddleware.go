package taxMiddlewares

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/montheankul-k/assessment-tax/config"
	"github.com/montheankul-k/assessment-tax/modules/servers"
	"github.com/montheankul-k/assessment-tax/modules/tax/taxHandlers"
	"net/http"
)

type ITaxMiddleware interface {
	ValidateCalculateTaxRequest(next echo.HandlerFunc) echo.HandlerFunc
}

type taxMiddleware struct {
	config     config.IConfig
	taxHandler taxHandlers.ITaxHandler
}

func Middleware(config config.IConfig, taxHandler taxHandlers.ITaxHandler) ITaxMiddleware {
	return &taxMiddleware{
		config:     config,
		taxHandler: taxHandler,
	}
}

func (m *taxMiddleware) ValidateCalculateTaxRequest(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req = taxHandlers.NewCalculateTaxRequest()
		err := c.Bind(&req)
		if err != nil {
			return servers.NewResponse(c).ResponseError(http.StatusBadRequest, err.Error())
		}

		if err := m.validateCalculateTaxRequest(req); err != nil {
			return servers.NewResponse(c).ResponseError(http.StatusBadRequest, err.Error())
		}

		return next(c)
	}
}

func (m *taxMiddleware) validateCalculateTaxRequest(req *taxHandlers.CalculateTaxRequest) error {
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

func (m *taxMiddleware) validateAllowance(allowance *taxHandlers.TaxAllowanceDetails) error {
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

func (m *taxMiddleware) findBaselineAmount(allowanceType string) (float64, float64, error) {
	return m.taxHandler.FindBaselineAllowance(allowanceType)
}

func (m *taxMiddleware) validateDonationAllowance(amount, minAmount, maxAmount float64) error {
	if amount < minAmount || amount > maxAmount {
		return errors.New("donation amount must be between 0 and 100000")
	}

	return nil
}

func (m *taxMiddleware) validateKReceiptAllowance(amount, minAmount, maxAmount float64) error {
	if amount < minAmount || amount > maxAmount {
		return fmt.Errorf("k-receipt amount must be between 0 and %.1f", maxAmount)
	}

	return nil
}
