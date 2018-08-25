package schedule

type Replay struct {
	R                    *ResourceManagement
	InstanceMoveCommands []*InstanceMoveCommand
	JobDeployCommands    []*JobDeployCommand
}

func NewReplay(
	r *ResourceManagement,
	instanceMoveCommands []*InstanceMoveCommand,
	jobDeployCommands []*JobDeployCommand) (replay *Replay) {

	replay = &Replay{}
	replay.R = r
	replay.InstanceMoveCommands = instanceMoveCommands
	replay.JobDeployCommands = jobDeployCommands

	return replay
}

func (r *Replay) Run() (err error) {
	return nil
}
