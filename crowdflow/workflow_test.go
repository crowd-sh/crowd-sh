package crowdflow_test

import (
	. "github.com/abhiyerra/workmachine/crowdflow"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Workflow", func() {

	var (
		tasks    []Research
		research TaskDesc
		batch    *Batch
	)

	BeforeEach(func() {
		tasks = []Research{
			Research{Name: "Michelangelo"},
			Research{Name: "Leonardo"},
		}

		research = TaskDesc{
			Title:       "Research the fields",
			Description: "Research the person.",
			Price:       1,
			Tasks:       tasks,
			Write:       func(j *MetaJob) {},
		}

		batch = NewBatch(research)
	})

	Describe("Batch", func() {
		Describe("TaskDesc", func() {
			It("has the correct TaskDesc", func() {
				Expect(batch.TaskDesc.Title).To(Equal("Research the fields"))
				Expect(batch.TaskDesc.Description).To(Equal("Research the person."))
				Expect(batch.TaskDesc.Price).To(Equal(uint(1)))
			})
		})

		Describe("MetaJobs", func() {
			It("has the correct MetaJobs", func() {
				Expect(len(batch.MetaJobs)).To(Equal(len(tasks)))
			})

			It("has the correct titles and descriptions", func() {
				Expect(len(batch.MetaJobs)).To(Equal(len(tasks)))

				for _, m := range batch.MetaJobs {
					Expect(m.Title).To(Equal("Research the fields"))
					Expect(m.Description).To(Equal("Research the person."))
				}
			})

			It("has the correct input fields", func() {
				m := batch.MetaJobs

				Expect(m[0].InputFields[0].Id).To(Equal("name"))
				Expect(m[0].InputFields[0].Value).To(Equal("Michelangelo"))
				Expect(m[0].InputFields[0].Description).To(Equal("The name of the person to research."))
				Expect(m[0].InputFields[0].Type).To(Equal(JobFieldType("")))

				Expect(m[1].InputFields[0].Id).To(Equal("name"))
				Expect(m[1].InputFields[0].Value).To(Equal("Leonardo"))
				Expect(m[1].InputFields[0].Description).To(Equal("The name of the person to research."))
				Expect(m[1].InputFields[0].Type).To(Equal(JobFieldType("")))
			})

			It("has the correct output fields", func() {
				for _, m := range batch.MetaJobs {
					Expect(m.OutputFields[0].Id).To(Equal("born"))
					Expect(m.OutputFields[1].Id).To(Equal("is_painter"))
					Expect(m.OutputFields[1].Value).To(Equal(""))
					Expect(m.OutputFields[1].Description).To(Equal("Was the person a painter?"))
					Expect(m.OutputFields[1].Type).To(Equal(JobFieldType("checkbox")))
				}

			})

			It("has an emptt Jobs field", func() {
				for _, m := range batch.MetaJobs {
					Expect(len(m.Jobs)).To(Equal(0))
				}
			})
		})

		PDescribe("Run", func() {
			It("Starts a new running job", func() {
				// batch.Run()
				// Has correct number of jobs
			})

			PDescribe("StartJobs", func() {

			})

			PDescribe("Job", func() {

			})
		})

	})
})
