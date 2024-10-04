package participant_test

import (
	"fmt"
	. "github.com/igodwin/secretsanta/internal/participant"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Participant", func() {
	var ind0, ind1, ind2 *Participant
	BeforeEach(func() {
		ind0 = &Participant{
			Name:       "John Doe",
			Email:      []string{"johndoe@example.com"},
			Exclusions: []string{""},
		}
		ind1 = &Participant{
			Name:       "Jane Doe",
			Email:      []string{"janedoe@example.com"},
			Exclusions: []string{"John Doe"},
		}
		ind2 = &Participant{
			Name:       "Jill Doe",
			Email:      []string{"jilldoe@example.com"},
			Exclusions: []string{""},
		}
	})

	Describe("UpdateRecipient", func() {
		Context("without exclusions", func() {
			It("should allow any individual to be set", func() {
				Expect(ind0.UpdateRecipient(ind1)).To(Succeed())
				Expect(ind0.UpdateRecipient(ind2)).To(Succeed())
			})
		})
		Context("with exclusions", func() {
			It("should allow an individual that is not excluded to be set", func() {
				Expect(ind1.UpdateRecipient(ind2)).To(Succeed())
			})

			It("should not allow an individual that is excluded to be set", func() {
				Expect(ind1.UpdateRecipient(ind0)).To(MatchError(fmt.Errorf("participant %s is excluded", ind0.Name)))
			})

			It("should not allow self to be set", func() {
				Expect(ind1.UpdateRecipient(ind1)).To(MatchError("cannot update match with self"))
			})
		})
	})
})
