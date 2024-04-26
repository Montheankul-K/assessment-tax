package taxUsecases

import (
	"github.com/montheankul-k/assessment-tax/modules/tax"
	"github.com/montheankul-k/assessment-tax/modules/tax/taxRepositories"
)

type ITaxUsecase interface {
	FindBaselineAllowance(req *tax.AllowanceFilter) (float64, float64, error)
	FindTaxPercent(req *tax.TaxLevelFilter) (float64, error)
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
