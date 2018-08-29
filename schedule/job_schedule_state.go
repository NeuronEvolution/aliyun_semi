package schedule

type JobScheduleState struct {
	Jobs      []*Job
	StartTime int
	EndTime   int
}

func (s *JobScheduleState) UpdateTime() {
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
