package crowdflow

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http/httptest"
)

var _ = Describe("Helper", func() {
	Describe("renderJson", func() {
		var (
			writer = httptest.NewRecorder()
		)

		BeforeEach(func() {
			renderJson(writer, "123")
		})

		It("returns the correct result", func() {
			Expect(writer.Header().Get("Content-Type")).To(Equal("application/json"))
			Expect(writer.Header().Get("Access-Control-Allow-Origin")).To(Equal("*"))
			Expect(writer.Body.String()).To(ContainSubstring("123"))
		})
	})
})
