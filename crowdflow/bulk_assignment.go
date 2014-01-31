package crowdflow

func GetFormAssignment() *Assignment {
	for _, a := range AvailableAssignments {
		if a.SharedAssignment.TryToAssign() {
			return &a
		}
	}

	return nil
}

func FindFormAssignment(id string) *Assignment {
	for _, a := range AvailableAssignments {
		if a.Id == id {
			return &a
		}
	}

	return nil
}
