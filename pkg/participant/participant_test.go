package participant_test

import (
	"fmt"
	. "github.com/igodwin/secretsanta/pkg/participant"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Participant", func() {
	var ind0, ind1, ind2 *Participant
	BeforeEach(func() {
		ind0 = &Participant{
			Name:             "John Doe",
			NotificationType: "email",
			ContactInfo:      []string{"johndoe@example.com"},
			Exclusions:       []string{""},
		}
		ind1 = &Participant{
			Name:             "Jane Doe",
			NotificationType: "email",
			ContactInfo:      []string{"janedoe@example.com"},
			Exclusions:       []string{"John Doe"},
		}
		ind2 = &Participant{
			Name:             "Jill Doe",
			NotificationType: "email",
			ContactInfo:      []string{"jilldoe@example.com"},
			Exclusions:       []string{""},
		}
	})

	Describe("UpdateRecipient", func() {
		It("should allow any non-excluded individual to be set", func() {
			Expect(ind0.UpdateRecipient(ind1)).To(Succeed())
			Expect(ind0.UpdateRecipient(ind2)).To(Succeed())
		})

		It("should not allow an excluded individual to be set", func() {
			Expect(ind1.UpdateRecipient(ind0)).To(MatchError(fmt.Errorf("participant %s is excluded", ind0.Name)))
		})

		It("should not allow self to be set", func() {
			Expect(ind1.UpdateRecipient(ind1)).To(MatchError("cannot update match with self"))
		})
	})
})
