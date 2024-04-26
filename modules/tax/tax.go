package tax

import "gorm.io/gorm"

type TaxAllowance struct {
	gorm.Model
	AllowanceType      string  `gorm:"not null"`
	MinAllowanceAmount float64 `gorm:"not null"`
	MaxAllowanceAmount float64 `gorm:"not null"`
}

type TaxLevel struct {
	gorm.Model
	MinIncome  float64 `gorm:"type:decimal(10,2) not null"`
	MaxIncome  float64 `gorm:"type:decimal(10,2) not null"`
	TaxPercent float64 `gorm:"type:decimal(10,2) not null"`
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