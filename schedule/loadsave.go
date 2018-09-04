package schedule

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func (r *ResourceManagement) getInstanceSaveFilepath() string {
	return r.OutputDir + fmt.Sprintf("/save_%d.json", r.GetDatasetMachineCount())
}

func (r *ResourceManagement) loadInstanceMoveCommands() (moveCommands []*InstanceMoveCommand, err error) {
	data, err := ioutil.ReadFile(r.getInstanceSaveFilepath())
	if err != nil {
		r.log("loadInstanceMoveCommands ReadFile failed,%s\n", err.Error())
		return
	}

	err = json.Unmarshal(data, &moveCommands)
	if err != nil {
		r.log("loadInstanceMoveCommands json.Unmarshal failed,%s\n", err.Error())
		return
	}

	//初始状态部署
	for _, config := range r.InstanceDeployConfigList {
		instance := r.InstanceMap[config.InstanceId]
		m := r.MachineMap[config.MachineId]
		m.AddInstance(instance)
		r.DeployMap[instance.InstanceId] = m
	}

	//迁移
	for _, move := range moveCommands {
		//fmt.Println(move.Round, move.InstanceId, move.MachineId)
		instance := r.InstanceMap[move.InstanceId]
		r.DeployMap[instance.InstanceId].RemoveInstance(instance.InstanceId)
		m := r.MachineMap[move.MachineId]
		if !m.ConstraintCheck(instance, 1) {
			return nil,
				fmt.Errorf("loadInstanceMoveCommands ConstraintCheck failed machineId=%d,instanceId=%d",
					m.MachineId, instance.InstanceId)
		}
		m.AddInstance(instance)
		r.DeployMap[instance.InstanceId] = m
	}

	r.log("loadInstanceMoveCommands ok,totalScore=%f,file=%s\n", MachinesGetScore(r.MachineList), r.getInstanceSaveFilepath())

	return moveCommands, nil
}

func (r *ResourceManagement) saveInstanceMoveCommands(moveCommands []*InstanceMoveCommand) {
	data, err := json.Marshal(moveCommands)
	if err != nil {
		r.log("saveInstanceMoveCommands failed,%s\n", err.Error())
		return
	}

	err = ioutil.WriteFile(r.getInstanceSaveFilepath(), data, os.ModePerm)
	if err != nil {
		r.log("saveInstanceMoveCommands failed,%s\n", err.Error())
		return
	}

	r.log("saveInstanceMoveCommands ok,commandCount=%d,file=%s\n", len(moveCommands), r.getInstanceSaveFilepath())
}

func (r *ResourceManagement) buildJobDeployCommands(machines []*Machine) (commands []*JobDeployCommand) {
	commands = make([]*JobDeployCommand, 0)
	for _, m := range machines {
		if m.JobListCount == 0 {
			continue
		}

		for _, job := range m.JobList[:m.JobListCount] {
			commands = append(commands, &JobDeployCommand{
				JobInstanceId: job.JobInstanceId,
				JobId:         job.Config.RealJobId,
				MachineId:     m.MachineId,
				Count:         job.InstanceCount,
				StartMinutes:  job.StartMinutes,
			})
		}
	}

	return commands
}

func (r *ResourceManagement) getJobDeploySaveFilepath() string {
	return r.OutputDir + fmt.Sprintf("/save_%d_job_%f_%d_%d_%d.json",
		r.GetDatasetMachineCount(), JobScheduleCpuLimitStep, JobPackCpu, JobPackMem, JobPackLimit)
}

func (r *ResourceManagement) loadJobDeployCommands(
	machines []*Machine, scheduleState []*JobScheduleState) (
	commands []*JobDeployCommand, err error) {
	data, err := ioutil.ReadFile(r.getJobDeploySaveFilepath())
	if err != nil {
		r.log("loadJobDeployCommands ReadFile failed,%s\n", err.Error())
		return
	}

	err = json.Unmarshal(data, &commands)
	if err != nil {
		r.log("loadJobDeployCommands json.Unmarshal failed,%s\n", err.Error())
		return
	}

	for _, cmd := range commands {
		job := r.JobMap[cmd.JobInstanceId]
		job.StartMinutes = cmd.StartMinutes
		scheduleState[job.Config.JobId].UpdateTime()
		for _, m := range machines {
			if m.MachineId == cmd.MachineId {
				m.AddJob(job)
				break
			}
		}
	}

	r.log("loadJobDeployCommands ok,totalScore=%f,file=%s\n", MachinesGetScore(machines), r.getJobDeploySaveFilepath())

	return commands, nil
}

func (r *ResourceManagement) saveJobDeployCommands(commands []*JobDeployCommand) {
	data, err := json.Marshal(commands)
	if err != nil {
		r.log("saveJobDeployCommands failed,%s\n", err.Error())
		return
	}

	err = ioutil.WriteFile(r.getJobDeploySaveFilepath(), data, os.ModePerm)
	if err != nil {
		r.log("saveJobDeployCommands failed,%s\n", err.Error())
		return
	}

	r.log("saveJobDeployCommands ok,commandCount=%d,file=%s\n", len(commands), r.getJobDeploySaveFilepath())
}
