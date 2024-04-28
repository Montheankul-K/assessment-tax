package adminHandlers

import (
	"github.com/Montheankul-K/assessment-tax/config"
	"github.com/Montheankul-K/assessment-tax/modules/admin"
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

func (m *MockTaxUsecase) CalculateTaxByTaxLevel(income float64) (float64, error) {
	args := m.Called(income)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockTaxUsecase) GetTaxLevel() ([]taxUsecases.EachTaxLevel, error) {
	args := m.Called()
	return args.Get(0).([]taxUsecases.EachTaxLevel), args.Error(1)
}

func (m *MockTaxUsecase) SetDeduction(req *tax.SetNewDeductionAmount) (float64, error) {
	args := m.Called(req)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockTaxUsecase) DecreasePersonalAllowance(totalIncome float64) (float64, error) {
	args := m.Called(totalIncome)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockTaxUsecase) DecreaseWHT(tax, wht float64) float64 {
	args := m.Called(tax, wht)
	return args.Get(0).(float64)
}

func (m *MockTaxUsecase) DecreaseAllowance(tax float64, allowances []taxUsecases.TaxAllowanceDetails) float64 {
	args := m.Called(tax, allowances)
	return args.Get(0).(float64)
}

func (m *MockTaxUsecase) ConstructTaxLevels(maxIncomeAmount float64, taxLevels []tax.TaxLevel) []taxUsecases.EachTaxLevel {
	args := m.Called(maxIncomeAmount, taxLevels)
	return args.Get(0).([]taxUsecases.EachTaxLevel)
}

func (m *MockTaxUsecase) SetValueToTaxLevel(taxLevels []taxUsecases.EachTaxLevel, tax float64) ([]taxUsecases.TaxLevelResponse, error) {
	args := m.Called(taxLevels, tax)
	return args.Get(0).([]taxUsecases.TaxLevelResponse), args.Error(1)
}

func (m *MockTaxUsecase) GetTaxLevelDetails(tax float64) ([]taxUsecases.TaxLevelResponse, error) {
	args := m.Called(tax)
	return args.Get(0).([]taxUsecases.TaxLevelResponse), args.Error(1)
}

func (m *MockTaxUsecase) CalculateTaxWithoutWHT(req *taxUsecases.CalculateTaxRequest) (float64, error) {
	args := m.Called(req)
	return args.Get(0).(float64), args.Error(1)
}

func setupEchoContext() (ctx echo.Context, recorder *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec), rec
}

func TestAdminHandler_SetPersonalDeduction(t *testing.T) {
	mockConfig := &MockConfig{}
	mockTaxUsecase := &MockTaxUsecase{}

	handler := &adminHandler{
		config:     mockConfig,
		taxUsecase: mockTaxUsecase,
	}

	c, rec := setupEchoContext()
	requestData := &admin.DeductionAmount{Amount: 70000.0}
	c.Set("request", requestData)

	mockTaxUsecase.On("SetDeduction", mock.Anything).Return(float64(70000.0), nil)
	err := handler.SetPersonalDeduction(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAdminHandler_SetKReceiptDeduction(t *testing.T) {
	mockConfig := &MockConfig{}
	mockTaxUsecase := &MockTaxUsecase{}

	handler := &adminHandler{
		config:     mockConfig,
		taxUsecase: mockTaxUsecase,
	}

	c, rec := setupEchoContext()
	requestData := &admin.DeductionAmount{Amount: 70000.0}
	c.Set("request", requestData)

	mockTaxUsecase.On("SetDeduction", mock.Anything).Return(float64(70000.0), nil)
	err := handler.SetKReceiptDeduction(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}
