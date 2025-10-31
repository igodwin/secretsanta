package notification

import (
	"fmt"
	"os"

	"github.com/igodwin/secretsanta/pkg/config"
	"github.com/igodwin/secretsanta/pkg/notifier"
	"github.com/igodwin/secretsanta/pkg/participant"
)

func Send(participants []*participant.Participant, appConfig *config.Config) error {
	// Check for service address from config first, then environment
	notifierServiceAddr := appConfig.Notifier.ServiceAddr
	if notifierServiceAddr == "" {
		notifierServiceAddr = os.Getenv("NOTIFIER_SERVICE_ADDR")
	}

	if notifierServiceAddr != "" {
		return sendViaGRPC(participants, notifierServiceAddr, appConfig.Notifier.ArchiveEmail, appConfig.Notifier.APIKey, appConfig.SMTP.ContentType)
	}

	return sendViaLegacy(participants, appConfig)
}

func sendViaGRPC(participants []*participant.Participant, serverAddr, archiveEmail, apiKey string, contentType string) error {
	grpcNotifier, err := NewGRPCNotifierWithAPIKey(serverAddr, apiKey, &PapaElfTemplate{})
	if err != nil {
		return fmt.Errorf("failed to create gRPC notifier: %w", err)
	}
	defer grpcNotifier.Close()

	return grpcNotifier.SendBatchNotifications(participants, archiveEmail, contentType)
}

func sendViaLegacy(participants []*participant.Participant, appConfig *config.Config) error {
	var notifierInstance notifier.Notifier
	var emailNotifier = &notifier.EmailNotifier{
		Host:        appConfig.SMTP.Host,
		Port:        appConfig.SMTP.Port,
		Identity:    appConfig.SMTP.Identity,
		Username:    appConfig.SMTP.Username,
		Password:    appConfig.SMTP.Password,
		FromAddress: appConfig.SMTP.FromAddress,
		FromName:    appConfig.SMTP.FromName,
		ContentType: appConfig.SMTP.ContentType,
	}

	for _, participant := range participants {
		switch participant.NotificationType {
		case "email":
			notifierInstance = emailNotifier
		default:
			notifierInstance = &notifier.Stdout{}
		}

		err := notifierInstance.IsConfigured()
		if err != nil {
			return err
		}
		err = notifierInstance.SendNotification(participant)
		if err != nil {
			return err
		}
	}
	return nil
}