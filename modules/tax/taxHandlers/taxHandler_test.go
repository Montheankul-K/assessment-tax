package taxHandlers

import (
	"errors"
	"github.com/Montheankul-K/assessment-tax/config"
	"github.com/Montheankul-K/assessment-tax/modules/tax"
	"github.com/Montheankul-K/assessment-tax/modules/tax/taxUsecases"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockConfig struct {
	mockAppConfig *MockAppConfig
	mockDBConfig  *MockDBConfig
	mockAdminAuth *MockAdminAuth
}

type MockAppConfig struct{}
type MockDBConfig struct{}
type MockAdminAuth struct{}

func (m *MockConfig) App() config.IAppConfig {
	return m.mockAppConfig
}

func (m *MockConfig) DB() config.IDBConfig {
	return m.mockDBConfig
}

func (m *MockConfig) AdminAuth() config.IAdminAuth {
	return m.mockAdminAuth
}

func (m *MockAppConfig) Name() string {
	return "name"
}

func (m *MockAppConfig) Port() string {
	return "port"
}

func (m *MockAppConfig) Version() string {
	return "version"
}

func (m *MockDBConfig) Url() string {
	return "url"
}

func (m *MockAdminAuth) Username() string {
	return "username"
}

func (m *MockAdminAuth) Password() string {
	return "password"
}

type MockTaxUsecase struct {
	mock.Mock
}

func (m *MockTaxUsecase) FindBaseline(allowanceType string) (float64, float64, error) {
	args := m.Called(allowanceType)
	return args.Get(0).(float64), args.Get(1).(float64), args.Error(2)
}

func (m *MockTaxUsecase) FindTaxPercent(totalIncome float64) (float64, error) {
	args := m.Called(totalIncome)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockTaxUsecase) FindMaxIncomeAndPercent() (float64, float64, error) {
	args := m.Called()
	return args.Get(0).(float64), args.Get(1).(float64), args.Error(2)
}

func (m *MockTaxUsecase) GetTaxLevel() ([]taxUsecases.EachTaxLevel, error) {
	args := m.Called()
	return args.Get(0).([]taxUsecases.EachTaxLevel), args.Error(1)
}

func (m *MockTaxUsecase) SetDeduction(req *tax.SetNewDeductionAmount) (float64, error) {
	args := m.Called(req)
	return args.Get(0).(float64), args.Error(1)
}

type CustomTaxHandlerMock struct {
	CalculateTaxWithoutWHTFunc func(req *CalculateTaxRequest) (float64, error)
	GetTaxLevelDetailsFunc     func(tax float64) ([]TaxLevelResponse, error)
	DecreaseWHTFunc            func(tax, wht float64) float64
	RoundToOneDecimalFunc      func(num float64) float64
}

func (m *CustomTaxHandlerMock) CalculateTaxWithoutWHT(req *CalculateTaxRequest) (float64, error) {
	if m.CalculateTaxWithoutWHTFunc != nil {
		return m.CalculateTaxWithoutWHTFunc(req)
	}
	return 0, nil
}

func (m *CustomTaxHandlerMock) GetTaxLevelDetails(tax float64) ([]TaxLevelResponse, error) {
	if m.GetTaxLevelDetailsFunc != nil {
		return m.GetTaxLevelDetailsFunc(tax)
	}
	return nil, nil
}

func (m *CustomTaxHandlerMock) DecreaseWHT(tax, wht float64) float64 {
	if m.DecreaseWHTFunc != nil {
		return m.DecreaseWHTFunc(tax, wht)
	}
	return 0
}

func (m *CustomTaxHandlerMock) RoundToOneDecimal(num float64) float64 {
	if m.RoundToOneDecimalFunc != nil {
		return m.RoundToOneDecimalFunc(num)
	}
	return 0
}

func setupEchoContext() (ctx echo.Context, recorder *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec), rec
}

func TestTaxHandler_findMaxIncomeAndPercent(t *testing.T) {
	cfg := &MockConfig{}
	usecase := &MockTaxUsecase{}
	handler := taxHandler{cfg, usecase}

	expectedMaxIncome := 2000001.0
	expectedMaxPercent := 35.0
	usecase.On("FindMaxIncomeAndPercent", mock.Anything).Return(expectedMaxIncome, expectedMaxPercent, nil)

	maxIncome, maxPercent, err := handler.findMaxIncomeAndPercent()

	assert.NoError(t, err)
	assert.Equal(t, expectedMaxIncome, maxIncome)
	assert.Equal(t, expectedMaxPercent, maxPercent)
}

func TestTaxHandler_findMaxIncomeAndPercent_Error(t *testing.T) {
	cfg := &MockConfig{}
	usecase := &MockTaxUsecase{}
	handler := taxHandler{cfg, usecase}

	expectErr := errors.New("failed to find max income and percent: error")
	usecase.On("FindMaxIncomeAndPercent", mock.Anything).Return(0.0, 0.0, errors.New("error"))

	_, _, err := handler.findMaxIncomeAndPercent()

	assert.Error(t, err)
	assert.Equal(t, expectErr.Error(), err.Error())
}

func TestTaxHandler_calculateTaxByTaxLevel(t *testing.T) {
	cfg := &MockConfig{}
	usecase := &MockTaxUsecase{}
	handler := taxHandler{cfg, usecase}

	maxIncome := 2000001.0
	maxPercent := 35.0

	income := 500000.0
	expectedTax := income * (10.0 / 100.0)

	usecase.On("FindMaxIncomeAndPercent", mock.Anything).Return(maxIncome, maxPercent, nil)
	usecase.On("FindTaxPercent", mock.Anything).Return(10.0, nil)

	resultTax, err := handler.calculateTaxByTaxLevel(income)

	assert.NoError(t, err)
	assert.Equal(t, expectedTax, resultTax)

	income = 3000000.0
	expectedTax = income * (35.0 / 100.0)

	resultTax, err = handler.calculateTaxByTaxLevel(income)

	assert.NoError(t, err)
	assert.Equal(t, expectedTax, resultTax)
}

func TestTaxHandler_decreaseWHT(t *testing.T) {
	handler := &taxHandler{}

	texWithoutWHT := 10000.0
	wht := 5000.0
	expectedTotalTax := texWithoutWHT - wht

	result := handler.decreaseWHT(texWithoutWHT, wht)
	if result != expectedTotalTax {
		t.Errorf("expected: %f but got: %f", expectedTotalTax, result)
	}
}

func TestTaxHandler_decreaseAllowance(t *testing.T) {
	handler := &taxHandler{}

	taxBeforeDecreaseAllowance := 100000.0
	allowances := []TaxAllowanceDetails{
		{AllowanceType: "donation", Amount: 20000.0},
		{AllowanceType: "k-receipt", Amount: 10000.0},
	}
	expectedTax := taxBeforeDecreaseAllowance - (allowances[0].Amount + allowances[1].Amount)

	result := handler.decreaseAllowance(taxBeforeDecreaseAllowance, allowances)

	if result != expectedTax {
		t.Errorf("expected: %f but got: %f", expectedTax, result)
	}
}

func TestTaxHandler_getTaxLevel(t *testing.T) {
	cfg := &MockConfig{}
	usecase := &MockTaxUsecase{}
	handler := taxHandler{cfg, usecase}

	mockResult := []taxUsecases.EachTaxLevel{
		{MinMax: []float64{0.0, 150000.0}, Level: "0-150000", Tax: 0.0},
		{MinMax: []float64{150001.0, 500000.0}, Level: "150001-500000", Tax: 0.0},
		{MinMax: []float64{500001.0, 1000000.0}, Level: "500001-1000000", Tax: 0.0},
		{MinMax: []float64{1000001.0, 2000000.0}, Level: "1000001-2000000", Tax: 0.0},
		{MinMax: []float64{2000001.0, 2000001.0}, Level: "2000001 ขึ้นไป", Tax: 0.0},
	}
	usecase.On("GetTaxLevel").Return(mockResult, nil)

	result, err := handler.getTaxLevel()

	assert.NoError(t, err)
	assert.Equal(t, mockResult, result)
}

func TestTaxHandler_getTaxLevel_Error(t *testing.T) {
	cfg := &MockConfig{}
	usecase := &MockTaxUsecase{}
	handler := taxHandler{cfg, usecase}

	mockResult := []taxUsecases.EachTaxLevel{}

	expectErr := errors.New("failed to get tax level")
	usecase.On("GetTaxLevel").Return(mockResult, expectErr)

	_, err := handler.getTaxLevel()

	assert.Error(t, err)
	assert.Equal(t, expectErr.Error(), err.Error())
}

func TestTaxHandler_setValueToTaxLevel(t *testing.T) {
	handler := &taxHandler{}

	taxLevels := []taxUsecases.EachTaxLevel{
		{MinMax: []float64{0.0, 150000.0}, Level: "0-150000", Tax: 0.0},
		{MinMax: []float64{150001.0, 500000.0}, Level: "150001-500000", Tax: 0.0},
		{MinMax: []float64{500001.0, 1000000.0}, Level: "500001-1000000", Tax: 0.0},
		{MinMax: []float64{1000001.0, 2000000.0}, Level: "1000001-2000000", Tax: 0.0},
		{MinMax: []float64{2000001.0, 2000001.0}, Level: "2000001 ขึ้นไป", Tax: 0.0},
	}

	income := 100000.0
	expected := []TaxLevelResponse{
		{Level: "0-150000", Tax: income},
		{Level: "150001-500000", Tax: 0.0},
		{Level: "500001-1000000", Tax: 0.0},
		{Level: "1000001-2000000", Tax: 0.0},
		{Level: "2000001 ขึ้นไป", Tax: 0.0},
	}

	result, err := handler.setValueToTaxLevel(taxLevels, income)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestTaxHandler_roundToOneDecimal(t *testing.T) {
	handler := &taxHandler{}

	number := 100000.00
	expected := 100000.0

	result := handler.roundToOneDecimal(number)

	if result != expected {
		t.Errorf("Expected result: %.1f but got: %.1f", expected, result)
	}
}
