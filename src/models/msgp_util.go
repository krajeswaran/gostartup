package models

import (
	"github.com/shopspring/decimal"
	"github.com/satori/go.uuid"
)

func DecimalToString(decimal decimal.Decimal) string {
	return decimal.String()
}

func StringToDecimal(str string) decimal.Decimal {
	v, err := decimal.NewFromString(str)
	if err != nil {
		return decimal.Zero
	}
	return v
}

func UuidToString(uuid uuid.UUID) string {
	return uuid.String()
}

func StringToUuid(str string) uuid.UUID {
	v, err := uuid.FromString(str)
	if err != nil {
		return uuid.Nil
	}
	return v
}
