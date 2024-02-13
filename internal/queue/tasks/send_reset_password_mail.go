package tasks

import (
	"context"
	"encoding/json"
	"keeper/internal/config"
	"keeper/internal/pkg/mailer"

	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
)

func SendResetPasswordMail(ctx context.Context, t *asynq.Task) error {
	var p UserResetPasswordMailPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		logrus.WithError(err).Error("failed to unmarshal reset password payload")
		return err
	}

	cfg := config.New()
	mailSvc := mailer.NewMailer(cfg)

	err := mailSvc.SendEmailVerificationMail(
		p.ReceiverEmailAddr,
		p.ReceiverName,
		p.Subject,
		p.TemplateData,
	)

	if err != nil {
		logrus.WithError(err).Error("failed to send reset password mail")
		return err
	}

	return nil
}
