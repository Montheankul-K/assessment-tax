package taxRepositories

import (
	"errors"
	"fmt"
	"github.com/montheankul-k/assessment-tax/modules/tax"
	"gorm.io/gorm"
)

type ITaxRepository interface {
	FindBaselineAllowanceAmount(req *tax.AllowanceFilter) (float64, float64, error)
	FindTaxPercentByIncome(req *tax.TaxLevelFilter) (float64, error)
	FindMaxIncomeAndPercent() (float64, float64, error)
}

type taxRepository struct {
	db *gorm.DB
}

func TaxRepository(db *gorm.DB) ITaxRepository {
	return &taxRepository{
		db: db,
	}
}

func (t *taxRepository) FindBaselineAllowanceAmount(req *tax.AllowanceFilter) (float64, float64, error) {
	var taxAllowance tax.TaxAllowance
	if result := t.db.Select("min_allowance_amount", "max_allowance_amount").Where(req.AllowanceType).First(&taxAllowance); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0, 0, fmt.Errorf("baseline amount for %s not found", req.AllowanceType)
		}

		return 0, 0, fmt.Errorf("can't find baseline amount for %s", req.AllowanceType)
	}

	return taxAllowance.MinAllowanceAmount, taxAllowance.MaxAllowanceAmount, nil
}

func (t *taxRepository) FindTaxPercentByIncome(req *tax.TaxLevelFilter) (float64, error) {
	var taxLevel tax.TaxLevel
	if result := t.db.Select("tax_percent").Where("min_income <= ? AND max_income >= ?", req.Income, req.Income).First(&taxLevel); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0, fmt.Errorf("income for %s not found", req.Income)
		}

		return 0, fmt.Errorf("can't find income for %s", req.Income)
	}

	return taxLevel.TaxPercent, nil
}

func (t *taxRepository) FindMaxIncomeAndPercent() (float64, float64, error) {
	var taxLevel tax.TaxLevel
	if result := t.db.Order("max_income DESC").Select("max_income, tax_percent").First(&taxLevel); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return 0, 0, fmt.Errorf("max income and max tax percent not found")
		}

		return 0, 0, fmt.Errorf("can't find max income and max tax percent")
	}

	return taxLevel.MaxIncome, taxLevel.TaxPercent, nil
}