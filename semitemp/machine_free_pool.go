package semitemp

import (
	"sort"
)

type MachineLevelFree struct {
	LevelConfig       *MachineLevelConfig
	MachineCollection *MachineCollection
}

func NewMachineLevelFree(level *MachineLevelConfig) *MachineLevelFree {
	p := &MachineLevelFree{}
	p.LevelConfig = level
	p.MachineCollection = NewMachineCollection()

	return p
}

func (p *MachineLevelFree) PushMachine(m *Machine) {
	//debugLog("MachineLevelFree.PushMachine %d %v", m.MachineId, m.Config)
	p.MachineCollection.Add(m)
}

func (p *MachineLevelFree) PeekMachine() (m *Machine) {
	//debugLog("MachineLevelFree.PeekMachine MachineListCount=%d %v",
	//p.MachineCollection.ListCount, p.Config)
	return p.MachineCollection.Peek()
}

func (p *MachineLevelFree) RemoveMachine(machineId int) {
	//debugLog("MachineLevelFree.RemoveMachine %d %v", machineId, p.Config)
	p.MachineCollection.Remove(machineId)
}

type MachineLevelFreeListSortByMachineLevelDesc []*MachineLevelFree

func (p MachineLevelFreeListSortByMachineLevelDesc) Len() int {
	return len(p)
}

func (p MachineLevelFreeListSortByMachineLevelDesc) Less(i, j int) bool {
	return !p[i].LevelConfig.Less(p[j].LevelConfig)
}

func (p MachineLevelFreeListSortByMachineLevelDesc) Swap(i, j int) {
	temp := p[i]
	p[i] = p[j]
	p[j] = temp
}

type MachineFreePool struct {
	MachineMap            map[int]*Machine
	MachineLevelFreeArray MachineLevelFreeListSortByMachineLevelDesc
}

func NewMachineFreePool() *MachineFreePool {
	p := &MachineFreePool{}
	p.MachineMap = make(map[int]*Machine)

	return p
}

func (p *MachineFreePool) AddMachine(m *Machine) {
	//debugLog("MachineFreePool.AddMachine %d", m.MachineId)
	p.MachineMap[m.MachineId] = m

	var pool *MachineLevelFree
	for _, v := range p.MachineLevelFreeArray {
		if v.LevelConfig == m.LevelConfig {
			pool = v
			break
		}
	}
	if pool == nil {
		//debugLog("MachineFreePool.AddMachine new level %v", m.Config)
		pool = NewMachineLevelFree(m.LevelConfig)
		p.MachineLevelFreeArray = append(p.MachineLevelFreeArray, pool)
		sort.Sort(p.MachineLevelFreeArray)
	}
	pool.PushMachine(m)
}

func (p *MachineFreePool) RemoveMachine(machineId int) *Machine {
	//debugLog("MachineFreePool.RemoveMachine %s", machineId)
	m, has := p.MachineMap[machineId]
	if !has {
		return nil
	}

	delete(p.MachineMap, machineId)

	for _, v := range p.MachineLevelFreeArray {
		if v.LevelConfig == m.LevelConfig {
			v.RemoveMachine(machineId)
			return m
		}
	}

	return m
}

func (p *MachineFreePool) PeekMachine() (m *Machine) {
	//debugLog("MachineFreePool.PeekMachine")
	for _, v := range p.MachineLevelFreeArray {
		if v.MachineCollection.ListCount > 0 {
			m = v.PeekMachine()
			if m != nil {
				//debugLog("MachineFreePool.PeekMachine success,machineId=%s remain=%d %v",
				//	m.MachineId, len(p.MachineMap), m.Config)
				return m
			}
		}
	}

	return nil
}

func (p *MachineFreePool) PeekMachineList(count int) []*Machine {
	machineList := make([]*Machine, 0)
	for _, v := range p.MachineLevelFreeArray {
		if v.MachineCollection.ListCount > 0 {
			if v.MachineCollection.ListCount >= count-len(machineList) {
				machineList = append(machineList, v.MachineCollection.List[:count-len(machineList)]...)
				break
			} else {
				machineList = append(machineList, v.MachineCollection.List[:v.MachineCollection.ListCount]...)
			}
		}
	}

	return machineList
}
