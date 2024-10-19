package notifier_test

import (
	"bytes"
	"fmt"
	"github.com/igodwin/secretsanta/pkg/notifier"
	"github.com/igodwin/secretsanta/pkg/participant"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("Stdout Notifier", func() {
	var stdoutNotifier *notifier.Stdout
	var testParticipant *participant.Participant

	BeforeEach(func() {
		stdoutNotifier = &notifier.Stdout{}
		testParticipant = &participant.Participant{
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
	})

	Context("SendNotification", func() {
		It("should send the message to stdout", func() {
			originalStdout := os.Stdout

			r, w, _ := os.Pipe()
			os.Stdout = w

			err := stdoutNotifier.SendNotification(testParticipant)
			Expect(err).NotTo(HaveOccurred())

			Expect(w.Close()).To(Succeed())
			os.Stdout = originalStdout

			var buf bytes.Buffer
			_, _ = buf.ReadFrom(r)
			Expect(buf.String()).To(Equal(fmt.Sprintf("%s has %s\n", testParticipant.Name, testParticipant.Recipient.Name)))
		})
	})

	Context("IsConfigured", func() {
		It("does not return an error", func() {
			Expect(stdoutNotifier.IsConfigured()).To(BeNil())
		})
	})
})
