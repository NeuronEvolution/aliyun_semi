package schedule

import (
	"fmt"
	"sort"
)

type JobScheduler struct {
	R        *ResourceManagement
	Machines []*Machine
}

func NewJobScheduler(r *ResourceManagement, machines []*Machine) (s *JobScheduler) {
	s = &JobScheduler{}
	s.R = r
	s.Machines = machines

	return s
}

func (s *JobScheduler) bestFitJobs(machines []*Machine, jobs []*Job) (result []*Machine, err error) {
	//复制机器
	result = MachinesCloneWithInstances(machines)

	//调度状态
	scheduleState := NewJobScheduleState(s.R, s.R.JobList)
	for _, job := range s.R.JobList {
		job.StartMinutes = -1
	}

	//BFD
	for i, job := range jobs {
		if i > 0 && i%10000 == 0 {
			s.R.log("bestFitJobs %d\n", i)
		}

		min := TimeSampleCount * 15
		var minMachine *Machine
		startTimeMin, startTimeMax, endTimeMin, endTimeMax := job.GetTimeRange(scheduleState)
		for _, m := range result {
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
			return nil, fmt.Errorf("bestFitJobs failed")
		}
		job.SetStartMinutes(min)
		scheduleState[job.Config.JobId].UpdateTime()
		minMachine.AddJob(job)
	}

	return result, nil
}

func (s *JobScheduler) Run() (err error) {
	s.R.log("JobScheduler.Run")
	if len(s.R.JobList) == 0 {
		return nil
	}

	//按照最早结束时间排序，FF插入
	sort.Slice(s.R.JobList, func(i, j int) bool {
		job1 := s.R.JobList[i]
		job2 := s.R.JobList[j]
		if job1.Config.EndTimeMin == job2.Config.EndTimeMin {
			return job1.Cpu*float64(job1.Config.ExecMinutes) > job2.Cpu*float64(job1.Config.ExecMinutes)
		} else {
			return job1.Config.EndTimeMin < job2.Config.EndTimeMin
		}
	})

	var lastResult []*Machine
	lastSucceed := false
	scaleCurrent := s.R.DeployedMachineCount
	scaleUp := true
	scaleDividing := false
	scaleStep := 512
	machineCount := 0
	for {
		//如果已经开始二分搜索，区间大小减半
		if scaleDividing {
			scaleStep /= 2
		}

		//正向或反向搜索
		if scaleUp {
			s.R.log("JobScheduler.Run scale up last=%d,scaleStep=%d,now=%d\n",
				scaleCurrent, scaleStep, scaleCurrent+scaleStep)
			scaleCurrent += scaleStep

		} else {
			s.R.log("JobScheduler.Run scale down last=%d,scaleStep=%d,now=%d\n",
				scaleCurrent, scaleStep, scaleCurrent-scaleStep)
			scaleCurrent -= scaleStep
		}

		//已到最大机器数//todo这里需要调整步长，暂不优化
		if scaleCurrent > len(s.Machines) {
			scaleCurrent = len(s.Machines)
			s.R.log("JobScheduler.Run reach max scaleCurrent=%d\n", scaleCurrent)
		}

		result, err := s.bestFitJobs(s.Machines[:scaleCurrent], s.R.JobList)
		if err != nil {
			s.R.log("JobScheduler.Run failed scaleCurrent=%d\n", scaleCurrent)
			//已达到最大机器数，并且调度失败
			if scaleCurrent == len(s.Machines) {
				return fmt.Errorf("JobScheduler.Run failed,max machine used")
			}

			//部署失败，减少机器数量
			if lastSucceed {
				//开始分割
				scaleDividing = true
				//反向
				scaleUp = !scaleUp
			}
			lastSucceed = false

			//已分割完
			if scaleStep == 1 {
				s.R.log("JobScheduler.Run scaleStep=1 failed\n")
				machineCount = scaleCurrent + 1
				//todo这里可以优化
				lastResult, err = s.bestFitJobs(s.Machines[:machineCount], s.R.JobList)
				if err != nil {
					panic("bestFitJobs last one failed")
				}
				break
			}
		} else {
			s.R.log("JobScheduler.Run succeed scaleCurrent=%d\n", scaleCurrent)
			if !lastSucceed {
				//开始分割
				scaleDividing = true
				//反向
				scaleUp = !scaleUp
			}
			lastSucceed = true

			//保存最后成功结果
			lastResult = result

			//已分割完
			if scaleStep == 1 {
				s.R.log("JobScheduler.Run scaleStep=1 ok\n")
				machineCount = scaleCurrent
				break
			}
		}
	}

	for i, m := range lastResult {
		for _, job := range m.JobList[:m.JobListCount] {
			s.Machines[i].AddJob(job)
		}
	}

	s.R.log("JobScheduler.Run totalScore=%f,machineCount=%d\n", MachinesGetScore(s.Machines), machineCount)

	return nil
}
