package schedule

type Job struct {
	R             *ResourceManagement
	JobInstanceId int
	Config        *JobConfig
	StartMinutes  int
}

func NewJob(r *ResourceManagement, jobInstanceId int, config *JobConfig) *Job {
	j := &Job{}
	j.R = r
	j.JobInstanceId = jobInstanceId
	j.Config = config

	return j
}

func JobsCopy(p []*Job) (r []*Job) {
	if p == nil {
		return nil
	}

	r = make([]*Job, len(p))
	for i, v := range p {
		r[i] = v
	}

	return r
}
