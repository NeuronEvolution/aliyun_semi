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

func (r *ResourceManagement) jobSchedule() (err error) {
	r.log("jobSchedule start\n")

	//之后实例不再调度，先计算出实例迁移指令
	instanceMoveCommands, err := NewOnlineMerge(r).Run()
	if err != nil {
		return err
	}

	//todo 这里需要考虑在线迁移时的实例交换,改为从初始状态迁移后再部署任务,暂时不需要优化，除了e数据不需要固定实例
	//重新插入实例，避免浮点精度问题
	machines := MachinesCloneWithInstances(r.MachineList)
	r.log("jobSchedule init totalScore=%f\n", MachinesGetScore(machines))

	//任务调度
	//err = r.firstFitJobs(machines)
	err = NewJobScheduler(r, machines).Run()
	if err != nil {
		return err
	}
	r.log("jobSchedule totalScore=%f\n", MachinesGetScore(machines))

	//构造任务调度指令
	jobDeployCommands := r.buildJobDeployCommands(machines)

	//验证结果
	err = NewReplay(r, instanceMoveCommands, jobDeployCommands).Run(machines)
	if err != nil {
		return nil
	}

	//输出结果
	return r.output(machines, instanceMoveCommands, jobDeployCommands)
}
