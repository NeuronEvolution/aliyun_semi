package schedule

import (
	"fmt"
	"sort"
)

func (r *ResourceManagement) calcJobsMachineNeed(jobs []*Job) (count int) {
	cpu := float64(0)
	for _, job := range jobs {
		cpu += job.Cpu * float64(job.Config.ExecMinutes)
	}
	cpu = cpu / (TimeSampleCount * 15)
	count = int(cpu / float64(46))
	if count < 1 {
		count = 1
	}

	return count
}

func (r *ResourceManagement) bestFitJobs(machines []*Machine, jobs []*Job) (deploy map[*Job]int, restJobs []*Job) {
	return nil, nil
}

func (r *ResourceManagement) jobsScheduleLoop(machines []*Machine) (err error) {
	//按照最早结束时间排序，FF插入
	sort.Slice(r.JobList, func(i, j int) bool {
		job1 := r.JobList[i]
		job2 := r.JobList[j]
		return job1.Config.EndTimeMin < job2.Config.EndTimeMin
	})

	machinesMap := make(map[int]*Machine)
	for _, m := range machines {
		machinesMap[m.MachineId] = m
	}

	currentMachineCount := r.DeployedMachineCount
	scaleCount := r.calcJobsMachineNeed(r.JobList)
	r.log("jobsScheduleLoop init scaleCount=%d\n", scaleCount)
	currentMachineCount += scaleCount

	for {
		tempMachines := MachinesClone(machines)
		deploy, restJobs := r.bestFitJobs(tempMachines[:currentMachineCount], r.JobList)
		if len(restJobs) == 0 {
			for job, machineId := range deploy {
				fmt.Println(job, machineId)
			}
			break
		}
		if currentMachineCount == len(tempMachines) {
			return fmt.Errorf("jobsScheduleLoop failed")
		}

		scaleCount := r.calcJobsMachineNeed(restJobs)
		currentMachineCount += scaleCount
		if currentMachineCount > len(tempMachines) {
			currentMachineCount = len(tempMachines)
		}
		r.log("jobsScheduleLoop scaleCount=%d,currentMachineCount=%d\n", scaleCount, currentMachineCount)
	}

	return nil
}
