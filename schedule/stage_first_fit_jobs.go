package schedule

import (
	"fmt"
	"sort"
)

func (r *ResourceManagement) firstFitJobs(machines []*Machine) (err error) {
	//按照最早结束时间排序，FF插入
	sort.Slice(r.JobList, func(i, j int) bool {
		job1 := r.JobList[i]
		job2 := r.JobList[j]
		return job1.Config.EndTimeMin < job2.Config.EndTimeMin
	})

	for i, job := range r.JobList {
		if i > 0 && i%1000 == 0 {
			r.log("firstFitJobs %d\n", i)
		}

		startTimeMin, startTimeMax, endTimeMin, endTimeMax := job.GetTimeRange()
		deployed := false
		for machineIndex, m := range machines {
			if m.InstanceListCount > 0 {
				continue
			}
			ok, startMinutes := m.CanFirstFitJob(job, startTimeMin, startTimeMax, endTimeMin, endTimeMax)
			if ok {
				job.SetStartMinutes(startMinutes)
				m.AddJob(job)
				r.JobDeployMap[job.JobInstanceId] = m
				//更新部署数量
				if machineIndex+1 > r.DeployedMachineCount {
					r.DeployedMachineCount = machineIndex + 1
				}
				deployed = true
				break
			}
		}
		if !deployed {
			return fmt.Errorf(fmt.Sprintf("firstFitJobs failed,%d jobInstanceId=%d,instanceCount=%d,cpu=%f,mem=%f",
				i, job.JobInstanceId, job.InstanceCount, job.Cpu, job.Mem))
		}
	}

	totalScore := float64(0)
	for _, m := range machines[:r.DeployedMachineCount] {
		totalScore += m.GetCpuCostReal()
	}

	r.log("firstFitJobs deployedMachineCount=%d,score=%f\n", r.DeployedMachineCount, totalScore)

	return nil
}
