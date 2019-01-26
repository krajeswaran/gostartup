package adapters

import (
	"errors"
	"gostartup/src/models"
	"github.com/satori/go.uuid"
	"github.com/jinzhu/gorm"
)

// DBAdapter - Struct to logically bind all the database related functions
type DBAdapter struct{}

// DeepStatus checks health of database
func (d *DBAdapter) DeepStatus() error {
	if db.HasTable(&models.Account{}) == false {
		return errors.New("SERVICE_DB_DOWN")
	}
	return nil
}

func (d *DBAdapter) GetRoutesByMccMnc(mcc string, mnc string) ([]models.SmsRoute, error) {
	// get all possible routes for mcc/mnc
	var routeCosts []models.SmsCost
	if err := db.Where(models.SmsCost{Mcc: mcc, Mnc: mnc}).Find(&routeCosts).Error; err != nil {
		return nil, err
	}

	if len(routeCosts) == 0 {
		return nil, errors.New("no costs found")
	}

	var routes []models.SmsRoute
	// in case of any cache miss, abort the whole thing and let db fill up
	for _, routeCost := range routeCosts {
		route := models.SmsRoute{}

		route.SmsCost = routeCost

		// ignore empty stats for routes
		var stat models.SmsStat
		db.Where(models.SmsStat{SmsCostId: routeCost.ID}).First(&stat)
		route.SmsStat = stat

		// fetch each provider
		var provider models.Provider
		if err := db.Where(models.Provider {
			ID: routeCost.ProviderId,
			Status:models.PROVIDER_STATUS_ENABLED}).First(&provider).Error;
			err != nil {
			return nil, err
		}
		route.Provider = provider

		// append super struct back
		routes = append(routes, route)
	}

	return routes, nil
}

func (d *DBAdapter) WriteProvider(provider *models.Provider) error {
	provider.Modified = gorm.NowFunc()
	if provider.ID <= 0 {
		provider.Created = gorm.NowFunc()
	}

	if err := db.Save(&provider).Error; err != nil {
		return err
	}

	return nil
}

func (d *DBAdapter) GetProvider(providerId uint, shouldLoadAll bool) (*models.Provider, error) {
	var provider models.Provider
	if shouldLoadAll {
		if err := db.Where(models.Provider{ID:providerId}).Find(&provider).Error; err != nil {
			return nil, err
		}
	} else {
		if err := db.Where(models.Provider{ID:providerId, Status:models.PROVIDER_STATUS_ENABLED}).Find(&provider).Error; err != nil {
			return nil, err
		}
	}

	return &provider, nil
}

func (d *DBAdapter) GetAllProviders() ([]models.Provider, error) {
	var providers []models.Provider
	if err := db.Find(&providers).Error; err != nil {
		return nil, err
	}

	return providers, nil
}

func (d *DBAdapter) GetMDR(msgId string) (*models.MessageDetailRecord, error) {
	var mdr models.MessageDetailRecord
	msgUuid, _  := uuid.FromString(msgId)
	if err := db.Where(models.MessageDetailRecord{ID: msgUuid}).First(&mdr).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// return empty response all the way back
			return nil, nil
		}
		return nil, err
	}

	return &mdr, nil
}

func (d *DBAdapter) UpdateMDR(mdr *models.MessageDetailRecord) error {
	mdr.Modified = gorm.NowFunc()
	if err := db.Update(&mdr).Error; err != nil {
		return err
	}

	return nil
}

func (d *DBAdapter) InsertMDR(mdr *models.MessageDetailRecord) error {
	mdr.Modified = gorm.NowFunc()
	mdr.Created = mdr.QueuedTime
	if err := db.Create(&mdr).Error; err != nil {
		return err
	}

	return nil
}

func (d *DBAdapter) GetBalance(acctId uuid.UUID) (*models.Balance, error) {
	var balance models.Balance
	if err := db.Where(models.Balance{AccountId: acctId}).First(&balance).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// return empty response all the way back
			return nil, nil
		}
		return nil, err
	}

	return &balance, nil
}

func (d *DBAdapter) UpdateBalances(balances []models.Balance) error {
	for i := range balances {
		bal := balances[i]
		bal.Modified = gorm.NowFunc()
		if err := db.Save(&bal).Error; err != nil {
			return err
		}
	}

	return nil
}

func (d *DBAdapter) UpdateStats(stats []models.SmsStat) error {
	for i := range stats {
		stat := stats[i]
		stat.Modified = gorm.NowFunc()
		if err := db.Save(&stat).Error; err != nil {
			return err
		}
	}

	return nil
}

func (d *DBAdapter) FetchPlanAndAccount(acctId uuid.UUID) (*models.Account, error) {
	var acctPlan models.Account
	if err := db.Where(models.Account{ID: acctId, Status:models.ACCT_STATUS_ENABLED}).First(&acctPlan).Error; err != nil {
		return nil, err
	}
	var plan models.Plan
	if err := db.Where(models.Plan{ID:acctPlan.PlanId}).First(&plan).Error; err != nil {
		return nil, err
	}
	acctPlan.Plan = plan

	return &acctPlan, nil
}

func (d *DBAdapter) LoadRateForRoute(route *models.SmsRoute) (*models.SmsRate, error) {
	var rate models.SmsRate
	if err := db.Where("sms_cost_id = ?", route.SmsCost.ID).First(&rate).Error; err != nil {
		return nil, err
	}

	return &rate, nil
}
