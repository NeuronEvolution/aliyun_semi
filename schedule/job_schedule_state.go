package schedule

type JobScheduleState struct {
	Jobs      []*Job
	StartTime int
	EndTime   int

	StartTimeMin int
	StartTimeMax int
	EndTimeMin   int
	EndTimeMax   int
}

type JobScheduleStateManager struct {
	States []*JobScheduleState
}

func NewJobScheduleStateManager(r *ResourceManagement, jobs []*Job) (m *JobScheduleStateManager) {
	m = &JobScheduleStateManager{}
	m.States = NewJobScheduleState(r, jobs)

	return m
}

func NewJobScheduleState(r *ResourceManagement, jobs []*Job) (result []*JobScheduleState) {
	result = make([]*JobScheduleState, len(r.JobConfigMap))
	for i := 0; i < len(result); i++ {
		config := r.JobConfigMap[i]
		if config != nil {
			s := &JobScheduleState{}
			s.StartTime = TimeSampleCount * 15
			s.EndTime = 0
			s.StartTimeMin = 0
			s.EndTimeMin = config.ExecMinutes
			s.EndTimeMax = TimeSampleCount * 15
			s.StartTimeMax = s.EndTimeMax - config.ExecMinutes
			result[i] = s
		}
	}

	for _, job := range jobs {
		result[job.Config.JobId].Jobs = append(result[job.Config.JobId].Jobs, job)
	}

	return result
}

func (m *JobScheduleStateManager) UpdateTime(job *Job) {
	s := m.States[job.Config.JobId]
	//更新自身时间区间
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

	//更新所有父节点的最大时间
	if s.StartTime < TimeSampleCount*15 && job.Config.Parents != nil {
		for _, p := range job.Config.Parents {
			if p.EndTimeMax > s.StartTime {
				p.EndTimeMax = s.StartTime
				p.StartTimeMax = p.EndTimeMax - p.ExecMinutes
			}
		}
	}

	//更新所有子节点的最小时间
	if s.EndTime > 0 && job.Config.Children != nil {
		for _, c := range job.Config.Children {
			if c.StartTimeMin < s.EndTime {
				c.StartTimeMin = s.EndTime
				c.EndTimeMin = c.StartTimeMin + c.ExecMinutes
			}
		}
	}
}
