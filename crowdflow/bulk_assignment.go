package crowdflow

type BulkServe struct{}

func (ss BulkServe) Execute(jobs chan MetaJob, j MetaJob) {
	assignment := &BulkAssignment{
		JobsChan: jobs,
		Job:      j,
		Assignment: Assignment{
			Assigned: false,
			Finished: false,
		},
	}

	bulk_assignments = append(bulk_assignments, assignment)
}

type BulkAssignment struct {
	Assignment
	JobsChan chan MetaJob `json:"-"`
	Job      MetaJob      `json:"job"`
}

type BulkAssignments []*BulkAssignment

var (
	bulk_assignments BulkAssignments
)

func (as BulkAssignments) Get() *BulkAssignment {
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

func (as BulkAssignments) Find(id string) *BulkAssignment {
	for _, a := range as {
		if a.Id == id {
			return a
		}
	}

	return nil
}
