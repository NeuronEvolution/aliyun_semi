package schedule

import (
	"sync"
)

func (r *ResourceManagement) scheduleTwoMachine(machines []*Machine, deadLoop int) (ok bool) {
	instances := make([]*Instance, 0)
	for _, m := range machines {
		instances = append(instances, m.InstanceList[:m.InstanceListCount]...)
	}

	cost := float64(0)
	for _, m := range machines {
		cost += m.GetCpuCost()
	}

	bestPos, bestCost := r.InstanceDeployForceBest(machines, instances, deadLoop)
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

func (r *ResourceManagement) scheduleMachines(machines []*Machine, deadLoop int) (has bool) {
	//两两分组并行调度
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

//从指定的机器中一大一小地随机一批机器
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

//计算实例部署最佳机器自动增长数量
func (r *ResourceManagement) instanceDeployCheckMachinesScale() (machineCountAllocate int) {
	h1 := 0
	h2 := 0
	h3 := 0
	for _, m := range r.MachineList[:r.DeployedMachineCount] {
		if m.GetCpuCost() > ScaleLimitH1 {
			h1++
		}
		if m.GetCpuCost() > ScaleLimitH2 {
			h2++
		}
		if m.GetCpuCost() > ScaleLimitH3 {
			h3++
		}
	}

	count := (h1-h2)/ScaleRatioH1 + (h2-h3)/ScaleRatioH2 + h3/ScaleRatioH3

	r.log("instanceDeployCheckMachinesScale h1=%4d,h2=%4d,h3=%4d,count=%4d\n", h1, h2, h3, count)

	return count
}

func (r *ResourceManagement) instanceSchedule() (err error) {
	startCost := MachinesGetScore(r.MachineList)
	totalLoop := 0
	for scaleCount := 0; ; scaleCount++ {
		currentCost := MachinesGetScore(r.MachineList)
		r.log("instanceSchedule scale=%2d start cost=%f\n", scaleCount, currentCost)
		pTableBigSmall := randBigSmall(r.DeployedMachineCount)
		pTableSmallBig := randSmallBig(r.DeployedMachineCount)
		loop := 0
		deadLoop := 0
		for ; ; loop++ {
			if r.Dataset == "e" {
				if totalLoop > 0 && totalLoop%128 == 0 && currentCost < 8460 {
					//todo fix
					return nil
				}
			} else if (r.Dataset == "a" && totalLoop > MachineALoop) ||
				(r.Dataset == "b" && totalLoop > MachineBLoop) ||
				(r.Dataset == "c" && totalLoop > MachineCLoop) ||
				(r.Dataset == "d" && totalLoop > MachineDLoop) {
				return nil
			}
			totalLoop++

			SortMachineByCpuCost(r.MachineList[:r.DeployedMachineCount])
			machinesByCpu := r.randomMachines(r.MachineList[:r.DeployedMachineCount], 32, pTableBigSmall, pTableSmallBig)
			ok := r.scheduleMachines(machinesByCpu, deadLoop)
			if !ok {
				r.log("instanceSchedule scale=%2d dead loop=%8d,totalLoop=%8d\n", scaleCount, deadLoop, totalLoop)
				deadLoop++
				continue
			}
			deadLoop = 0
			r.log("instanceSchedule scale=%2d loop=%8d,totalLoop=%8d %d %f %f\n",
				scaleCount, loop, totalLoop, r.DeployedMachineCount, startCost, MachinesGetScore(r.MachineList))
		}
	}

	return nil
}
