package taxHandlers

import (
	"github.com/montheankul-k/assessment-tax/config"
	"github.com/montheankul-k/assessment-tax/modules/tax"
	"github.com/montheankul-k/assessment-tax/modules/tax/taxUsecases"
)

type taxHandlerErrCode string

type ITaxHandler interface {
	FindBaselineAllowance(allowanceType string) (float64, float64, error)
	FindTaxPercent(totalIncome float64) (float64, error)
}

type taxHandler struct {
	config     config.IConfig
	taxUsecase taxUsecases.ITaxUsecase
}

func TaxHandler(config config.IConfig, taxUsecase taxUsecases.ITaxUsecase) ITaxHandler {
	return &taxHandler{
		config:     config,
		taxUsecase: taxUsecase,
	}
}

func (h *taxHandler) FindBaselineAllowance(allowanceType string) (float64, float64, error) {
	req := tax.AllowanceFilter{
		AllowanceType: allowanceType,
	}

	minAllowanceAmount, maxAllowanceAmount, err := h.taxUsecase.FindBaselineAllowance(&req)
	if err != nil {
		return 0, 0, err
	}

	return minAllowanceAmount, maxAllowanceAmount, nil
}

func (h *taxHandler) FindTaxPercent(totalIncome float64) (float64, error) {
	req := tax.TaxLevelFilter{
		Income: totalIncome,
	}

	taxPercent, err := h.taxUsecase.FindTaxPercent(&req)
	if err != nil {
		return 0, nil
	}

	return taxPercent, nil
}
