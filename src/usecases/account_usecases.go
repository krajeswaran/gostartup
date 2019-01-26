package usecases

import (
	"gostartup/src/adapters"
	"gostartup/src/models"
	"github.com/satori/go.uuid"
	"github.com/rs/zerolog/log"
)

// SmsUsecase - Struct to logically bind all the use case functionality
type AccountUsecase struct{}

var dbAdapter = adapters.DBAdapter{}

// GetAccountByAuthID - get account / sub account based on auth id
func (a *AccountUsecase) AuthenticateUser(authID string, authToken string) (*models.Account, error) {
	// TODO FIXME implement
	acctId, _ := uuid.FromString("1e8531f7-f5cb-4d6f-984c-8295afa19451")
	return a.GetAccountByAcctId(acctId)
}

func (a *AccountUsecase) GetAccountByAcctId(acctId uuid.UUID) (*models.Account, error) {
	acct, err := dbAdapter.FetchPlanAndAccount(acctId)
	if err != nil {
		log.Error().Err(err).Msgf("Can't find account for id: %s", acctId.String())
		return nil, err
	}

	return acct, nil
}