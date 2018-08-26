package main

import (
	"fmt"
	"github.com/NeuronEvolution/aliyun_semi/schedule"
	"time"
)

const appResourcesFile = "./data/app_resources.csv"
const appInterferenceFile = "./data/app_interference.csv"

func run(appResourcesConfigMap []*schedule.AppResourcesConfig, appInferenceConfigMap [][]int, data string) {
	machineResourceConfigList, err := schedule.LoadMachineResourcesConfig("./data/machine_resources." + data + ".csv")
	if err != nil {
		fmt.Println(data, err)
		return
	}

	instanceDeployConfigList, err := schedule.LoadInstanceDeployConfig("./data/instance_deploy." + data + ".csv")
	if err != nil {
		fmt.Println(data, err)
		return
	}

	jobConfigMap, jobConfigDag, err := schedule.LoadJobDAG("./data/job_info." + data + ".csv")
	if err != nil {
		fmt.Println(data, err)
		return
	}

	r := schedule.NewResourceManagement(
		appResourcesConfigMap,
		appInferenceConfigMap,
		machineResourceConfigList,
		instanceDeployConfigList,
		jobConfigMap,
		jobConfigDag, data, "./_output/"+data)

	err = r.Run()
	if err != nil {
		fmt.Println(data, err)
	}
}

func main() {
	appResourceConfigMap, appInferenceConfigMap, err := schedule.LoadAppConfig(appResourcesFile, appInterferenceFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	//go run(appResourceConfigMap, appInferenceConfigMap, "a")
	//go run(appResourceConfigMap, appInferenceConfigMap, "b")
	//go run(appResourceConfigMap, appInferenceConfigMap, "c")
	go run(appResourceConfigMap, appInferenceConfigMap, "d")
	//go run(appResourceConfigMap, appInferenceConfigMap, "e")

	for {
		time.Sleep(time.Second * 30)
		output()
	}
}
