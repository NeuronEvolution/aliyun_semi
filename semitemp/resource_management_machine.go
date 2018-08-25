package semitemp

import "fmt"

//添加机器
//todo 将触发调度
//todo 异步将最低资源机器释放，相关instance重新拉入

func (r *ResourceManagement) AddMachine(config *MachineResourcesConfig) error {
	r.MachineConfigMap[config.MachineId] = config
	machine := NewMachine(r, config.MachineId, r.MachineLevelConfigPool.GetConfig(&config.MachineLevelConfig))
	r.MachineMap[machine.MachineId] = machine
	r.addFreeMachine(machine)

	return nil
}

//删除机器
//不能删除有实例部署的机器
//todo 异步将最低资源机器释放，相关instance重新拉入
func (r *ResourceManagement) RemoveMachine(machineId int) error {
	m := r.MachineMap[machineId]
	if m.InstanceArrayCount > 0 {
		return fmt.Errorf("R.RemoveMachine 机器%d已部署%d个实例",
			machineId, m.InstanceArrayCount)
	}

	r.MachineConfigMap[machineId] = nil
	r.MachineMap[machineId] = nil
	r.removeFromFreeMachine(machineId)

	return nil
}

func (r *ResourceManagement) addFreeMachine(m *Machine) {
	r.MachineFreePool.AddMachine(m)
}

func (r *ResourceManagement) removeFromFreeMachine(machineId int) *Machine {
	return r.MachineFreePool.RemoveMachine(machineId)
}
