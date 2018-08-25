package schedule

type InstanceMoveCommand struct {
	Round      int
	InstanceId int
	MachineId  int
}

type JobDeployCommand struct {
	JobId        string
	MachineId    int
	Count        int
	StartMinutes int
}
