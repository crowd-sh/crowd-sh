package crowdflow

// Interface for serving crowdsourcing jobs.
type Assigner interface {
	// Publish a Batch to be done by workers.
	Execute(aChan chan Assignment, assign Assignment)
}
