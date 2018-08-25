package main

import (
	"fmt"
	"github.com/NeuronEvolution/aliyun_x/cloud"
)

func (c *AnalysisContext) AnalysisArraySize() {
	var maxMachineId, maxAppId, maxInstanceId int
	for _, v := range c.machineResourcesList {
		if v.MachineId > maxMachineId {
			maxMachineId = v.MachineId
		}
	}
	for _, v := range c.appResourcesList {
		if v.AppId > maxAppId {
			maxAppId = v.AppId
		}
	}
	for _, v := range c.appInterferenceList {
		if v.AppId1 > maxAppId {
			maxAppId = v.AppId1
		}
		if v.AppId2 > maxAppId {
			maxAppId = v.AppId2
		}
	}

	instanceList := make([]*cloud.InstanceDeployConfig, 0)
	instanceMachineList := make([]*cloud.InstanceDeployConfig, 0)
	for _, v := range c.instanceDeployList {
		if v.MachineId == 0 {
			instanceList = append(instanceList, v)
		} else {
			instanceMachineList = append(instanceMachineList, v)
		}
		if v.AppId > maxAppId {
			maxAppId = v.AppId
		}
		if v.MachineId > maxMachineId {
			maxMachineId = v.MachineId
		}
		if v.InstanceId > maxInstanceId {
			maxInstanceId = v.InstanceId
		}
	}

	var appCount [cloud.MaxAppId]int
	for _, v := range c.instanceDeployList {
		appCount[v.AppId]++
	}
	maxAppCount := 0
	maxAppCountAppId := 0
	for i, v := range appCount {
		if v > maxAppCount {
			maxAppCount = v
			maxAppCountAppId = i
		}
	}

	fmt.Printf("AnalysisArraySize\n")
	fmt.Printf("    maxAppId=                    %d\n", maxAppId)
	fmt.Printf("    maxInstanceId=               %d\n", maxInstanceId)
	fmt.Printf("    maxMachineId=                %d\n", maxMachineId)
	fmt.Printf("    maxAppCount=                 %d,appId=%d\n", maxAppCount, maxAppCountAppId)
	fmt.Printf("    deployedInstanceCount=       %d\n", len(instanceMachineList))
	fmt.Printf("    non-deployedInstanceCount=   %d\n", len(instanceList))
	fmt.Printf("    totalInstantCount=           %d\n", len(c.instanceDeployList))
}
