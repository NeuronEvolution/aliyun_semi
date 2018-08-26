package schedule

type Job struct {
	R             *ResourceManagement
	JobInstanceId int
	Config        *JobConfig
	InstanceCount int
	Cpu           float64
	Mem           float64

	StartMinutes int
}

func NewJob(r *ResourceManagement, jobInstanceId int, config *JobConfig, instanceCount int) *Job {
	j := &Job{}
	j.R = r
	j.JobInstanceId = jobInstanceId
	j.Config = config
	j.InstanceCount = instanceCount
	j.Cpu = j.Config.Cpu * float64(j.InstanceCount)
	j.Mem = j.Config.Mem * float64(j.InstanceCount)

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
