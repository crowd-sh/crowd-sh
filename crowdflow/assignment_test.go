package crowdflow

import (
//	. "github.com/onsi/ginkgo"
//	. "github.com/onsi/gomega"
)

// var _ = Describe("Assignments", func() {
// 	var (
// 		AvailAssignments Assignments
// 		EmptyAssignments Assignments
// 		Assign1          Assignment
// 	)

// 	BeforeEach(func() {
// 		Assign1 = Assignment{
// 			SharedAssignment: SharedAssignment{
// 				Assigned: false,
// 			},
// 		}

// 		AvailAssignments = append(AvailAssignments, Assign1)
// 	})

// 	Describe("GetUnfinished", func() {
// 		It("returns an unassigned assignment if available", func() {
// 			var Assign *Assignment = AvailAssignments.GetUnfinished()
// 			Expect(Assign).To(Equal(&Assign1))

// 			Assign1.SharedAssignment.Finished = true
// 			Assign1.SharedAssignment.Assigned = true

// 			Expect(AvailableAssignments.GetUnfinished()).To(BeNil())
// 		})

// 		It("returns nil if assignment if unavailable", func() {
// 			Expect(EmptyAssignments.GetUnfinished()).To(BeNil())
// 		})

// 	})
// })

// var _ = Describe("Assignment", func() {
// 	Describe("NewAssignment", func() {
// 		It("adds the assignment to AvailableAssignments", func() {
// 			oldLen := len(AvailableAssignments)

// 			batch := Batch{}
// 			meta_job := MetaJob{}
// 			assignDone := make(chan bool)

// 			NewAssignment(assignDone, &batch, &meta_job)
// 			Expect(len(AvailableAssignments)).To(Equal(oldLen + 1))
// 		})
// 	})

// 	PDescribe("Execute", func() {

// 	})

// 	PDescribe("Finish", func() {

// 	})
// })
