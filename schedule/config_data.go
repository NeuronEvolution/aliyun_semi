package schedule

type AppInterferenceConfig struct {
	AppId1       int
	AppId2       int
	Interference int
}

type AppResourcesConfig struct {
	AppId int
	Resource

	InferenceAppCount int
}

type InstanceDeployConfig struct {
	InstanceId int
	AppId      int
	MachineId  int
}

type MachineResourcesConfig struct {
	MachineId int
	MachineConfig
}

type JobConfig struct {
	JobId         int
	RealJobId     string
	Cpu           float64
	Mem           float64
	InstanceCount int
	ExecMinutes   int
	PreJobs       []int

	Parents                 []*JobConfig
	Children                []*JobConfig
	StartTimeMin            int
	StartTimeMax            int
	EndTimeMin              int
	EndTimeMax              int
	TimeRangeMinInitialized bool
	TimeRangeMaxInitialized bool
}

//任务打包部署
func (c *JobConfig) getPackCount(totalJobCount int) (count int) {
	if totalJobCount < 320000 {
		return 1
	}

	maxCpu := float64(JobPackCpu)
	maxMem := float64(JobPackMem)

	if c.Cpu >= maxCpu {
		return 1
	}

	if c.Cpu >= maxMem {
		return 1
	}

	count = 32
	cpuCount := int(maxCpu / c.Cpu)
	memCount := int(maxMem / c.Mem)
	if cpuCount < count {
		count = cpuCount
	}
	if memCount < count {
		count = memCount
	}

	return count
}

func (c *JobConfig) isParentOf(p *JobConfig) bool {
	for _, v := range c.Children {
		if v == p {
			return true
		}
	}

	return false
}

func (c *JobConfig) isChildOf(p *JobConfig) bool {
	for _, v := range c.Parents {
		if v == p {
			return true
		}
	}

	return false
}
