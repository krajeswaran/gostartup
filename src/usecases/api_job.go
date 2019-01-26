package usecases

import "gostartup/src/models"

type ApiPostJob struct {
	Routes []models.SmsRoute
	Acct   *models.Account
	Msg    *models.Message
}

func (c ApiPostJob) Name() string {
	return "Api Post process for Msg id: " + c.Msg.Id
}

func (c ApiPostJob) Execute() error {
	smsUsecase.ProcessOutboundMessage(c.Routes, c.Msg, c.Acct)

	return nil
}

