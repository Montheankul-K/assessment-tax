package tax

import "gorm.io/gorm"

type TaxAllowance struct {
	gorm.Model
	AllowanceType      string  `gorm:"not null" json:"allowance_type"`
	MinAllowanceAmount float64 `gorm:"not null" json:"min_allowance_amount"`
	MaxAllowanceAmount float64 `gorm:"not null" json:"max_allowance_amount"`
}

type TaxLevel struct {
	gorm.Model
	TaxLevel   uint    `gorm:"not null" json:"tax_level"`
	MinIncome  float64 `gorm:"type:decimal(10,2) not null" json:"min_income"`
	MaxIncome  float64 `gorm:"type:decimal(10,2) not null" json:"max_income"`
	TaxPercent float64 `gorm:"type:decimal(10,2) not null" json:"tax_percent"`
}

type AllowanceFilter struct {
	AllowanceType string
}

type TaxLevelFilter struct {
	Income float64
}

func (TaxAllowance) TableName() string {
	return "tax_allowance"
}

func (TaxLevel) TableName() string {
	return "tax_level"
}
