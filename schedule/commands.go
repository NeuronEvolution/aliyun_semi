package schedule

type InstanceMoveCommand struct {
	Round      int `json:"round"`
	InstanceId int `json:"instance_id"`
	MachineId  int `json:"machine_id"`
}

type JobDeployCommand struct {
	JobId        string
	MachineId    int
	Count        int
	StartMinutes int
}
