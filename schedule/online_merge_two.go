package schedule

import (
	"math"
	"math/rand"
)

func (r *ResourceManagement) RandomBest(machines []*Machine, instances []*Instance) (bestPos []int, bestCost float64) {
	random := rand.New(rand.NewSource(int64(len(instances))))

	machineCount := len(machines)
	instanceCount := len(instances)
	pos := make([]int, instanceCount)
	bestPos = make([]int, instanceCount)
	bestCost = math.MaxFloat64

	for totalLoop := 0; totalLoop < 1000; totalLoop++ {
		deploy := make([]*Machine, len(machines))
		for i := 0; i < len(deploy); i++ {
			deploy[i] = NewMachine(r, machines[i].MachineId, machines[i].Config)
		}
		failed := false
		for instanceIndex := 0; instanceIndex < instanceCount; instanceIndex++ {
			instance := instances[instanceIndex]
			machineIndex := random.Intn(machineCount)
			m := deploy[machineIndex]
			if !m.ConstraintCheck(instance, 1) {
				failed = true
				continue
			}
			m.AddInstance(instance)
			pos[instanceIndex] = machineIndex
		}
		if !failed {
			totalCost := float64(0)
			for _, m := range deploy {
				totalCost += m.GetCpuCost()
			}

			//fmt.Println("BEST", totalCost)

			//最优解
			if totalCost < bestCost {
				//fmt.Println("BEST", bestCost, totalCost)
				bestCost = totalCost
				for i, v := range pos {
					bestPos[i] = v
				}
				//fmt.Println(bestPos)
			}
		}
	}

	return bestPos, bestCost
}

func (r *ResourceManagement) Best(machines []*Machine, instances []*Instance, deadLoop int) (bestPos []int, bestCost float64) {
	if deadLoop == 0 && r.Dataset != "e" {
		//return r.RandomBest(machines, instances)
	}

	e := deadLoop
	if e > 8 {
		e = 8
	}
	totalLoopLimit := 1024 * 8 * int(math.Pow(float64(2), float64(e)))

	machineCount := len(machines)
	instanceCount := len(instances)
	pos := make([]int, instanceCount)
	bestPos = make([]int, instanceCount)
	bestCost = math.MaxFloat64
	deploy := make([]*Machine, len(machines))
	for i := 0; i < len(deploy); i++ {
		deploy[i] = NewMachine(r, machines[i].MachineId, machines[i].Config)
	}

	totalLoop := 0

	for instanceIndex := 0; instanceIndex < instanceCount; instanceIndex++ {
		instance := instances[instanceIndex]
		added := false

		for ; pos[instanceIndex] < machineCount; pos[instanceIndex]++ {
			totalLoop++
			machineIndex := pos[instanceIndex]
			m := deploy[machineIndex]
			if !m.ConstraintCheck(instance, 1) {
				continue
			}
			m.AddInstance(instance)
			added = true
			break
		}

		if added {
			//有效解,回退
			if instanceIndex == instanceCount-1 {
				totalCost := float64(0)
				for _, m := range deploy {
					totalCost += m.GetCpuCost()
				}

				//fmt.Println("BEST", bestCost, totalCost, pos)

				//最优解
				if totalCost < bestCost {
					//fmt.Println("BEST", bestCost, totalCost)
					bestCost = totalCost
					for i, v := range pos {
						bestPos[i] = v
					}
					//fmt.Println(bestPos)
				}

				//回退
				deploy[pos[instanceIndex]].RemoveInstance(instance.InstanceId)
				pos[instanceIndex] = 0
			}
		} else {
			//回退
			pos[instanceIndex] = 0
		}

		end := false
		if !added || instanceIndex == instanceCount-1 {
			for {
				//已到最后
				instanceIndex--
				if instanceIndex < 0 {
					end = true
					break
				}

				deploy[pos[instanceIndex]].RemoveInstance(instances[instanceIndex].InstanceId)
				pos[instanceIndex]++
				if pos[instanceIndex] < machineCount {
					//进位成功
					instanceIndex--
					break
				} else {
					pos[instanceIndex] = 0
				}
			}
		}

		if end || (instanceCount > 20 && totalLoop > totalLoopLimit) {
			break
		}
	}

	return bestPos, bestCost
}
