package middlewareHandlers

import (
	"github.com/Montheankul-K/assessment-tax/modules/tax"
	"github.com/Montheankul-K/assessment-tax/modules/tax/taxHandlers"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupEchoContext() (ctx echo.Context, recorder *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()

	return e.NewContext(req, rec), rec
}

func TestMiddlewareHandler_ValidateCalculateTaxRequest(t *testing.T) {
	c, _ := setupEchoContext()
	handler := &middlewareHandler{}

	req := &taxHandlers.CalculateTaxRequest{
		TotalIncome: 500000,
		Wht:         50000,
		Allowances: []taxHandlers.TaxAllowanceDetails{
			{AllowanceType: "donation", Amount: 10000},
			{AllowanceType: "k-receipt", Amount: 20000},
		},
	}

	c.Set("request", req)
	err := handler.ValidateCalculateTaxRequest(func(c echo.Context) error {
		return nil
	})(c)

	assert.NoError(t, err)
	assert.NotNil(t, c.Get("request"))
}

func TestMiddlewareHandler_ChangeStructFormat(t *testing.T) {
	c, _ := setupEchoContext()
	handler := &middlewareHandler{}

	req := []tax.TaxFromCSV{
		{TotalIncome: 500000, Wht: 0, Donation: 0},
		{TotalIncome: 600000, Wht: 40000, Donation: 20000},
		{TotalIncome: 750000, Wht: 50000, Donation: 15000},
	}

	c.Set("request", req)
	err := handler.ChangeStructFormat(func(c echo.Context) error {
		return nil
	})(c)

	result := c.Get("request").([]taxHandlers.CalculateTaxRequest)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, len(req))
}

func TestMiddlewareHandler_ValidateTaxFromCSV(t *testing.T) {
	c, _ := setupEchoContext()
	handler := &middlewareHandler{}

	req := &[]taxHandlers.CalculateTaxRequest{
		{TotalIncome: 500000, Wht: 0, Allowances: []taxHandlers.TaxAllowanceDetails{
			{AllowanceType: "donation", Amount: 0},
		}},
		{TotalIncome: 600000, Wht: 40000, Allowances: []taxHandlers.TaxAllowanceDetails{
			{AllowanceType: "donation", Amount: 20000},
		}},
		{TotalIncome: 750000, Wht: 50000, Allowances: []taxHandlers.TaxAllowanceDetails{
			{AllowanceType: "donation", Amount: 15000},
		}},
	}

	c.Set("request", req)
	err := handler.ValidateTaxFromCSV(func(c echo.Context) error {
		return nil
	})(c)

	assert.NoError(t, err)
	assert.NotNil(t, c.Get("request"))
}
