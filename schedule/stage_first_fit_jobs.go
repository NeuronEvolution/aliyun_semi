package schedule

import (
	"fmt"
	"sort"
)

func (r *ResourceManagement) firstFitJobs() (err error) {
	sort.Slice(r.JobList, func(i, j int) bool {
		job1 := r.JobList[i]
		job2 := r.JobList[j]
		return job1.Config.EndTimeMin < job2.Config.EndTimeMin
	})

	for i, job := range r.JobList {
		if i > 0 && i%10000 == 0 {
			r.log("firstFitJobs %d\n", i)
		}

		startTimeMin, startTimeMax, _, _ := job.GetTimeRange()

		deployed := false
		for machineIndex, m := range r.MachineList {
			ok, startMinutes := m.CanFirstFitJob(job, startTimeMin, startTimeMax)
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

	r.log("firstFitJobs deployedMachineCount=%d,score=%f\n", r.DeployedMachineCount, r.CalcTotalScoreReal())

	return nil
}
