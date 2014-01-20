package crowdflow

type SplitHtmlServe struct{}

func (ss SplitHtmlServe) Execute(jobs chan Job, j Job) {
	assignment := &SplitAssignment{
		JobsChan: jobs,
		Job:      j,
		Assignment: Assignment{
			Assigned: false,
			Finished: false,
		},
	}

	split_assignments = append(split_assignments, assignment)
}

type SplitAssignment struct {
	Assignment
	JobsChan chan Job `json:"-"`
	Job      Job      `json:"job"`
}

func (as SplitAssignment) Finish(value string) {
	as.Mutex.Lock()
	as.Job.OutputField.Value = value

	as.Finished = true
	as.Mutex.Unlock()

	as.JobsChan <- as.Job
}

type SplitAssignments []*SplitAssignment

var (
	split_assignments SplitAssignments
)

func (as SplitAssignments) Get() *SplitAssignment {
	for _, a := range as {
		a.Mutex.Lock()

		a.UnassignIfExpired()

		if !a.Assigned && !a.Finished {
			defer a.Mutex.Unlock()
			a.Assign()
			return a
		}
		a.Mutex.Unlock()
	}

	return nil
}

func (as SplitAssignments) Find(id string) *SplitAssignment {
	for _, a := range as {
		if a.Id == id {
			return a
		}
	}

	return nil
}
