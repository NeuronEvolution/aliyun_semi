package main

import (
	"fmt"
)

func (r *ResourceManagement) firstFitInstances() (err error) {
	r.log("firstFitInstances start\n")

	instances := make([]*Instance, 0)
	for _, v := range r.InstanceMap {
		if v != nil {
			instances = append(instances, v)
		}
	}

	SortInstanceByTotalMaxLowWithInference(instances, 16)

	for i, instance := range instances {
		if i > 0 && i%10000 == 0 {
			r.log("firstFitInstances %d\n", i)
		}

		deployed := false
		for machineIndex, m := range r.MachineList {
			if m.ConstraintCheck(instance, 1) {
				m.AddInstance(instance)
				r.DeployMap[instance.InstanceId] = m
				//更新部署数量
				if machineIndex+1 > r.DeployedMachineCount {
					r.DeployedMachineCount = machineIndex + 1
				}
				deployed = true
				break
			}
		}
		if !deployed {
			return fmt.Errorf(fmt.Sprintf("firstFitInstances failed,%d instanceId=%d,appId=%d",
				i, instance.InstanceId, instance.Config.AppId))
		}
	}

	r.log("firstFitInstances deployedMachineCount=%d,score=%f\n",
		r.DeployedMachineCount, MachinesGetScore(r.MachineList))

	return nil
}
