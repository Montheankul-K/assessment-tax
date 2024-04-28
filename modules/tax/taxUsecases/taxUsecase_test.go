package taxUsecases

import (
	"github.com/montheankul-k/assessment-tax/modules/tax"
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
	taxRepository := taxUsecase{&mockTaxRepository{}}
	minAllowanceAmount, maxAllowanceAmount, err := taxRepository.FindBaselineAllowance(&tax.AllowanceFilter{})

	assert.NoError(t, err)
	assert.Equal(t, 0.0, minAllowanceAmount)
	assert.Equal(t, 100000.0, maxAllowanceAmount)
}

func TestFindTaxPercent(t *testing.T) {
	taxRepository := taxUsecase{&mockTaxRepository{}}
	taxPercent, err := taxRepository.FindTaxPercent(&tax.TaxLevelFilter{})

	assert.NoError(t, err)
	assert.Equal(t, 35.0, taxPercent)
}

func TestFindMaxIncomeAndPercent(t *testing.T) {
	taxRepository := taxUsecase{&mockTaxRepository{}}
	maxIncome, taxPercent, err := taxRepository.FindMaxIncomeAndPercent()

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

	result := usecase.constructTaxLevels(maxIncomeAmount, taxLevels)
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
