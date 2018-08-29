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

func (r *ResourceManagement) bestFitJobs(machines []*Machine, jobs []*Job, scheduleState []*JobScheduleState) (
	ok bool, deploy map[*Job]int, restJobs []*Job) {
	for i, job := range jobs {
		if i > 0 && i%1000 == 0 {
			r.log("bestFitJobs %d\n", i)
		}

		min := TimeSampleCount * 15
		var minMachine *Machine
		startTimeMin, startTimeMax, endTimeMin, endTimeMax := job.GetTimeRange(scheduleState)
		for _, m := range machines {
			ok, startTime := m.CanFirstFitJob(job, startTimeMin, startTimeMax, endTimeMin, endTimeMax)
			if !ok {
				continue
			}
			if startTime < min {
				min = startTime
				minMachine = m
			}
		}
		if minMachine == nil {
			return false, nil, JobsCopy(jobs[i:])
		}
		job.SetStartMinutes(min)
		scheduleState[job.Config.JobId].UpdateTime()
		minMachine.AddJob(job)
	}

	return true, nil, nil
}

func (r *ResourceManagement) jobsScheduleLoop(machines []*Machine) (err error) {
	r.log("jobsScheduleLoop start\n")

	//按照最早结束时间排序，FF插入
	sort.Slice(r.JobList, func(i, j int) bool {
		job1 := r.JobList[i]
		job2 := r.JobList[j]
		if job1.Config.EndTimeMin == job2.Config.EndTimeMin {
			return job1.Cpu*float64(job1.Config.ExecMinutes) > job2.Cpu*float64(job1.Config.ExecMinutes)
		} else {
			return job1.Config.EndTimeMin < job2.Config.EndTimeMin
		}
	})

	scheduleState := NewJobScheduleState(r, r.JobList)
	for _, job := range r.JobList {
		job.StartMinutes = -1
	}
	ok, _, restJobs := r.bestFitJobs(machines[:r.DeployedMachineCount+772], r.JobList, scheduleState)
	if !ok {
		fmt.Printf("bestFitJobs failed restJobs=%d\n", len(restJobs))
		panic("aaaa")
	}

	totalScore := float64(0)
	for _, m := range machines {
		if m.JobListCount > 0 {
			totalScore += m.GetCpuCost()
		}
	}
	r.log("jobsScheduleLoop totalScore=%f\n", totalScore)

	return

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
		ok, deploy, restJobs := r.bestFitJobs(tempMachines[:currentMachineCount], r.JobList, scheduleState)
		if ok {
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
