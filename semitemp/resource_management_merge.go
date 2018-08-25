package semitemp

import (
	"fmt"
)

func (r *ResourceManagement) MergeTo(status *ResourceManagement) (err error) {
	n := 0
	for _, v := range status.InstanceList {
		if v != nil {
			n++
		}
	}
	fmt.Println("MergeTo totalInstance", n)

	for _, m := range r.MachineMap {
		if m == nil {
			continue
		}

		//fmt.Println("marge machine", i)
		err = r.mergeMachine(m, status)
		if err != nil {
			return err
		}
	}

	notMerged := 0
	for _, m := range r.MachineMap {
		if m == nil {
			continue
		}

		for _, instance := range m.InstanceArray[:m.InstanceArrayCount] {
			if status.InstanceDeployedMachineMap[instance.InstanceId].MachineId != m.MachineId {
				notMerged++
			}
		}
	}
	if notMerged != 0 {
		return fmt.Errorf("not merged %d", notMerged)
	}

	return r.mergeNonDeployedMachines(status)
}

func (r *ResourceManagement) mergeMachine(m *Machine, status *ResourceManagement) (err error) {
	//迁出不匹配的
	instances := InstancesCopy(m.InstanceArray[:m.InstanceArrayCount])
	for _, instance := range instances {
		if status.InstanceDeployedMachineMap[instance.InstanceId].MachineId == m.MachineId {
			continue
		}

		mergedOut := false
		for _, targetMachine := range r.MachineMap {
			if targetMachine == nil || targetMachine.MachineId == m.MachineId {
				continue
			}

			if targetMachine.ConstraintCheck(instance, 1) {
				m.RemoveInstance(instance.InstanceId)
				r.CommandDeployInstance(instance, targetMachine)
				mergedOut = true
				break
			}
		}
		if !mergedOut {
			return fmt.Errorf("mergedOut failed %d %d", m.MachineId, instance.InstanceId)
		}
	}

	//迁入匹配的
	statusMachine := status.MachineMap[m.MachineId]
	for _, statusInstance := range statusMachine.InstanceArray[:statusMachine.InstanceArrayCount] {
		instance := r.InstanceList[statusInstance.InstanceId]
		if instance == nil {
			continue
		}

		machine := r.InstanceDeployedMachineMap[instance.InstanceId]
		if machine.MachineId == m.MachineId {
			continue
		}

		machine.RemoveInstance(instance.InstanceId)
		r.CommandDeployInstance(instance, m)
	}

	return nil
}

func (r *ResourceManagement) mergeNonDeployedMachines(status *ResourceManagement) (err error) {
	for _, instance := range status.InstanceList {
		if instance == nil {
			continue
		}

		if r.InstanceList[instance.InstanceId] == nil {
			newInstance := r.CreateInstance(instance.InstanceId, instance.Config)
			m := r.MachineMap[status.InstanceDeployedMachineMap[instance.InstanceId].MachineId]
			r.CommandDeployInstance(newInstance, m)
		}
	}

	return nil
}
