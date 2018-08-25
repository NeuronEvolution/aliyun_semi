package semitemp

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
	MachineLevelConfig
}
