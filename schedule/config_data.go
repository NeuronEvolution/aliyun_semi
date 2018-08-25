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
