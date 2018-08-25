package schedule

import (
	"fmt"
	"math"
	"os"
	"sync"
	"time"
)

func (r *ResourceManagement) scheduleMachines(machines []*Machine, deadLoop int) (has bool) {
	wg := &sync.WaitGroup{}
	max := len(machines)
	if len(machines)%2 == 1 {
		max = len(machines) - 1
	}
	for i := 0; i < max; i += 2 {
		batchMachines := []*Machine{machines[i], machines[i+1]}
		wg.Add(1)
		go func() {
			defer wg.Done()

			ok := r.scheduleTwoMachine(batchMachines, deadLoop)
			if ok {
				has = true
			}
		}()
	}

	wg.Wait()

	return has
}

func (r *ResourceManagement) scheduleTwoMachine(machines []*Machine, deadLoop int) (ok bool) {
	instances := make([]*Instance, 0)
	for _, m := range machines {
		instances = append(instances, m.InstanceList[:m.InstanceListCount]...)
	}

	cost := float64(0)
	for _, m := range machines {
		cost += m.GetCpuCost()
	}

	bestPos, bestCost := r.Best(machines, instances, deadLoop)
	if bestCost >= cost {
		return false
	}

	//将所有实例迁出
	for _, m := range machines {
		for _, inst := range InstancesCopy(m.InstanceList[:m.InstanceListCount]) {
			m.RemoveInstance(inst.InstanceId)
		}
	}

	for i, instance := range instances {
		m := machines[bestPos[i]]
		if !m.ConstraintCheck(instance, m.Config.Cpu) {
			panic("ConstraintCheck")
		}
		m.AddInstance(instance)
		r.DeployMap[instance.InstanceId] = m
	}

	return true
}

func (r *ResourceManagement) randomMachines(pool []*Machine, count int, bigSmall []float64, smallBig []float64) (machines []*Machine) {
	machines = make([]*Machine, 0)
	for i := 0; i < count; i++ {
		var table []float64
		if len(machines)%2 == 0 {
			table = bigSmall
		} else {
			table = smallBig
		}

		maxP := table[len(table)-1]
		r := r.Rand.Float64() * maxP
		for machineIndex, p := range table {
			if p < r {
				continue
			}

			if MachinesContains(machines, pool[machineIndex].MachineId) {
				if machineIndex == count-1 {
					machineIndex = -1
				}
				continue
			}

			machines = append(machines, pool[machineIndex])
			break
		}
	}

	return machines
}

func (r *ResourceManagement) checkScale() (machineCountAllocate int) {
	h1 := 0
	h2 := 0
	h3 := 0
	for _, m := range r.MachineList[:r.DeployedMachineCount] {
		if m.GetCpuCost() > 1.01 {
			h1++
		}
		if m.GetCpuCost() > 1.1 {
			h2++
		}
		if m.GetCpuCost() > 1.2 {
			h3++
		}
	}

	return (h1-h2)/400 + (h2-h3)/40 + h3/20
}

func (r *ResourceManagement) scheduleLoop() {
	if r.Dataset == "e" {
		r.DeployedMachineCount = 8000
	}

	totalLoop := 0

	for scaleCount := 0; ; scaleCount++ {
		startCost := r.CalcTotalScore()
		r.log("scheduleLoop scale=%2d start cost=%f\n", scaleCount, startCost)

		lastSaveTime := time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local)

		pTableBigSmall := randBigSmall(r.DeployedMachineCount)
		pTableSmallBig := randSmallBig(r.DeployedMachineCount)

		currentCost := startCost
		loop := 0
		deadLoop := 0
		stop := false
		for ; ; loop++ {
			totalLoop++
			if totalLoop%256 == 0 {
				r.onlineMerge()
			}

			SortMachineByCpuCost(r.MachineList[:r.DeployedMachineCount])
			machinesByCpu := r.randomMachines(r.MachineList[:r.DeployedMachineCount], 32, pTableBigSmall, pTableSmallBig)
			ok := r.scheduleMachines(machinesByCpu, deadLoop)
			if !ok {
				r.log("scheduleLoop scale=%2d dead loop=%8d\n", scaleCount, deadLoop)
				deadLoop++
				continue
			}

			deadLoop = 0
			currentCost = r.CalcTotalScore()
			r.log("scheduleLoop scale=%2d loop=%8d %d %f %f\n",
				scaleCount, loop, r.DeployedMachineCount, startCost, currentCost)

			now := time.Now()
			if now.Sub(lastSaveTime).Seconds() > 60*30 {
				r.log("scheduleLoop scale=%2d save loop=%8d %d %f %f\n",
					scaleCount, loop, r.DeployedMachineCount, startCost, currentCost)
				err := r.mergeOutput()
				if err != nil {
					fmt.Printf("[%s]scheduleLoop save failed,%s\n",
						time.Now().Format(time.RFC3339), err)
				}
				lastSaveTime = now
			}

			_, err := os.Stat(r.OutputDir + "/" + "aliyun_stop")
			if err == nil {
				r.log("scheduleLoop scale=%2d loop=%8d aliyun_stop\n", scaleCount, loop)
				stop = true
				break
			}

			if loop == 128*int(math.Pow(1.414, float64(scaleCount))) {
				machineCountAllocate := r.checkScale()
				if machineCountAllocate > 0 {
					r.log("scheduleLoop scale=%2d machineCountAllocate %d\n", scaleCount, machineCountAllocate)
					r.DeployedMachineCount += machineCountAllocate
					if r.DeployedMachineCount > len(r.MachineList) {
						r.DeployedMachineCount = len(r.MachineList)
					}
					break
				}
			}
		}
		if stop {
			break
		}
	}

	r.log("scheduleLoop end  %d %f\n", r.DeployedMachineCount, r.CalcTotalScore())
}
