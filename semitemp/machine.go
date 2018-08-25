package semitemp

import (
	"fmt"
	"math"
	"math/rand"
)

type Machine struct {
	Resource
	R                  *ResourceManagement
	Rand               *rand.Rand
	MachineId          int
	LevelConfig        *MachineLevelConfig
	InstanceArray      []*Instance
	InstanceArrayCount int
	appCountCollection *AppCountCollection

	cpuCost      float64
	cpuCostValid bool
}

func NewMachine(r *ResourceManagement, machineId int, levelConfig *MachineLevelConfig) *Machine {
	m := &Machine{}
	m.R = r
	m.Rand = rand.New(rand.NewSource(0))
	m.MachineId = machineId
	m.LevelConfig = levelConfig
	m.InstanceArray = make([]*Instance, MaxInstancePerMachine)
	m.appCountCollection = NewAppCountCollection()

	return m
}

func (m *Machine) ClearInstances() {
	m.InstanceArrayCount = 0
	m.appCountCollection.Reset()
	for i := 0; i < len(m.Cpu); i++ {
		m.Cpu[i] = 0
	}
	for i := 0; i < len(m.Mem); i++ {
		m.Mem[i] = 0
	}
	m.Disk = 0
	m.P = 0
	m.M = 0
	m.PM = 0

	m.cpuCostValid = false
}

func (m *Machine) AddInstance(instance *Instance) {
	//debugLog("Machine.AddInstance %d %d", m.MachineId, instance.InstanceId)

	m.InstanceArray[m.InstanceArrayCount] = instance
	m.InstanceArrayCount++
	m.appCountCollection.Add(instance.Config.AppId)
	m.AddResource(&instance.Config.Resource)

	m.R.SetInstanceDeployedMachine(instance, m)
	if m.InstanceArrayCount == 1 {
		m.R.MachineFreePool.RemoveMachine(m.MachineId)
		m.R.MachineDeployPool.AddMachine(m)
	}

	m.cpuCostValid = false

	m.calcCostEval(m.LevelConfig)

	if DebugEnabled {
		//m.debugValidation()
	}
}

func (m *Machine) RemoveInstance(instanceId int) {
	//debugLog("Machine.RemoveInstance machineId=%d,instanceId=%d", m.MachineId, instanceId)
	for i, v := range m.InstanceArray[:m.InstanceArrayCount] {
		if v.InstanceId == instanceId {
			instance := m.InstanceArray[i]
			//debugLog("Machine.RemoveInstance appId=%d", instance.Config.AppId)
			m.InstanceArray[i] = nil
			if m.InstanceArrayCount > 1 && i < m.InstanceArrayCount-1 {
				for j := i; j < m.InstanceArrayCount-1; j++ {
					m.InstanceArray[j] = m.InstanceArray[j+1]
				}
				m.InstanceArray[m.InstanceArrayCount-1] = nil
			}

			m.InstanceArrayCount--
			m.appCountCollection.Remove(instance.Config.AppId)
			m.RemoveResource(&instance.Config.Resource)

			if m.InstanceArrayCount == 0 {
				m.R.MachineDeployPool.RemoveMachine(m.MachineId)
				m.R.MachineFreePool.AddMachine(m)
			}

			m.cpuCostValid = false

			break
		}
	}

	m.calcCostEval(m.LevelConfig)

	if DebugEnabled {
		//m.debugValidation()
	}
}

func (m *Machine) IsEmpty() bool {
	return len(m.InstanceArray) == 0
}

func (m *Machine) ConstraintCheck(instance *Instance, maxCpuRatio float64) bool {
	//debugLog("Machine.ConstraintCheck %s %s", m.MachineId, instance.InstanceId)

	if !ConstraintCheckResourceLimit(&m.Resource, &instance.Config.Resource, m.LevelConfig, maxCpuRatio) {
		//debugLog("Machine.ConstraintCheck constraintCheckResourceLimit failed")
		return false
	}

	if !ConstraintCheckAppInterferenceAddInstance(
		instance.Config.AppId,
		m.appCountCollection,
		m.R.AppInterferenceConfigMap) {
		//debugLog("Machine.ConstraintCheck constraintCheckAppInterferenceAddInstance failed")
		return false
	}

	return true
}

func (m *Machine) ConstraintCheckResourceLimit(instance *Instance, maxCpuRatio float64) bool {
	return ConstraintCheckResourceLimit(&m.Resource, &instance.Config.Resource, m.LevelConfig, maxCpuRatio)
}

func (m *Machine) ConstraintCheckAppInterferenceAddInstance(instance *Instance) bool {
	return ConstraintCheckAppInterferenceAddInstance(instance.Config.AppId,
		m.appCountCollection,
		m.R.AppInterferenceConfigMap)
}

func (m *Machine) HasBadConstraint() bool {
	return !ConstraintCheckAppInterference(m.appCountCollection, m.R.AppInterferenceConfigMap)
}

func (m *Machine) GetCpuCostReal() float64 {
	totalCost := float64(0)
	for i := 0; i < TimeSampleCount; i++ {
		r := m.Cpu[i] / m.LevelConfig.Cpu
		if r > 0.5 {
			totalCost += 1 + 10*(math.Exp(r-0.5)-1)
		} else {
			totalCost += 1
		}
	}

	return totalCost / TimeSampleCount
}

func (m *Machine) GetCpuCost() float64 {
	if m.cpuCostValid {
		return m.cpuCost
	}
	m.cpuCostValid = true

	totalCost := float64(0)
	for i := 0; i < TimeSampleCount; i++ {
		r := m.Cpu[i] / m.LevelConfig.Cpu
		if r > 0.5 {
			totalCost += 1 + 10*(Exp(r-0.5)-1)
		} else {
			totalCost += 1
		}
	}

	m.cpuCost = totalCost / TimeSampleCount

	return m.cpuCost
}

func (m *Machine) GetCostWithInstance(instance *Instance) float64 {
	totalCost := float64(0)
	for i := 0; i < TimeSampleCount; i++ {
		r := (m.Cpu[i] + instance.Config.Cpu[i]) / m.LevelConfig.Cpu
		if r > 0.5 {
			totalCost += 1 + 10*(Exp(r-0.5)-1)
		} else {
			totalCost += 1
		}
	}

	return totalCost / TimeSampleCount
}

func (m *Machine) GetDerivationWithInstance(instance *Instance) float64 {
	avg := float64(0)
	for i := 0; i < TimeSampleCount; i++ {
		r := (m.Cpu[i] + instance.Config.Cpu[i]) / m.LevelConfig.Cpu
		avg += r
	}

	avg = avg / float64(TimeSampleCount)
	d := float64(0)
	for i := 0; i < TimeSampleCount; i++ {
		r := (m.Cpu[i] + instance.Config.Cpu[i]) / m.LevelConfig.Cpu
		d += (r - avg) * (r - avg)
	}
	d = math.Sqrt(d / TimeSampleCount)

	return d
}

func (r *Resource) GetDerivationWithInstances(instances []*Instance) float64 {
	var cpu [TimeSampleCount]float64
	avg := float64(0)
	for i := 0; i < TimeSampleCount; i++ {
		total := r.Cpu[i]
		for _, instance := range instances {
			total += instance.Config.Cpu[i]
		}

		if total > HighCpu {
			return math.MaxFloat64
		}

		cpu[i] = total / HighCpu
		avg += cpu[i]
	}
	avg = avg / float64(TimeSampleCount)

	d := float64(0)
	for i := 0; i < TimeSampleCount; i++ {
		d += (cpu[i] - avg) * (cpu[i] - avg)
	}
	d = math.Sqrt(d / TimeSampleCount)

	return d / avg
}

func (m *Machine) GetLinearCostWithInstance(instance *Instance) float64 {
	totalCost := float64(0)
	for i := 0; i < TimeSampleCount; i++ {
		r := (m.Cpu[i] + instance.Config.Cpu[i]) / m.LevelConfig.Cpu
		if r > 0.5 {
			totalCost += 1 + 10*(Exp(r-0.5)-1)
		} else {
			totalCost += 2 * r
		}
	}
	totalCost = totalCost / TimeSampleCount

	return totalCost
}

func (m *Machine) CalcDeviationWithInstance(inst *Instance) float64 {
	cpuMax := float64(0)
	for i, v := range m.Cpu {
		cpu := v + inst.Config.Cpu[i]
		if cpu > cpuMax {
			cpuMax = cpu
		}
	}

	avgMem := float64(0)
	for i, v := range m.Mem {
		avgMem += v + inst.Config.Mem[i]
	}
	avgMem = avgMem / float64(len(m.Mem))

	cpu := cpuMax / (m.LevelConfig.Cpu * MaxCpuRatio)
	mem := avgMem / m.LevelConfig.Mem
	disk := float64(m.Disk+inst.Config.Disk) / float64(m.LevelConfig.Disk)
	p := float64(m.P+inst.Config.P) / float64(m.LevelConfig.P)
	mCost := float64(m.M+inst.Config.M) / float64(m.LevelConfig.M)
	pm := float64(m.PM+inst.Config.PM) / float64(m.LevelConfig.PM)

	//debugLog("calcResourceCostDeviation %f %f %f %f %f %f\n", cpu, mem, disk, p, mCost, pm)

	return calcResourceCostDeviation(cpu, mem, disk, p, mCost, pm)
}

func (m *Machine) GetResourceCostWithInstance(inst *Instance) float64 {
	avgCpu := float64(0)
	for i, v := range m.Cpu {
		avgCpu += v + inst.Config.Cpu[i]
	}
	avgCpu = avgCpu / float64(len(m.Cpu))

	avgMem := float64(0)
	for i, v := range m.Mem {
		avgMem += v + inst.Config.Mem[i]
	}
	avgMem = avgMem / float64(len(m.Mem))

	cpu := avgCpu * ParamCpuCostMultiply / m.LevelConfig.Cpu
	mem := avgMem / m.LevelConfig.Mem
	disk := float64(m.Disk+inst.Config.Disk) / float64(m.LevelConfig.Disk)
	p := float64(m.P+inst.Config.P) / float64(m.LevelConfig.P)
	mCost := float64(m.M+inst.Config.M) / float64(m.LevelConfig.M)
	pm := float64(m.PM+inst.Config.PM) / float64(m.LevelConfig.PM)

	cost := m.scaleCost(cpu) +
		m.scaleCost(mem) +
		m.scaleCost(disk) +
		m.scaleCost(p) +
		m.scaleCost(mCost) +
		m.scaleCost(pm)

	deviation := calcResourceCostDeviation(cpu, mem, disk, p, mCost, pm)
	deviation = 1 + deviation - m.ResourceCostDeviation
	appCount := m.appCountCollection.GetAppCount(inst.Config.AppId)

	if appCount > MaxAppCountInMachine {
		MaxAppCountInMachine = appCount
		fmt.Println("MaxAppCountInMachine", MaxAppCountInMachine)
	}

	return cost * Exp(ParamDeviationMultiply*deviation+ParamAppInferenceMultiply*float64(appCount))
}

var MaxAppCountInMachine = 0

func (m *Machine) scaleCost(r float64) float64 {
	return Exp(ParamMachineCostMultiply * r)
}

func (m *Machine) debugValidation() {
	for i := 0; i < m.InstanceArrayCount; i++ {
		if m.InstanceArray[i] == nil {
			panic(fmt.Errorf("Machine.debugValidation machineId=%d,i=%d", m.MachineId, i))
		}
	}

	m.appCountCollection.debugValidation()
}

func (m *Machine) DebugPrint() {
	fmt.Printf("Machine.DebugPrint %d %v cost=%f linearCost=%f\n",
		m.MachineId, m.LevelConfig, m.GetCpuCostReal(), m.GetLinearCpuCost(m.LevelConfig.Cpu))
	for i := 0; i < m.appCountCollection.ListCount; i++ {
		fmt.Printf("    %v,", m.appCountCollection.List[i])
		m.R.AppResourcesConfigMap[m.appCountCollection.List[i].AppId].DebugPrint()
	}

	m.Resource.DebugPrint()
}

func MachinesContains(machines []*Machine, machineId int) bool {
	for _, v := range machines {
		if v.MachineId == machineId {
			return true
		}
	}

	return false
}

func MachinesRemove(machines []*Machine, removes []*Machine) (rest []*Machine) {
	rest = make([]*Machine, 0)
	for _, v := range machines {
		has := false
		for _, i := range removes {
			if i.MachineId == v.MachineId {
				has = true
				break
			}
		}
		if !has {
			rest = append(rest, v)
		}
	}

	return rest
}
