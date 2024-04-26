package taxHandlers

type TaxAllowanceDetails struct {
	AllowanceType string
	Amount        float64
}

type CalculateTaxRequest struct {
	TotalIncome float64
	Wht         float64
	Allowances  []TaxAllowanceDetails
}

func NewCalculateTaxRequest() *CalculateTaxRequest {
	return &CalculateTaxRequest{}
}
