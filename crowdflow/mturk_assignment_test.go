package crowdflow

import (
	"github.com/crowdmob/goamz/exp/mturk"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MTurkAssignment", func() {
	Describe("mturkPrice", func() {
		var (
			price mturk.Price
		)

		BeforeEach(func() {
			price = mturkPrice(1)
		})

		It("returns the correct price", func() {
			Expect(price.Amount).To(Equal("0.01"))
			Expect(price.CurrencyCode).To(Equal("USD"))
		})
	})
})
