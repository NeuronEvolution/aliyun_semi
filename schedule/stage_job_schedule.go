package schedule

func (r *ResourceManagement) buildJobDeployCommands(machines []*Machine) (commands []*JobDeployCommand) {
	commands = make([]*JobDeployCommand, 0)
	for _, m := range machines {
		if m.JobListCount == 0 {
			continue
		}

		for _, job := range m.JobList[:m.JobListCount] {
			commands = append(commands, &JobDeployCommand{
				JobId:        job.Config.RealJobId,
				MachineId:    m.MachineId,
				Count:        job.InstanceCount,
				StartMinutes: job.StartMinutes,
			})
		}
	}

	return commands
}
