package schedule

import (
	"math"
	"sort"
	"sync"
)

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

func (s *JobMerge) bestFit(machines []*Machine, job *Job) (bestMachine *Machine, bestStartTime int, minScoreAdd float64) {
	var minScoreAddMachine *Machine
	minScoreAdd = math.MaxFloat64
	bestStartTime = TimeSampleCount * 15
	startTimeMin, startTimeMax, _, _ := job.RecursiveGetTimeRange(s.ScheduleState)
	for _, m := range machines {
		ok, startMinutes, scoreAdd := m.bestFitJob(job, startTimeMin, startTimeMax)
		if !ok {
			continue
		}

		update := false
		if scoreAdd < minScoreAdd {
			//fmt.Println("bestFit scoreAdd < minScoreAdd", scoreAdd, minScoreAdd, startMinutes, m.MachineId)
			update = true
		} else if scoreAdd == minScoreAdd {
			//fmt.Println("bestFit scoreAdd == minScoreAdd", scoreAdd, minScoreAdd, startMinutes, m.MachineId)
			if startMinutes < bestStartTime {
				update = true
			}
		}
		if update {
			minScoreAdd = scoreAdd
			bestStartTime = startMinutes
			minScoreAddMachine = m
		}
	}

	//fmt.Println("bestFit", minScoreAdd, bestStartTime, minScoreAddMachine.MachineId)

	return minScoreAddMachine, bestStartTime, minScoreAdd
}

func (s *JobMerge) parallelBestFit(
	machines []*Machine, job *Job) (
	bestMachine *Machine, bestStartTime int, minScoreAdd float64) {

	minScoreAdd = math.MaxFloat64

	//分割机器，并发BestFit
	const parallelCount = ParallelCpuCount * 2
	var bestMachineList [parallelCount]*Machine
	var bestStartTimeList [parallelCount]int
	var minScoreAddList [parallelCount]float64
	size := len(machines) / parallelCount

	wg := &sync.WaitGroup{}
	for pI := 0; pI < parallelCount; pI++ {
		start := pI * size
		count := size
		if pI == parallelCount-1 {
			count = len(machines) - start
		}
		wg.Add(1)
		go func(index int, subMachines []*Machine) {
			defer wg.Done()

			bestMachineList[index], bestStartTimeList[index], minScoreAddList[index] = s.bestFit(subMachines, job)
		}(pI, machines[start:start+count])
	}
	wg.Wait()
	for i, v := range bestMachineList {
		if v == nil {
			continue
		}

		if minScoreAddList[i] < minScoreAdd {
			minScoreAdd = minScoreAddList[i]
			bestMachine = bestMachineList[i]
			bestStartTime = bestStartTimeList[i]
		}
	}

	return bestMachine, bestStartTime, minScoreAdd
}

func (s *JobMerge) Run(outputCallback func() (err error)) (err error) {
	s.R.log("JobMerge.Run totalScore=%f\n", MachinesGetScore(s.Machines))

	for {
		moved := false
		for i, m := range s.Machines {
			if i > 0 && i%100 == 0 {
				s.R.log("JobMerge.Run %d,totalScore=%f\n", i, MachinesGetScore(s.Machines))
			}

			if m.JobListCount == 0 {
				continue
			}

			//获取每个机器cpu最高且部署了任务的时刻
			maxCpu, _, jobs := m.GetMaxCpuTimeWithJobs()
			if maxCpu <= m.Config.Cpu*0.5 {
				//fmt.Println("merge small")
				continue
			}

			//对job按面积排序
			sort.Slice(jobs, func(i, j int) bool {
				return jobs[i].Cpu*float64(jobs[i].Config.ExecMinutes) > jobs[j].Cpu*float64(jobs[j].Config.ExecMinutes)
			})

			for _, job := range jobs {
				m.RemoveJob(job.JobInstanceId)
				//fmt.Println("remove", MachinesGetScore(s.Machines))
				bestMachine, bestStartTime, _ := s.parallelBestFit(s.Machines, job)
				//跳过最佳位置是原来的位置
				if bestMachine == m && bestStartTime == job.StartMinutes {
					m.AddJob(job)
					//fmt.Println("merge keep", MachinesGetScore(s.Machines))
					continue
				}

				//fmt.Println("merge new", job.StartMinutes, bestStartTime, m.MachineId, bestMachine.MachineId)

				//迁移job到最佳位置
				job.StartMinutes = bestStartTime
				s.ScheduleState[job.Config.JobId].UpdateTime()
				bestMachine.AddJob(job)
				moved = true

				//fmt.Println("merge new", MachinesGetScore(s.Machines))

				//每轮只处理一个，避免过度优化
				break
			}
		}

		err := outputCallback()
		if err != nil {
			return err
		}

		if !moved {
			break
		}
	}

	s.R.log("JobMerge.Run ok totalScore=%f\n", MachinesGetScore(s.Machines))

	return nil
}
