package taxUsecases

import (
	"fmt"
	"github.com/montheankul-k/assessment-tax/modules/tax"
	"github.com/montheankul-k/assessment-tax/modules/tax/taxRepositories"
)

type ITaxUsecase interface {
	FindBaselineAllowance(req *tax.AllowanceFilter) (float64, float64, error)
	FindTaxPercent(req *tax.TaxLevelFilter) (float64, error)
	FindMaxIncomeAndPercent() (float64, float64, error)
	GetTaxLevel() ([]EachTaxLevel, error)
}

type taxUsecase struct {
	taxRepository taxRepositories.ITaxRepository
}

func TaxUsecase(taxRepository taxRepositories.ITaxRepository) ITaxUsecase {
	return &taxUsecase{
		taxRepository: taxRepository,
	}
}

func (u *taxUsecase) FindBaselineAllowance(req *tax.AllowanceFilter) (float64, float64, error) {
	minAllowanceAmount, maxAllowanceAmount, err := u.taxRepository.FindBaselineAllowanceAmount(req)
	if err != nil {
		return 0, 0, err
	}

	return minAllowanceAmount, maxAllowanceAmount, nil
}

func (u *taxUsecase) FindTaxPercent(req *tax.TaxLevelFilter) (float64, error) {
	taxPercent, err := u.taxRepository.FindTaxPercentByIncome(req)
	if err != nil {
		return 0, err
	}

	return taxPercent, nil
}

func (u *taxUsecase) FindMaxIncomeAndPercent() (float64, float64, error) {
	maxIncome, taxPercent, err := u.taxRepository.FindMaxIncomeAndPercent()
	if err != nil {
		return 0, 0, err
	}

	return maxIncome, taxPercent, nil
}

type EachTaxLevel struct {
	MinMax []float64
	Level  string
	Tax    float64
}

func (u *taxUsecase) GetTaxLevel() ([]EachTaxLevel, error) {
	maxIncomeAmount, _, err := u.FindMaxIncomeAndPercent()
	if err != nil {
		return nil, err
	}

	taxLevels, err := u.taxRepository.GetTaxLevel()
	if err != nil {
		return nil, err
	}

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

	return newTaxLevel, nil
}
