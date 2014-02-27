package crowdflow

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	Describe("Default Config", func() {
		It("has the default values", func() {
			Expect(config.WithMTurk).To(Equal(true))
		})
	})

	Describe("SetConfig", func() {
		var (
			c Config
		)

		BeforeEach(func() {
			c = Config{
				WithMTurk: false,
			}
		})

		It("has correct config", func() {
			Expect(config.WithMTurk).To(Equal(true))
		})
	})
})
