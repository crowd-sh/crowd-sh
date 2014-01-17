package crowdflow

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Helper", func() {
	var (
		assignment Assignment
	)

	BeforeEach(func() {
		assignment = Assignment{}
	})

	Describe("Assign", func() {
		BeforeEach(func() {
			assignment.Assign()
		})

		It("has an Id", func() {
			Expect(assignment.Id).ToNot(BeEmpty())
		})

		It("is assigned", func() {
			Expect(assignment.Assigned).To(BeTrue())
		})
	})

	PDescribe("UnassigneIfExpired", func() {
		It("returns the correct result", func() {
			// Expect(writer.Header().Get("Content-Type")).To(Equal("application/json"))
			// Expect(writer.Header().Get("Access-Control-Allow-Origin")).To(Equal("*"))
			// Expect(writer.Body.String()).To(ContainSubstring("123"))
		})
	})

})
