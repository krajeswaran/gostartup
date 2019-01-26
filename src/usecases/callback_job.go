package usecases

import (
	"gostartup/src/adapters"
	"github.com/rs/zerolog/log"
)

var callbackUsecase = CallbackUsecase{}

type CallbackJob struct {
	MsgId string
	Payload map[string]interface{}
}

func (c CallbackJob) Name() string {
	return "Callback for Msg id: " + c.MsgId
}

func (c CallbackJob) Execute() error {
	mdr, err := callbackUsecase.MessageCallback(c.MsgId, c.Provider, c.Payload)
	if err != nil {
		log.Error().Err(err).Msgf("Message callback failed for %s", c.MsgId)
		return err
	}

	err = smsUsecase.WriteMessageDetails(mdr)
	if err != nil {
		log.Error().Err(err).Msgf("MDR update failed for %s", mdr.ID.String())
		return err
	}

	return nil
}

