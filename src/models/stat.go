package models

import (
	"time"
	"github.com/jinzhu/gorm/dialects/postgres"
)

type SmsStat struct {
	ID       uint    `gorm:"primary_key"`
	Created  time.Time `msg:"-"`
	Modified time.Time `msg:"-"`

	SmsCostId           uint

	ManualOverride  bool

	// computed
	SuccessRateComputed  float64
	LatencyComputed             float64

	// counters
	ApiSuccessCount uint64
	ApiTotalCount   uint64

	DeliverySuccessCount uint64
	DeliveryTotalCount   uint64

	// TODO DLR error codes aggregate
	CarrierErrorCount uint64
	ProviderErrorCount uint64

	// sliding windows
	Api      ApiSlidingWindow      `gorm:"-"`
	Delivery DeliverySlidingWindow `gorm:"-"`

	Meta   postgres.Jsonb `msg:"-"`
}

type ApiSlidingWindow struct {
	// rolling list
	ApiSuccessRatio   []float64
	AverageApiLatency []float64
}

type DeliverySlidingWindow struct {
	DeliverySuccessRatio []float64
	AverageDeliveryTime  []float64

	DlrCodes []float64
}
