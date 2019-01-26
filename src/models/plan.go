package models

import (
	"github.com/shopspring/decimal"
	"time"
)

type Plan struct {
	ID       uint `gorm:"primary_key"`
	Created  time.Time
	Modified time.Time

	Name                 string
	SetupFee             decimal.Decimal
	SubscriptionFee      decimal.Decimal
	CloudUnitRate        decimal.Decimal
	SmsFlatRate          decimal.Decimal
	NumberFlatRate       decimal.Decimal
	VoiceFlatRate        decimal.Decimal
	NumberSetupRate      decimal.Decimal
	ShouldApplyFlatrates bool
}
