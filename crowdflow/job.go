package crowdflow

// Interface for serving crowdsourcing jobs.
type SplitAssigner interface {
	// Publish a Batch to be done by workers.
	Execute(jobs chan Job, j Job)
}

// Interface for serving crowdsourcing jobs.
type Assigner interface {
	// Publish a Batch to be done by workers.
	Execute(jobs chan MetaJob, j MetaJob)
}

type Job struct {
	Parent      *MetaJob  `json:"info"`
	OutputField *JobField `json:"output"`
}
