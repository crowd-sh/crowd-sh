package crowdflow

func NewSplitAssigner(assignDone chan bool, a *Assignment) {
	<-assignDone
	// No-op
}

type SplitAssignment struct {
	SharedAssignment

	Assignment     *Assignment
	AssignmentDone chan bool

	OutputField *JobField `json:"output"`
}

func (as SplitAssignment) Finish(value string) {
	as.SharedAssignment.Mutex.Lock()

	as.OutputField.Value = value
	as.SharedAssignment.Finished = true

	as.SharedAssignment.Mutex.Unlock()

	as.AssignmentDone <- true
}

type SplitAssignments []SplitAssignment

var (
	AvailableSplitAssignments SplitAssignments
)

func GetSplitAssignment() (assign *SplitAssignment) {
	if len(AvailableSplitAssignments) > 0 {
		for _, a := range AvailableSplitAssignments {
			if a.SharedAssignment.TryToAssign() {
				return &a
			}
		}
	} else {
		// Create a new split assignment.

		available := AvailableAssignments.GetUnfinished()
		if available != nil {
			for _, f := range available.Job.OutputFields {
				ss := SplitAssignment{
					Assignment:  available,
					OutputField: &f,
				}

				AvailableSplitAssignments = append(AvailableSplitAssignments, ss)
			}

			return GetSplitAssignment()
		}

	}

	return nil
}

func FindSplitAssignment(id string) *SplitAssignment {
	for _, a := range AvailableSplitAssignments {
		if a.Id == id {
			return &a
		}
	}

	return nil
}
