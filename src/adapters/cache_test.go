package adapters

import (
	"errors"
	"github.com/bouk/monkey"
	"github.com/go-redis/redis"
	"gostartup/src/config"
	"gostartup/src/models"
	"reflect"
	"testing"
)

func setup() {
	config.Init()
	Init()
}

func tearDown() {

	//c := Cache()
	//c.FlushAll()

	monkey.UnpatchAll()

}

func createDummyAccountByAuthID(authID string, accountID int) {
	//Cache().Set(authID, accountID, 0).Result()
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

func TestCacheAdapter_SetAccount(t *testing.T) {
	setup()

	accountID := 123
	account := map[string]interface{}{
		"id":         accountID,
		"name":       "sample name",
		"auth_id":    "MANGY5OWI0OWI0ODC2OW",
		"auth_token": "MDA3ZGEyZDJiMzFjY2NjNjYwZTNlY2JmYzkwNGFl",
		"enabled":    true,
		"country_id": "US",
	}

	err := cacheAdapter.SetAccount(accountID, models.AccountPrefix, account)

	if err != nil {
		t.Error(err.Error())
	}

	tearDown()
}

func TestCacheAdapter_SetAccount2(t *testing.T) {
	setup()

	accountID := 123
	account := map[string]interface{}{
		"id":         accountID,
		"name":       "sample name",
		"auth_id":    "MANGY5OWI0OWI0ODC2OW",
		"auth_token": "MDA3ZGEyZDJiMzFjY2NjNjYwZTNlY2JmYzkwNGFl",
		"enabled":    true,
		"country_id": "US",
	}

	var cmd *redis.StatusCmd
	monkey.PatchInstanceMethod(reflect.TypeOf(cmd), "Result", func(_ *redis.StatusCmd) (string, error) {
		return "0", errors.New("Error saving account")
	})

	err := cacheAdapter.SetAccount(accountID, models.AccountPrefix, account)
	if err == nil {
		t.Error("Set account test failed")
	}

	tearDown()
}

func TestCacheAdapter_GetAccount(t *testing.T) {
	setup()

	accountID := 123
	authID := "MANGY5OWI0OWI0ODC2OW"
	createDummyAccount(accountID, authID)

	account, err := cacheAdapter.GetAccount(accountID, models.AccountPrefix)
	if err != nil {
		t.Error("Get account test failed")
	}

	if len(account) == 0 {
		t.Error("Failed: Account should not be empty")
	}

	tearDown()
}

func TestCacheAdapter_GetAccount2(t *testing.T) {
	setup()

	accountID := 123
	authID := "MANGY5OWI0OWI0ODC2OW"
	createDummyAccount(accountID, authID)

	var cmd *redis.StringStringMapCmd
	monkey.PatchInstanceMethod(reflect.TypeOf(cmd), "Result", func(_ *redis.StringStringMapCmd) (map[string]string, error) {
		return make(map[string]string), errors.New("Error saving account")
	})

	account, err := cacheAdapter.GetAccount(accountID, models.AccountPrefix)
	if err == nil {
		t.Error("Get account test failed")
	}

	if len(account) > 0 {
		t.Error("Failed: Account should be empty due to error")
	}

	tearDown()
}

func TestCacheAdapter_GetAccount3(t *testing.T) {
	setup()

	accountID := 123
	authID := "MANGY5OWI0OWI0ODC2OW"
	createDummyAccount(accountID, authID)
	subAccountID := 321
	subAuthID := "SAFDSFDSGDFGD7FD"
	createDummySubAccount(subAccountID, subAuthID)

	subAccount, err := cacheAdapter.GetAccount(subAccountID, models.SubAccountPrefix)

	if err != nil {
		t.Error("Get sub account failed")
	}

	if len(subAccount) == 0 {
		t.Error("Failed: Sub account should not be empty")
	}

	tearDown()
}

func TestCacheAdapter_GetAccountByAuthID(t *testing.T) {
	setup()

	accountID := 123
	authID := "MANGY5OWI0OWI0ODC2OW"
	createDummyAccount(accountID, authID)

	account, err := cacheAdapter.GetAccountByAuthID(authID, models.AccountPrefix)

	if err != nil {
		t.Error("Get account by Auth ID failed")
	}

	if len(account) == 0 {
		t.Error("Failed: account should not be empty")
	}

	tearDown()
}

func TestCacheAdapter_GetAccountByAuthID2(t *testing.T) {
	setup()

	accountID := 123
	authID := "MANGY5OWI0OWI0ODC2OW"
	createDummyAccount(accountID, authID)

	var cmd *redis.StringCmd
	monkey.PatchInstanceMethod(reflect.TypeOf(cmd), "Result", func(_ *redis.StringCmd) (string, error) {
		return "0", errors.New("Error saving account")
	})

	account, err := cacheAdapter.GetAccountByAuthID(authID, models.AccountPrefix)

	if err == nil {
		t.Error("Get account by auth ID failed")
	}

	if len(account) > 0 {
		t.Error("Failed: Account should be empty due to error")
	}

	tearDown()
}

func TestCacheAdapter_GetAccountByAuthID3(t *testing.T) {
	setup()

	_, err := cacheAdapter.GetAccountByAuthID("MAFGDGSDFGD", models.AccountPrefix)

	if err == nil {
		t.Error("Account should not have found")
	}

	tearDown()
}

func TestCacheAdapter_SetAccountByAuthID(t *testing.T) {
	setup()

	accountID := 123
	authID := "MANGY5OWI0OWI0ODC2OW"
	account := map[string]interface{}{
		"id":         accountID,
		"name":       "sample name",
		"auth_id":    authID,
		"auth_token": "MDA3ZGEyZDJiMzFjY2NjNjYwZTNlY2JmYzkwNGFl",
		"enabled":    true,
		"country_id": "US",
	}

	err := cacheAdapter.SetAccountByAuthID(authID, accountID, models.AccountPrefix, account)

	if err != nil {
		t.Error("Set account by auth id failed")
	}

	tearDown()
}

func TestCacheAdapter_SetAccountByAuthID2(t *testing.T) {
	setup()

	accountID := 123
	authID := "MANGY5OWI0OWI0ODC2OW"
	account := map[string]interface{}{
		"id":         accountID,
		"name":       "sample name",
		"auth_id":    authID,
		"auth_token": "MDA3ZGEyZDJiMzFjY2NjNjYwZTNlY2JmYzkwNGFl",
		"enabled":    true,
		"country_id": "US",
	}

	var cmd *redis.StatusCmd
	monkey.PatchInstanceMethod(reflect.TypeOf(cmd), "Result", func(_ *redis.StatusCmd) (string, error) {
		return "0", errors.New("Error saving account")
	})

	err := cacheAdapter.SetAccountByAuthID(authID, accountID, models.AccountPrefix, account)

	if err == nil {
		t.Error("Set account by auth id failed")
	}

	tearDown()
}

func TestCacheAdapter_DeleteAccountByID(t *testing.T) {
	setup()

	accountID := 123
	authID := "MANGY5OWI0OWI0ODC2OW"
	createDummyAccount(accountID, authID)

	err := cacheAdapter.DeleteAccountByID("123")
	if err != nil {
		t.Error("Delete account failed")
	}

	tearDown()
}

func TestCacheAdapter_DeleteAccountByID2(t *testing.T) {
	setup()

	err := cacheAdapter.DeleteAccountByID("1234")
	if err == nil {
		t.Error("Delete account failed")
	}

	tearDown()
}

func TestCacheAdapter_DeleteAccountByID3(t *testing.T) {
	setup()

	subAccountID := 321
	subAuthID := "SAFDSFDSGDFGD7FD"
	createDummySubAccount(subAccountID, subAuthID)

	err := cacheAdapter.DeleteAccountByID("321")
	if err != nil {
		t.Error("Deleting sub account failed")
	}

	tearDown()
}
