package crowdflow

// import (
// 	. "github.com/onsi/ginkgo"
// 	. "github.com/onsi/gomega"
// )

// var _ = Describe("SharedAssignment", func() {
// 	var (
// 		assignment Assignment
// 	)

// 	BeforeEach(func() {
// 		assignment = Assignment{}
// 	})

// 	Describe("generateId", func() {
// 		It("sets the correct id", func() {
// 			Expect(assignment.generateId()).ToNot(Equal(""))
// 		})
// 	})

// 	Describe("Assign", func() {
// 		BeforeEach(func() {
// 			assignment.Assign()
// 		})

// 		It("has the correct attributes", func() {
// 			Expect(assignment.Assigned).To(BeTrue())
// 			Expect(assignment.Id).ToNot(BeEmpty())
// 			Expect(assignment.StartedAt).ToNot(BeNil())
// 		})
// 	})

// 	Describe("TryToAssign", func() {
// 		It("sets the correct attributes", func() {
// 			Expect(assignment.TryToAssign()).To(BeTrue())
// 			Expect(assignment.TryToAssign()).To(BeFalse())
// 		})
// 	})

// 	Describe("UnassigneIfExpired", func() {
// 		It("returns the correct result", func() {
// 			assignment.UnassignIfExpired()

// 			Expect(assignment.Id).ToNot(BeNil())

// 			// Modify the time to more than 5 minutes.

// 			// Check again.
// 		})
// 	})

// })
