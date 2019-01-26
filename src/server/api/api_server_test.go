package api

import (
	"github.com/franela/goreq"
	"gostartup/src/adapters"
	"gostartup/src/config"
	"gostartup/src/models"
	"net/http"
	"testing"
)

const AccountServiceURL = "http://127.0.0.1:5000/gostartup/v1/account/"

var authPassword string
var authUserName string

var cacheAdapter *adapters.CacheAdapter

func setup() {
	config.Init()
	adapters.Init()

	c := config.GetConfig()
	authUserName = c.GetString("auth.auth_id")
	authPassword = c.GetString("auth.auth_token")

	go ApiInit()
}

func tearDown() {
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
	cacheAdapter.SetAccountByAuthID(authID, accountID, models.AccountPrefix, account)
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
	cacheAdapter.SetAccountByAuthID(authID, subAccountID, models.SubAccountPrefix, subAccount)
}

func request(authID string, phone string) (*goreq.Response, error) {
	if authID != "" {
		url := AccountServiceURL + authID + "/"
		res, err := goreq.Request{
			Method:            "GET",
			Uri:               url,
			BasicAuthPassword: authPassword,
			BasicAuthUsername: authUserName,
		}.Do()
		return res, err
	}
	url := AccountServiceURL + "?phone=" + phone
	res, err := goreq.Request{
		Method:            "GET",
		Uri:               url,
		BasicAuthPassword: authPassword,
		BasicAuthUsername: authUserName,
	}.Do()
	return res, err
}

func TestController_GetAccountByAuthID(t *testing.T) {
	setup()

	authID := "MANGY5OWI0OWI0ODC2OW"
	accountID := 123
	createDummyAccount(accountID, authID)

	res, err := request(authID, "")
	if err != nil {
		t.Error(err.Error())
	}

	if res.StatusCode != http.StatusOK {
		t.Error("Response status is not ok")
	}

	tearDown()
}

func TestController_GetAccountByAuthID2(t *testing.T) {
	setup()

	authID := "MANGY5OWI0OWI0ODC2OW"
	accountID := 123
	createDummyAccount(accountID, authID)

	res, err := request("MAFDSFGDSGSDG", "")
	if err != nil {
		t.Error(err.Error())
	}

	if res.StatusCode == http.StatusOK {
		t.Error("Failed: Response status is ok")
	}
	tearDown()
}

func TestController_DeleteAccount(t *testing.T) {
	setup()

	authID := "MANGY5OWI0OWI0ODC2OW"
	accountID := 123
	createDummyAccount(accountID, authID)

	res, err := goreq.Request{
		Method:            "DELETE",
		Uri:               AccountServiceURL + "123" + "/",
		BasicAuthPassword: authPassword,
		BasicAuthUsername: authUserName,
	}.Do()
	if err != nil {
		t.Error(err.Error())
	}

	if res.StatusCode != http.StatusOK {
		t.Error("Response status is not ok")
	}

	tearDown()
}

func TestController_DeleteAccount2(t *testing.T) {
	setup()

	authID := "MANGY5OWI0OWI0ODC2OW"
	accountID := 123
	createDummyAccount(accountID, authID)

	res, err := goreq.Request{
		Method:            "DELETE",
		Uri:               AccountServiceURL + "321" + "/",
		BasicAuthPassword: authPassword,
		BasicAuthUsername: authUserName,
	}.Do()
	if err != nil {
		t.Error(err.Error())
	}

	if res.StatusCode == http.StatusOK {
		t.Error("Response status is ok")
	}

	tearDown()
}
