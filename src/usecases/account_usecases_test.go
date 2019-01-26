package usecases

import (
	"gostartup/src/adapters"
	"gostartup/src/config"
	"gostartup/src/models"
	"strconv"
	"testing"
)

var accountUsecase *SmsUsecase

func createDummyAccountByAuthID(authID string, accountID int) {
	cache := adapters.Cache()
	cache.Set(authID, accountID, 0).Result()
}

func createDummyAccount(accountID int, authID string) {
	account := map[string]interface{}{
		"id":         accountID,
		"name":       "sample name",
		"auth_id":    authID,
		"auth_token": "MDA3ZGEyZDJiMzFjY2NjNjYwZTNlY2JmYzkwNGFl",
		"enabled":    true,
		"country_id": "US",
	}
	createDummyAccountByAuthID(authID, accountID)
	cacheAdapter.SetAccount(accountID, models.AccountPrefix, account)
}

func createDummySubAccount(subAccountID int, authID string) {
	subAccount := map[string]interface{}{
		"id":          subAccountID,
		"account_id":  123,
		"name":        "sample sub account",
		"auth_id":     authID,
		"auth_token":  "FDSFDSGFNJSDKFHJSDHFGJSDF",
		"mps_allowed": 5,
		"enabled":     true,
	}
	createDummyAccountByAuthID(authID, subAccountID)
	cacheAdapter.SetAccount(subAccountID, models.SubAccountPrefix, subAccount)
}

func setup() {
	config.Init()
	adapters.Init()
}

func tearDown() {
	c := adapters.Cache()
	c.FlushDB()
}

func TestAccountUsecase_GetAccountByAuthID(t *testing.T) {
	setup()

	accountID := 123
	authID := "MANGY5OWI0OWI0ODC2OW"
	createDummyAccount(accountID, authID)

	account, err := accountUsecase.GetAccountByAuthID(authID)
	if err != nil {
		t.Error(err.Error())
	}

	if len(account) == 0 {
		t.Error("Failed: account object should not be empty")
	}

	tearDown()
}

func TestAccountUsecase_GetAccountByAuthID2(t *testing.T) {
	setup()

	accountID := 123
	authID := "MANGY5OWI0OWI0ODC2OW"
	createDummyAccount(accountID, authID)

	account, err := accountUsecase.GetAccountByAuthID("MAFDSDSFSDDFFD")
	if err == nil {
		t.Error("Get account by auth id failed")
	}

	if len(account) > 0 {
		t.Error("Failed: account object should be empty")
	}

	tearDown()
}

func TestAccountUsecase_GetAccountByAuthID3(t *testing.T) {
	setup()

	account, err := accountUsecase.GetAccountByAuthID("MAMTVINTHMMDEXZGFMND")
	if err == nil {
		t.Error(err.Error())
	}

	if len(account) != 0 {
		t.Error("Failed: account object should not be empty")
	}

	account, err = accountUsecase.GetAccountByAuthID("SARKNDGZYTQ5YTLKMTGW")
	if err == nil {
		t.Error(err.Error())
	}

	if len(account) != 0 {
		t.Error("Failed: account object should not be empty")
	}

	tearDown()
}

func TestAccountUsecase_GetMainAccountByAuthID(t *testing.T) {
	setup()

	accountID := 123
	authID := "MANGY5OWI0OWI0ODC2OW"
	createDummyAccount(accountID, authID)

	account, err := accountUsecase.GetMainAccountByAuthID(authID)
	if err != nil {
		t.Error(err.Error())
	}

	if len(account) == 0 {
		t.Error("Failed: account object should not be empty")
	}

	tearDown()
}

func TestAccountUsecase_GetMainAccountByAuthID2(t *testing.T) {
	setup()

	accountID := 123
	authID := "MANGY5OWI0OWI0ODC2OW"
	createDummyAccount(accountID, authID)

	account, err := accountUsecase.GetMainAccountByAuthID("MAFDSDSFSDDFFD")
	if err == nil {
		t.Error("Get account by auth id failed")
	}

	if len(account) > 0 {
		t.Error("Failed: account object should be empty")
	}

	tearDown()
}

func TestAccountUsecase_GetSubAccountByAuthID(t *testing.T) {
	setup()

	accountID := 123
	authID := "MANGY5OWI0OWI0ODC2OW"
	createDummyAccount(accountID, authID)
	subAccountID := 321
	subAuthID := "SAFDSFDSGDFGD7FD"
	createDummySubAccount(subAccountID, subAuthID)

	account, err := accountUsecase.GetSubAccountByAuthID(subAuthID)
	if err != nil {
		t.Error(err.Error())
	}

	if len(account) == 0 {
		t.Error("Failed: account object should not be empty")
	}

	tearDown()
}

func TestAccountUsecase_GetSubAccountByAuthID2(t *testing.T) {
	setup()

	accountID := 123
	authID := "MANGY5OWI0OWI0ODC2OW"
	createDummyAccount(accountID, authID)
	subAccountID := 321
	subAuthID := "SAFDSFDSGDFGD7FD"
	createDummySubAccount(subAccountID, subAuthID)

	account, err := accountUsecase.GetSubAccountByAuthID("SAFDSDSFSDDFFDFD")
	if err == nil {
		t.Error("Get account by auth id failed")
	}

	if len(account) > 0 {
		t.Error("Failed: account object should be empty")
	}

	tearDown()
}

func TestAccountUsecase_GetMainAccountByID(t *testing.T) {
	setup()

	accountID := 123
	authID := "MANGY5OWI0OWI0ODC2OW"
	createDummyAccount(accountID, authID)

	account, err := accountUsecase.GetMainAccountByID(accountID)
	if err != nil {
		t.Error(err.Error())
	}

	if len(account) == 0 {
		t.Error("Failed: account object should not be empty")
	}

	tearDown()
}

func TestAccountUsecase_GetMainAccountByID2(t *testing.T) {
	setup()

	accountID := 123
	authID := "MANGY5OWI0OWI0ODC2OW"
	createDummyAccount(accountID, authID)

	account, err := accountUsecase.GetMainAccountByID(12345678)
	if err == nil {
		t.Error("Get account by auth id failed")
	}

	if len(account) > 0 {
		t.Error("Failed: account object should be empty")
	}

	tearDown()
}

func TestAccountUsecase_GetSubAccountByID(t *testing.T) {
	setup()

	accountID := 123
	authID := "MANGY5OWI0OWI0ODC2OW"
	createDummyAccount(accountID, authID)
	subAccountID := 321
	subAuthID := "SAFDSFDSGDFGD7FD"
	createDummySubAccount(subAccountID, subAuthID)

	account, err := accountUsecase.GetSubAccountByID(subAccountID)
	if err != nil {
		t.Error(err.Error())
	}

	if len(account) == 0 {
		t.Error("Failed: account object should not be empty")
	}

	tearDown()
}

func TestAccountUsecase_GetSubAccountByID2(t *testing.T) {
	setup()

	accountID := 123
	authID := "MANGY5OWI0OWI0ODC2OW"
	createDummyAccount(accountID, authID)
	subAccountID := 321
	subAuthID := "SAFDSFDSGDFGD7FD"
	createDummySubAccount(subAccountID, subAuthID)

	account, err := accountUsecase.GetSubAccountByID(9874321)
	if err == nil {
		t.Error("Get account by auth id failed")
	}

	if len(account) > 0 {
		t.Error("Failed: account object should be empty")
	}

	tearDown()
}

func TestAccountUsecase_DeleteAccountByID(t *testing.T) {
	setup()

	accountID := 123
	authID := "MANGY5OWI0OWI0ODC2OW"
	createDummyAccount(accountID, authID)

	err := accountUsecase.DeleteAccountByID(strconv.Itoa(accountID))
	if err != nil {
		t.Error(err.Error())
	}

	tearDown()
}

func TestAccountUsecase_DeleteAccountByID2(t *testing.T) {
	setup()

	accountID := 123
	authID := "MANGY5OWI0OWI0ODC2OW"
	createDummyAccount(accountID, authID)

	err := accountUsecase.DeleteAccountByID("12345678")
	if err == nil {
		t.Error("Get account by auth id failed")
	}

	tearDown()
}
