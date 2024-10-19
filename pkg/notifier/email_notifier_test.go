package notifier_test

import (
	"github.com/igodwin/secretsanta/pkg/notifier"
	"github.com/igodwin/secretsanta/pkg/participant"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"net/smtp"
)

var _ = Describe("Email Notifier", func() {
	var (
		emailNotifier, badEmailNotifier *notifier.EmailNotifier
		testParticipant                 = &participant.Participant{
			Name:             "Test",
			NotificationType: "stdout",
			ContactInfo:      []string{""},
			Exclusions:       []string{""},
			Recipient: &participant.Participant{
				Name:             "TestRecipient",
				NotificationType: "stdout",
				ContactInfo:      []string{""},
				Exclusions:       []string{""},
			},
		}
		mockSendMail = func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
			return nil
		}
	)

	BeforeEach(func() {
		badEmailNotifier = &notifier.EmailNotifier{
			Host:        "",
			Port:        "",
			Identity:    "",
			Username:    "",
			Password:    "",
			FromAddress: "",
			FromName:    "",
		}
		emailNotifier = &notifier.EmailNotifier{
			Host:         "smtp.example.com",
			Port:         "587",
			Identity:     "",
			Username:     "user@example.com",
			Password:     "password",
			FromAddress:  "noreply@example.com",
			FromName:     "",
			SendMailFunc: mockSendMail,
		}
	})

	Context("SendNotification", func() {
		It("should not error when smtp is configured", func() {
			Expect(emailNotifier.SendNotification(testParticipant)).NotTo(HaveOccurred())
		})

		It("should error when smtp is not configured", func() {
			err := badEmailNotifier.SendNotification(testParticipant)
			Expect(err).To(HaveOccurred())
		})
	})

	Context("IsConfigured", func() {
		It("should not error when smtp is configured", func() {
			Expect(emailNotifier.IsConfigured()).NotTo(HaveOccurred())
		})

		It("should error when smtp is not configured", func() {
			err := badEmailNotifier.IsConfigured()
			Expect(err).To(MatchError("smtp is not configured"))
		})
	})
})
