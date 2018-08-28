package schedule

import "fmt"

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

func (r *Replay) Run(finalState []*Machine) (err error) {
	//创建机器
	machines := make(map[int]*Machine)
	for _, m := range r.R.MachineList {
		machine := NewMachine(m.R, m.MachineId, m.Config)
		machines[machine.MachineId] = machine
	}

	//初始化
	deploys := make(map[int]*Machine)
	for _, config := range r.R.InstanceDeployConfigList {
		instance := r.R.InstanceMap[config.InstanceId]
		m := machines[config.MachineId]
		m.AddInstance(instance)
		deploys[instance.InstanceId] = m
	}

	totalScore := float64(0)
	for _, m := range machines {
		totalScore += m.GetCpuCost()
	}
	r.R.log("replay 1 score=%f\n", totalScore)

	for _, move := range r.InstanceMoveCommands {
		//fmt.Println(move.Round, move.InstanceId, move.MachineId)
		instance := r.R.InstanceMap[move.InstanceId]
		deploys[instance.InstanceId].RemoveInstance(instance.InstanceId)
		m := machines[move.MachineId]
		if !m.ConstraintCheck(instance, 1) {
			return fmt.Errorf("replay ConstraintCheck failed machineId=%d,instanceId=%d", m.MachineId, instance.InstanceId)
		}
		m.AddInstance(instance)
		deploys[instance.InstanceId] = m
	}

	for _, m := range machines {
		m.beginOffline()
	}

	totalScore = float64(0)
	for _, m := range machines {
		totalScore += m.GetCpuCostReal()
	}
	r.R.log("replay 2 score=%f\n", totalScore)

	for _, v := range r.JobDeployCommands {
		m := machines[v.MachineId]
		var config *JobConfig
		for _, c := range r.R.JobConfigMap {
			if c != nil && c.RealJobId == v.JobId {
				config = c
				break
			}
		}

		for i := v.StartMinutes; i < v.StartMinutes+config.ExecMinutes; i++ {
			m.Cpu[i] += config.Cpu * float64(v.Count)
			m.Mem[i] += config.Mem * float64(v.Count)
		}
	}

	totalScore = float64(0)
	for _, m := range machines {
		totalScore += m.GetCpuCostReal()
	}

	r.R.log("Replay ok,total score=%f\n", totalScore)

	return nil
}
