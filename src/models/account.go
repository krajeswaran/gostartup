package models

import (
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/satori/go.uuid"
	"time"
)

const (
	// type
	ACCT_TYPE_FREE = "FREE"
	ACCT_TYPE_STD = "STANDARD"
	ACCT_TYPE_BIZ = "BUSINESS"

	// status
	ACCT_STATUS_ENABLED = "ENABLED"
	ACCT_STATUS_DISABLED = "DISABLED"
	ACCT_STATUS_BLOCKED = "BLOCKED"
)

// Account model
type Account struct {
	ID       uuid.UUID `json:"message_id" sql:"type:uuid;default:uuid_generate_v4()"`
	Created  time.Time
	Modified time.Time

	Name string
	//UserID                 int
	AccountType string
	Address     postgres.Jsonb
	Timezone    string
	Country     string
	Plan      Plan
	PlanId      uint
	Status      string
}
