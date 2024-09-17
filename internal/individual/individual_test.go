package individual_test

import (
	. "github.com/igodwin/secretsanta/internal/individual"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Individual", func() {
	var ind0, ind1, ind2 *Individual
	BeforeEach(func() {
		ind0 = &Individual{
			Name:       "John Doe",
			Email:      []string{"johndoe@example.com"},
			Exclusions: []string{""},
			Match:      nil,
		}
		ind1 = &Individual{
			Name:       "Jane Doe",
			Email:      []string{"janedoe@example.com"},
			Exclusions: []string{"John Doe"},
			Match:      nil,
		}
		ind2 = &Individual{
			Name:       "Jill Doe",
			Email:      []string{"jilldoe@example.com"},
			Exclusions: []string{""},
			Match:      nil,
		}
	})

	Describe("SetMatch", func() {
		Context("without exclusions", func() {
			It("should allow any individual to be set", func() {
				Expect(ind0.SetMatch(ind1)).To(Succeed())
				Expect(ind0.SetMatch(ind2)).To(Succeed())
			})
		})
		Context("with exclusions", func() {
			It("should allow an individual that is not excluded to be set", func() {
				Expect(ind1.SetMatch(ind2)).To(Succeed())
			})

			It("should not allow an individual that is excluded to be set", func() {
				Expect(ind1.SetMatch(ind0)).NotTo(Succeed())
			})
		})
	})
})
