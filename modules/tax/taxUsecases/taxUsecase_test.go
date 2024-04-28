package taxUsecases

import (
	"github.com/Montheankul-K/assessment-tax/modules/tax"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockTaxRepository struct{}

func (m *mockTaxRepository) FindBaselineAllowanceAmount(req *tax.AllowanceFilter) (float64, float64, error) {
	return 0.0, 100000.0, nil
}

func (m *mockTaxRepository) FindTaxPercentByIncome(req *tax.TaxLevelFilter) (float64, error) {
	return 35.0, nil
}

func (m *mockTaxRepository) FindMaxIncomeAndPercent() (float64, float64, error) {
	return 2000001.0, 35.0, nil
}

func (m *mockTaxRepository) GetTaxLevel() ([]tax.TaxLevel, error) {
	return []tax.TaxLevel{
		{MinIncome: 0.0, MaxIncome: 150000.0},
		{MinIncome: 150001.0, MaxIncome: 500000.0},
		{MinIncome: 500001.0, MaxIncome: 1000000.0},
		{MinIncome: 1000001.0, MaxIncome: 2000000.0},
		{MinIncome: 2000001.0, MaxIncome: 2000001.0},
	}, nil
}

func (m *mockTaxRepository) SetDeduction(req *tax.SetNewDeductionAmount) (float64, error) {
	return 70000.0, nil
}

func TestFindBaselineAllowance(t *testing.T) {
	usecase := taxUsecase{&mockTaxRepository{}}
	minAllowanceAmount, maxAllowanceAmount, err := usecase.taxRepository.FindBaselineAllowanceAmount(&tax.AllowanceFilter{})

	assert.NoError(t, err)
	assert.Equal(t, 0.0, minAllowanceAmount)
	assert.Equal(t, 100000.0, maxAllowanceAmount)
}

func TestFindTaxPercent(t *testing.T) {
	usecase := taxUsecase{&mockTaxRepository{}}
	taxPercent, err := usecase.taxRepository.FindTaxPercentByIncome(&tax.TaxLevelFilter{})

	assert.NoError(t, err)
	assert.Equal(t, 35.0, taxPercent)
}

func TestFindMaxIncomeAndPercent(t *testing.T) {
	usecase := taxUsecase{&mockTaxRepository{}}
	maxIncome, taxPercent, err := usecase.taxRepository.FindMaxIncomeAndPercent()

	assert.NoError(t, err)
	assert.Equal(t, 2000001.0, maxIncome)
	assert.Equal(t, 35.0, taxPercent)
}

func TestTaxUsecase_constructTaxLevels(t *testing.T) {
	usecase := taxUsecase{&mockTaxRepository{}}

	maxIncomeAmount := 2000001.0
	taxLevels := []tax.TaxLevel{
		{MinIncome: 0.0, MaxIncome: 150000.0},
		{MinIncome: 150001.0, MaxIncome: 500000.0},
		{MinIncome: 500001.0, MaxIncome: 1000000.0},
		{MinIncome: 1000001.0, MaxIncome: 2000000.0},
		{MinIncome: 2000001.0, MaxIncome: 2000001.0},
	}

	expectedTaxLevels := []EachTaxLevel{
		{MinMax: []float64{0.0, 150000.0}, Level: "0-150000", Tax: 0.0},
		{MinMax: []float64{150001.0, 500000.0}, Level: "150001-500000", Tax: 0.0},
		{MinMax: []float64{500001.0, 1000000.0}, Level: "500001-1000000", Tax: 0.0},
		{MinMax: []float64{1000001.0, 2000000.0}, Level: "1000001-2000000", Tax: 0.0},
		{MinMax: []float64{2000001.0, 2000001.0}, Level: "2000001 ขึ้นไป", Tax: 0.0},
	}

	result := usecase.ConstructTaxLevels(maxIncomeAmount, taxLevels)
	assert.Equal(t, expectedTaxLevels, result)
}

func TestGetTaxLevel(t *testing.T) {
	usecase := taxUsecase{&mockTaxRepository{}}

	expectedTaxLevels := []EachTaxLevel{
		{MinMax: []float64{0.0, 150000.0}, Level: "0-150000", Tax: 0.0},
		{MinMax: []float64{150001.0, 500000.0}, Level: "150001-500000", Tax: 0.0},
		{MinMax: []float64{500001.0, 1000000.0}, Level: "500001-1000000", Tax: 0.0},
		{MinMax: []float64{1000001.0, 2000000.0}, Level: "1000001-2000000", Tax: 0.0},
		{MinMax: []float64{2000001.0, 2000001.0}, Level: "2000001 ขึ้นไป", Tax: 0.0},
	}

	taxLevels, err := usecase.GetTaxLevel()
	assert.NoError(t, err)
	assert.Equal(t, expectedTaxLevels, taxLevels)
}

func TestSetDeduction(t *testing.T) {
	taxRepository := taxUsecase{&mockTaxRepository{}}
	result, err := taxRepository.SetDeduction(&tax.SetNewDeductionAmount{})

	assert.NoError(t, err)
	assert.Equal(t, 70000.0, result)
}

/*
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
	allowances := []taxUsecases.TaxAllowanceDetails{
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
*/
