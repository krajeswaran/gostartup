package usecases

import (
	"github.com/dghubble/sling"
	"github.com/rs/zerolog/log"
	"gostartup/src/adapters"
	"gostartup/src/config"
	"gostartup/src/models"
	"time"
)

type CallbackUsecase struct{}

var acctUsecase = AccountUsecase{}
var commonUsecase = CommonUsecase{}

const HEADER_SIGNATURE = "x-strofo-signature"
const HEADER_SIGNATURE_NONCE = "x-strofo-signature-nonce"
const UA_STROFO_API = "Strofo API 1.0"

func (a *CallbackUsecase) MessageCallback(msgId string, provider adapters.ProviderAdapter,
	payload map[string]interface{}) (*models.MessageDetailRecord, error) {

	c := config.GetConfig()
	// transfer provider callback to our callback format
	p, err := provider.ParseMessageCallback(payload)
	if err != nil {
		log.Error().Err(err).Msgf("Error parsing callback payload for mdr: %s", msgId)
		return nil, err
	}

	// update mdr with current status
	mdr, err := smsUsecase.MergeMessageDetails(&p.Response, msgId)
	if err != nil {
		log.Error().Err(err).Msgf("Error updating MDR with callback payload for mdr: %s", msgId)
		return nil, err
	}

	if (mdr.ExternalStatus == models.MSG_REJECTED || mdr.ExternalStatus == models.MSG_FAILED) && mdr.Guaranteed {
		if mdr.Text != "" {
			msg := models.Message{}
			msg.Text = mdr.Text
			msg.From = mdr.FromNumber
			msg.To = mdr.ToNumber
			msg.IsGuarantee = mdr.Guaranteed
			msg.Id = msgId

			acct, err := acctUsecase.GetAccountByAcctId(mdr.AccountId)
			if err != nil {
				log.Error().Msgf("CALLBACK_FAIL:: can't deliver guaranteed msg for %s, "+
					"as we seem to have lost account number: %d", msgId, mdr.AccountId)
			} else {
				// retry this msg again
				_, err = smsUsecase.SendMessage(&msg, acct)
				if err != nil {
					log.Error().Msgf("CALLBACK_FAIL:: can't deliver guaranteed msg for %s, "+
						"as send message has failed. abandon all hopes!!!",
						msgId)

				}
			}

		} else {
			log.Error().Msgf("CALLBACK_FAIL:: can't deliver guaranteed msg for %s, "+
				"as we seem to have lost message text", msgId)
		}
	}

	// callback!
	if mdr.OriginalCallback == "" {
		log.Info().Msgf("Original callback empty for: %s", msgId)
		return mdr, nil
	}

	req := sling.New().Post{
		Uri:         mdr.OriginalCallback,
		ContentType: "application/json",
		Accept:      "application/json",
		Timeout:     time.Millisecond * time.Duration(c.GetInt("callback.timeout")),
		Body:        mdr,
		UserAgent:   UA_STROFO_API,
	}

	acct, err := acctUsecase.GetAccountByAcctId(mdr.AccountId)
	if err != nil {
		log.Error().Err(err).Msgf("Error getting account Id %d for mdr: %s", mdr.AccountId, msgId)
		return nil, err
	}
	commonUsecase.CreateApiSignature(&req, acct)

	retry := c.GetInt("callback.retry")
	for i := 0; i < retry; i++ {
		res, err := req.Do()
		if err != nil || res.StatusCode >= 300 {
			log.Error().Err(err).Msgf("Error sending callback to url: %s, retry %d", mdr.OriginalCallback, i+1)
			if res != nil {
				resStr, _ := res.Body.ToString()
				log.Debug().Err(err).Msgf("response = %s", resStr)
			}
		} else if res.StatusCode < 300 {
			// we are good
			log.Debug().Err(err).Msgf("callback successful, response code = %d", res.StatusCode)
			break
		}
	}

	return mdr, nil
}
