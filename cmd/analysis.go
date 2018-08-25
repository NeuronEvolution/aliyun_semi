package main

import (
	"github.com/NeuronEvolution/aliyun_x/cloud"
)

type AnalysisContext struct {
	appInterferenceList  []*cloud.AppInterferenceConfig
	appResourcesList     []*cloud.AppResourcesConfig
	machineResourcesList []*cloud.MachineResourcesConfig
	instanceDeployList   []*cloud.InstanceDeployConfig

	appResourcesMap        map[int]*cloud.AppResourcesConfig
	machineResourcesMap    map[int]*cloud.MachineResourcesConfig
	instanceMap            map[int]*cloud.InstanceDeployConfig
	instanceDeployedMap    map[int]*cloud.InstanceDeployConfig
	instanceNonDeployedMap map[int]*cloud.InstanceDeployConfig
}

func NewAnalysisContext(
	appInterferenceList []*cloud.AppInterferenceConfig,
	appResourcesList []*cloud.AppResourcesConfig,
	machineResourcesList []*cloud.MachineResourcesConfig,
	instanceDeployList []*cloud.InstanceDeployConfig) *AnalysisContext {
	c := &AnalysisContext{}
	c.appInterferenceList = appInterferenceList
	c.appResourcesList = appResourcesList
	c.machineResourcesList = machineResourcesList
	c.instanceDeployList = instanceDeployList

	return c
}

func (c *AnalysisContext) init() {
	c.appResourcesMap = make(map[int]*cloud.AppResourcesConfig)
	for _, v := range c.appResourcesList {
		c.appResourcesMap[v.AppId] = v
	}

	c.machineResourcesMap = make(map[int]*cloud.MachineResourcesConfig)
	for _, v := range c.machineResourcesList {
		c.machineResourcesMap[v.MachineId] = v
	}

	c.instanceMap = make(map[int]*cloud.InstanceDeployConfig)
	c.instanceDeployedMap = make(map[int]*cloud.InstanceDeployConfig)
	c.instanceNonDeployedMap = make(map[int]*cloud.InstanceDeployConfig)
	for _, v := range c.instanceDeployList {
		c.instanceMap[v.InstanceId] = v
		if v.MachineId == 0 {
			c.instanceNonDeployedMap[v.InstanceId] = v
		} else {
			c.instanceDeployedMap[v.InstanceId] = v
		}
	}
}

func (c *AnalysisContext) Run() {
	c.init()

	c.AnalysisArraySize()
	c.AnalysisDiskDistribution()
	//c.AnalysisAppInference()
}
