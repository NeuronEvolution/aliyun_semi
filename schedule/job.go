package schedule

type JobCommonState struct {
	Jobs      []*Job
	StartTime int
	EndTime   int
}

func (s *JobCommonState) UpdateTime() {
	s.StartTime = TimeSampleCount * 15
	s.EndTime = 0
	for _, job := range s.Jobs {
		if job.StartMinutes == -1 {
			continue
		}

		if job.StartMinutes < s.StartTime {
			s.StartTime = job.StartMinutes
		}

		endTime := job.StartMinutes + job.Config.ExecMinutes
		if endTime > s.EndTime {
			s.EndTime = endTime
		}
	}
}

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
	job := &Job{}
	job.R = r
	job.JobInstanceId = jobInstanceId
	job.Config = config
	job.InstanceCount = instanceCount
	job.Cpu = job.Config.Cpu * float64(job.InstanceCount)
	job.Mem = job.Config.Mem * float64(job.InstanceCount)
	job.StartMinutes = -1

	return job
}

func (job *Job) SetStartMinutes(t int) {
	job.StartMinutes = t
	job.Config.State.UpdateTime()
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
