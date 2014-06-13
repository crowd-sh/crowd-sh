package crowdflow

type Assignment struct {
	SharedAssignment

	AssignDone chan bool `json:"-"`
	Job        *MetaJob  `json:"job"`
}

type Assignments []Assignment

var (
	AvailableAssignments Assignments
)

func NewAssignment(assignDone chan bool, b *Batch, j *MetaJob) {
	assignment := Assignment{
		AssignDone: assignDone,
		Job:        j,
	}

	AvailableAssignments = append(AvailableAssignments, assignment)
}

func (as Assignments) GetUnfinished() *Assignment {
	for _, a := range as {
		if !a.SharedAssignment.Finished && !a.SharedAssignment.Assigned {
			return &a
		}
	}

	return nil
}

func (a Assignment) Finish(val string) {

	// b.TaskConfig.Write(a.Job)
}
