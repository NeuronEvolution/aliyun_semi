package semitemp

import (
	"bytes"
	"fmt"
	"sort"
)

type MachineLevelDeploy struct {
	LevelConfig       *MachineLevelConfig
	MachineCollection *MachineCollection
}

func NewMachineLevelDeploy(level *MachineLevelConfig) *MachineLevelDeploy {
	p := &MachineLevelDeploy{}
	p.LevelConfig = level
	p.MachineCollection = NewMachineCollection()

	return p
}

func (p *MachineLevelDeploy) AddMachine(m *Machine) {
	//debugLog("MachineLevelDeploy.AddMachine %s %v", m.MachineId, p.Config)
	p.MachineCollection.Add(m)
}

func (p *MachineLevelDeploy) RemoveMachine(machineId int) {
	//debugLog("MachineLevelDeploy.RemoveMachine %s %v", machineId, p.Config)
	p.MachineCollection.Remove(machineId)
}

type MachineLevelDeployListSortByMachineLevelDesc []*MachineLevelDeploy

func (p MachineLevelDeployListSortByMachineLevelDesc) Len() int {
	return len(p)
}

func (p MachineLevelDeployListSortByMachineLevelDesc) Less(i, j int) bool {
	return !p[i].LevelConfig.Less(p[j].LevelConfig)
}

func (p MachineLevelDeployListSortByMachineLevelDesc) Swap(i, j int) {
	temp := p[i]
	p[i] = p[j]
	p[j] = temp
}

func (p MachineLevelDeployListSortByMachineLevelDesc) First() *MachineLevelDeploy {
	if len(p) == 0 {
		return nil
	}

	return p[0]
}

func (p MachineLevelDeployListSortByMachineLevelDesc) Last() *MachineLevelDeploy {
	if len(p) == 0 {
		return nil
	}

	return p[len(p)-1]
}

type MachineDeployPool struct {
	R                       *ResourceManagement
	MachineLevelDeployArray MachineLevelDeployListSortByMachineLevelDesc
}

func NewMachineDeployPool(r *ResourceManagement) *MachineDeployPool {
	p := &MachineDeployPool{R: r}

	return p
}

func (p *MachineDeployPool) AddMachine(m *Machine) {
	//debugLog("MachineDeployPool.AddMachine %d", m.MachineId)

	var pool *MachineLevelDeploy
	for _, v := range p.MachineLevelDeployArray {
		if v.LevelConfig == m.LevelConfig {
			pool = v
			break
		}
	}
	if pool == nil {
		//debugLog("MachineDeployPool.AddMachine new level %v", m.Config)
		pool = NewMachineLevelDeploy(m.LevelConfig)
		p.MachineLevelDeployArray = append(p.MachineLevelDeployArray, pool)
		sort.Sort(p.MachineLevelDeployArray)
	}
	pool.AddMachine(m)
}

func (p *MachineDeployPool) RemoveMachine(machineId int) *Machine {
	//debugLog("MachineDeployPool.RemoveMachine %d", machineId)
	m := p.R.MachineMap[machineId]
	for _, v := range p.MachineLevelDeployArray {
		if v.LevelConfig == m.LevelConfig {
			v.RemoveMachine(machineId)
			return m
		}
	}

	return m
}

func (p *MachineDeployPool) DebugPrint(buf *bytes.Buffer) {
	buf.WriteString("MachineDeployPool.DebugPrint\n")
	instanceCount := 0
	machineCount := 0
	machineDeployed := make([]*Machine, 0)
	cpuHighCount := 0
	highCpuLimit := 1.01
	for _, m := range p.R.MachineMap {
		//v.DebugPrint()
		if m != nil && m.InstanceArrayCount > 0 {
			instanceCount += m.InstanceArrayCount
			machineCount++
			machineDeployed = append(machineDeployed, m)

			if m.GetCpuCost() > highCpuLimit {
				cpuHighCount++
			}
		}
	}
	for _, v := range p.MachineLevelDeployArray {
		buf.WriteString(fmt.Sprintf("    %v machineCount=%d\n",
			v.LevelConfig, v.MachineCollection.ListCount))
	}

	sort.Slice(machineDeployed, func(i, j int) bool {
		m1 := machineDeployed[i]
		m2 := machineDeployed[j]
		cpu1 := m1.GetCpuCostReal()
		cpu2 := m2.GetCpuCostReal()
		linearCpu1 := m1.GetLinearCpuCost(m1.LevelConfig.Cpu)
		linearCpu2 := m2.GetLinearCpuCost(m2.LevelConfig.Cpu)
		if cpu1 > 1.01 || cpu2 > 1.01 {
			return cpu1 > cpu2
		}

		return linearCpu1 > linearCpu2
	})

	//for i, m := range machineDeployed {
	//	if i < 32 {
	//buf.WriteString(fmt.Sprintf("    cpuCost=%f,machineId=%d\n", m.GetCpuCost(), m.MachineId))
	//m.DebugPrint()
	//	}

	//	if i >= len(machineDeployed)-32 {
	//m.DebugPrint()
	//	}
	//}
	buf.WriteString(fmt.Sprintf("total high cpu(%f) count=%d\n", highCpuLimit, cpuHighCount))

	buf.WriteString(fmt.Sprintf("MachineDeployPool.DebugPrint machineCount=%d,instanceCount=%d\n",
		machineCount, instanceCount))
}
