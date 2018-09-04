package schedule

type JobMerge struct {
	R             *ResourceManagement
	Machines      []*Machine
	ScheduleState []*JobScheduleState
}

func NewJobMerge(r *ResourceManagement, machines []*Machine, scheduleState []*JobScheduleState) *JobMerge {
	s := &JobMerge{}
	s.R = r
	s.Machines = machines
	s.ScheduleState = scheduleState

	return s
}

func (s *JobMerge) Run() (jobDeployCommands []*JobDeployCommand) {

	return s.R.buildJobDeployCommands(s.Machines)
}
