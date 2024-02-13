package tasks

import (
	"encoding/json"

	"github.com/hibiken/asynq"
)

const (
	TypeUserVerificationMail  = "email:user_verification"
	TypeUserResetPasswordMail = "email:reset_password"
)

type UserVerificationMailPayload struct {
	ReceiverEmailAddr string
	ReceiverName      string
	Subject           string
	TemplateData      interface{}
}

type UserResetPasswordMailPayload struct {
	ReceiverEmailAddr string
	ReceiverName      string
	Subject           string
	TemplateData      interface{}
}

// create the tasks
func NewUserVerificationMailTask(receiverEmailAddr string, receiverName string, subject string, templateData interface{}) (*asynq.Task, error) {
	payload, err := json.Marshal(UserVerificationMailPayload{
		ReceiverEmailAddr: receiverEmailAddr,
		ReceiverName:      receiverName,
		Subject:           subject,
		TemplateData:      templateData,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeUserVerificationMail, payload), nil
}

func NewUserResetPasswordMailTask(receiverEmailAddr string, receiverName string, subject string, templateData interface{}) (*asynq.Task, error) {
	payload, err := json.Marshal(UserResetPasswordMailPayload{
		ReceiverEmailAddr: receiverEmailAddr,
		ReceiverName:      receiverName,
		Subject:           subject,
		TemplateData:      templateData,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeUserResetPasswordMail, payload), nil
}
