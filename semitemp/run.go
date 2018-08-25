package semitemp

import (
	"fmt"
)

func Run(machineResourcesList []*MachineResourcesConfig,
	appResourcesList []*AppResourcesConfig,
	appInterferenceList []*AppInterferenceConfig,
	instanceDeployList []*InstanceDeployConfig,
	strategyCreator func(r *ResourceManagement) Strategy) (result *ResourceManagement, err error) {

	instanceList := make([]*InstanceDeployConfig, 0)
	instanceMachineList := make([]*InstanceDeployConfig, 0)
	for _, v := range instanceDeployList {
		if v.MachineId == 0 {
			instanceList = append(instanceList, v)
		} else {
			instanceMachineList = append(instanceMachineList, v)
		}
	}

	fmt.Printf("---------------------------------------------DEPLOY-----------------------------------------\n")
	deploy := NewResourceManagement()
	deploy.tempInstanceDeployConfig = instanceMachineList
	deploy.SetStrategy(strategyCreator(deploy))
	err = deploy.Init(nil, machineResourcesList, appResourcesList, appInterferenceList, nil)
	if err != nil {
		return nil, err
	}
	deploy.DebugPrintStatus()

	err = deploy.AddInstanceList(instanceDeployList)
	if err != nil {
		return nil, err
	}
	deploy.DebugPrintStatus()

	fmt.Printf("---------------------------------------------MERGE-----------------------------------------\n")

	merge := NewResourceManagement()
	err = merge.Init(deploy, machineResourcesList, appResourcesList, appInterferenceList, instanceMachineList)
	if err != nil {
		return nil, err
	}
	merge.DebugPrintStatus()

	err = merge.MergeTo(deploy)
	if err != nil {
		return nil, err
	}
	merge.DebugPrintStatus()

	fmt.Printf("---------------------------------------------REPLAY-----------------------------------------\n")

	playback := NewResourceManagement()
	err = playback.Init(deploy, machineResourcesList, appResourcesList, appInterferenceList, instanceMachineList)
	if err != nil {
		return nil, err
	}
	playback.DebugPrintStatus()

	err = playback.Play(merge.DeployCommandHistory)
	if err != nil {
		return nil, err
	}
	playback.DebugPrintStatus()

	return merge, nil
}
