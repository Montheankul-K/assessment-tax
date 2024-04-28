package taxUsecases

import (
	"fmt"
	"github.com/Montheankul-K/assessment-tax/modules/tax"
	"github.com/Montheankul-K/assessment-tax/modules/tax/taxRepositories"
	"math"
)

type ITaxUsecase interface {
	FindBaseline(allowanceType string) (float64, float64, error)
	FindTaxPercent(totalIncome float64) (float64, error)
	FindMaxIncomeAndPercent() (float64, float64, error)
	CalculateTaxByTaxLevel(income float64) (float64, error)
	GetTaxLevel() ([]EachTaxLevel, error)
	SetDeduction(req *tax.SetNewDeductionAmount) (float64, error)
	DecreasePersonalAllowance(totalIncome float64) (float64, error)
	DecreaseWHT(tax, wht float64) float64
	DecreaseAllowance(tax float64, allowances []TaxAllowanceDetails) float64
	ConstructTaxLevels(maxIncomeAmount float64, taxLevels []tax.TaxLevel) []EachTaxLevel
	SetValueToTaxLevel(taxLevels []EachTaxLevel, tax float64) ([]TaxLevelResponse, error)
	GetTaxLevelDetails(tax float64) ([]TaxLevelResponse, error)
	CalculateTaxWithoutWHT(req *CalculateTaxRequest) (float64, error)
}

type taxUsecase struct {
	taxRepository taxRepositories.ITaxRepository
}

type TaxLevelResponse struct {
	Level string  `json:"level"`
	Tax   float64 `json:"tax"`
}

func TaxUsecase(taxRepository taxRepositories.ITaxRepository) ITaxUsecase {
	return &taxUsecase{
		taxRepository: taxRepository,
	}
}

func (u *taxUsecase) FindBaseline(allowanceType string) (float64, float64, error) {
	req := tax.AllowanceFilter{
		AllowanceType: allowanceType,
	}

	minAllowanceAmount, maxAllowanceAmount, err := u.taxRepository.FindBaselineAllowanceAmount(&req)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to find baseline allowance: %v", err)
	}

	return minAllowanceAmount, maxAllowanceAmount, nil
}

func (u *taxUsecase) FindTaxPercent(totalIncome float64) (float64, error) {
	req := tax.TaxLevelFilter{
		Income: totalIncome,
	}

	taxPercent, err := u.taxRepository.FindTaxPercentByIncome(&req)
	if err != nil {
		return 0, fmt.Errorf("failed to find tax percent: %v", err)
	}

	return taxPercent, nil
}

func (u *taxUsecase) FindMaxIncomeAndPercent() (float64, float64, error) {
	maxIncome, taxPercent, err := u.taxRepository.FindMaxIncomeAndPercent()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to find max income and percent: %v", err)
	}

	return maxIncome, taxPercent, nil
}

type EachTaxLevel struct {
	MinMax []float64
	Level  string
	Tax    float64
}

func (u *taxUsecase) ConstructTaxLevels(maxIncomeAmount float64, taxLevels []tax.TaxLevel) []EachTaxLevel {
	newTaxLevel := make([]EachTaxLevel, 0, len(taxLevels))
	for _, level := range taxLevels {
		var levelDesc string
		if level.MinIncome == maxIncomeAmount {
			levelDesc = fmt.Sprintf("%d ขึ้นไป", int(level.MinIncome))
		} else {
			levelDesc = fmt.Sprintf("%d-%d", int(level.MinIncome), int(level.MaxIncome))
		}

		newTaxLevel = append(newTaxLevel, EachTaxLevel{
			MinMax: []float64{level.MinIncome, level.MaxIncome},
			Level:  levelDesc,
			Tax:    0.0,
		})
	}

	return newTaxLevel
}

func (u *taxUsecase) GetTaxLevel() ([]EachTaxLevel, error) {
	maxIncomeAmount, _, err := u.FindMaxIncomeAndPercent()
	if err != nil {
		return nil, fmt.Errorf("failed to get max income and percent: %v", err)
	}

	taxLevels, err := u.taxRepository.GetTaxLevel()
	if err != nil {
		return nil, fmt.Errorf("failed to get tax level: %v", err)
	}

	return u.ConstructTaxLevels(maxIncomeAmount, taxLevels), nil
}

func (u *taxUsecase) SetDeduction(req *tax.SetNewDeductionAmount) (float64, error) {
	result, err := u.taxRepository.SetDeduction(req)
	if err != nil {
		return 0, err
	}

	return result, nil
}

func (u *taxUsecase) DecreasePersonalAllowance(totalIncome float64) (float64, error) {
	_, maxAllowanceAmount, err := u.FindBaseline("personal")
	if err != nil {
		return 0, fmt.Errorf("failed to decrease personal allowance")
	}

	return totalIncome - maxAllowanceAmount, nil
}

func (u *taxUsecase) DecreaseWHT(tax, wht float64) float64 {
	return tax - wht
}

func (u *taxUsecase) DecreaseAllowance(tax float64, allowances []TaxAllowanceDetails) float64 {
	result := tax
	for _, allowance := range allowances {
		result -= allowance.Amount
	}

	return result
}

func (u *taxUsecase) CalculateTaxByTaxLevel(income float64) (float64, error) {
	maxIncome, maxPercent, err := u.FindMaxIncomeAndPercent()
	if err != nil {
		return 0, fmt.Errorf("failed to calculate tax: %v", err)
	}

	if income > maxIncome {
		return income * (maxPercent / 100), nil
	}

	taxPercent, err := u.FindTaxPercent(income)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate tax: %v", err)
	}

	return income * (taxPercent / 100), nil
}

func (u *taxUsecase) SetValueToTaxLevel(taxLevels []EachTaxLevel, tax float64) ([]TaxLevelResponse, error) {
	result := make([]TaxLevelResponse, 0, len(taxLevels))

	for _, level := range taxLevels {
		levelMinIncome, levelMaxIncome := level.MinMax[0], level.MinMax[1]
		if tax >= levelMinIncome && tax <= levelMaxIncome || levelMinIncome == levelMaxIncome && tax >= levelMaxIncome {
			result = append(result, TaxLevelResponse{
				Level: level.Level,
				Tax:   tax,
			})
			continue
		}

		result = append(result, TaxLevelResponse{
			Level: level.Level,
			Tax:   level.Tax,
		})
	}

	return result, nil
}

func (u *taxUsecase) GetTaxLevelDetails(tax float64) ([]TaxLevelResponse, error) {
	taxLevel, err := u.GetTaxLevel()
	if err != nil {
		return nil, fmt.Errorf("failed to get tax level")
	}

	result, err := u.SetValueToTaxLevel(taxLevel, tax)
	if err != nil {
		return nil, fmt.Errorf("failed to get tax level")
	}

	return result, nil
}

func (u *taxUsecase) CalculateTaxWithoutWHT(req *CalculateTaxRequest) (float64, error) {
	result, err := u.DecreasePersonalAllowance(req.TotalIncome)
	if err != nil {
		return 0, err
	}
	result = math.Max(0, result)

	result = u.DecreaseAllowance(result, req.Allowances)
	result = math.Max(0, result)

	result, err = u.CalculateTaxByTaxLevel(result)
	if err != nil {
		return 0, err
	}

	return result, nil
}
