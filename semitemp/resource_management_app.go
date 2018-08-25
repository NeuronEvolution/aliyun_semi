package semitemp

import "fmt"

//保存应用资源需求配置
//todo 将触发调度
//todo 异步将该app相关instance拉出再重新拉入
func (r *ResourceManagement) SaveAppResourceConfig(config *AppResourcesConfig) error {
	config.calcCostEval(&MachineLevelConfig{
		Cpu:  MachineCpuMax,
		Mem:  MachineMemMax,
		Disk: MachineDiskMax,
		P:    MachinePMax,
		M:    MachineMMax,
		PM:   MachinePMMax,
	})
	r.AppResourcesConfigMap[config.AppId] = config
	return nil
}

//保存应用冲突配置
//todo 将触发调度
//todo 异步将该app相关instance拉出再重新拉入
func (r *ResourceManagement) SaveAppInterferenceConfig(config *AppInterferenceConfig) error {
	appResource1 := r.AppResourcesConfigMap[config.AppId1]
	if appResource1 == nil {
		return fmt.Errorf("SaveAppInterferenceConfig app %d not exists", config.AppId1)
	}

	appResource2 := r.AppResourcesConfigMap[config.AppId2]
	if appResource2 == nil {
		return fmt.Errorf("SaveAppInterferenceConfig app %d not esists", config.AppId2)
	}

	r.AppInterferenceConfigMap[config.AppId1][config.AppId2] = config.Interference

	return nil
}
