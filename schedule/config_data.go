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
	Parents       []*JobConfig
	Children      []*JobConfig
}

func (c *JobConfig) getPackCount() (count int) {
	maxCpu := float64(16)
	maxMem := float64(32)

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
