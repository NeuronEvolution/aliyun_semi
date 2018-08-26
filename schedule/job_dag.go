package schedule

func recursiveInitInitJobTimeRangeMin(config *JobConfig) {
	if config.TimeRangeMinInitialized {
		return
	}

	if config.Parents == nil || len(config.Parents) == 0 {
		config.MinStartTime = 0
	} else {
		for _, p := range config.Parents {
			recursiveInitInitJobTimeRangeMin(p)
		}
		for _, p := range config.Parents {
			if config.MinStartTime < p.MinEndTime {
				config.MinStartTime = p.MinEndTime
			}
		}
	}

	config.MinEndTime = config.MinStartTime + config.ExecMinutes
	config.TimeRangeMinInitialized = true
}

func recursiveInitInitJobTimeRangeMax(config *JobConfig) {
	if config.TimeRangeMaxInitialized {
		return
	}

	if config.Children == nil || len(config.Children) == 0 {
		config.MaxEndTime = TimeSampleCount * 15
	} else {
		for _, c := range config.Children {
			recursiveInitInitJobTimeRangeMax(c)
		}
		for _, c := range config.Children {
			if config.MaxEndTime == 0 || config.MaxEndTime > c.MaxStartTime {
				config.MaxEndTime = c.MaxStartTime
			}
		}
	}

	config.MaxStartTime = config.MaxEndTime - config.ExecMinutes
	config.TimeRangeMaxInitialized = true
}

func (r *ResourceManagement) initJobConfigs() {
	for _, config := range r.JobConfigMap {
		if config == nil {
			continue
		}

		recursiveInitInitJobTimeRangeMin(config)
		recursiveInitInitJobTimeRangeMax(config)
	}
}
