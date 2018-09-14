package main

import (
	"fmt"
)

func (r *ResourceManagement) firstFitInstances() (err error) {
	r.log("firstFitInstances start\n")

	instances := InstancesCopy(r.InstanceList)

	SortInstanceByTotalMaxLowWithInference(instances, 16)

	for i, instance := range instances {
		//if i > 0 && i%10000 == 0 {
		//	r.log("firstFitInstances %d\n", i)
		//}

		deployed := false
		for _, m := range r.MachineList {
			if m.ConstraintCheck(instance, 1) {
				m.AddInstance(instance)
				r.DeployMap[instance.InstanceId] = m
				//更新部署数量,todo 这里注释掉是因为手工指定了机器数量，之后放开
				//if machineIndex+1 > r.DeployedMachineCount {
				//	r.DeployedMachineCount = machineIndex + 1
				//}
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
