package semitemp

func (r *ResourceManagement) Init(
	share *ResourceManagement,
	machineResourcesConfig []*MachineResourcesConfig,
	appResourcesConfig []*AppResourcesConfig,
	appInterferenceConfig []*AppInterferenceConfig,
	instanceDeployConfig []*InstanceDeployConfig) (err error) {

	r.machineResourcesConfig = machineResourcesConfig
	r.appResourcesConfig = appResourcesConfig
	r.appInterferenceConfig = appInterferenceConfig

	r.Initializing = true
	defer func() { r.Initializing = false }()

	r.MachineConfigMap = make([]*MachineResourcesConfig, MaxMachineId)
	if share == nil {
		r.AppResourcesConfigMap = make([]*AppResourcesConfig, MaxAppId)
		r.AppInterferenceConfigMap = make([][MaxAppId]int, MaxAppId)
		for i := 0; i < MaxAppId; i++ {
			for j := 0; j < MaxAppId; j++ {
				r.AppInterferenceConfigMap[i][j] = -1
			}
		}

		if appResourcesConfig != nil {
			for _, v := range appResourcesConfig {
				err = r.SaveAppResourceConfig(v)
				if err != nil {
					return err
				}
			}
		}

		if appInterferenceConfig != nil {
			for _, v := range appInterferenceConfig {
				err = r.SaveAppInterferenceConfig(v)
				if err != nil {
					return err
				}
			}
		}
	} else {
		r.AppResourcesConfigMap = share.AppResourcesConfigMap
		r.AppInterferenceConfigMap = share.AppInterferenceConfigMap
	}

	if machineResourcesConfig != nil {
		for _, v := range machineResourcesConfig {
			err = r.AddMachine(v)
			if err != nil {
				return err
			}
		}
	}

	if instanceDeployConfig != nil {
		err = r.InitInstanceDeploy(instanceDeployConfig)
		if err != nil {
			return err
		}
	}

	return nil
}
