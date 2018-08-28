package schedule

import (
	"math/rand"
	"sort"
)

type ResourceManagement struct {
	//Config
	AppResourcesConfigMap    []*AppResourcesConfig
	AppInterferenceConfigMap [][]int
	MachineConfigList        []*MachineResourcesConfig
	MachineConfigMap         []*MachineResourcesConfig
	InstanceDeployConfigList []*InstanceDeployConfig
	InstanceDeployConfigMap  []*InstanceDeployConfig
	JobConfigMap             []*JobConfig
	JobConfigDAG             []*JobConfig
	MachineConfigPool        *MachineConfigPool
	OutputDir                string
	Dataset                  string

	//Status
	Rand *rand.Rand

	MaxJobInstanceId     int
	MachineList          []*Machine
	MachineMap           []*Machine
	MaxMachineId         int
	DeployedMachineCount int
	InstanceList         []*Instance
	InstanceMap          []*Instance
	MaxInstanceId        int
	DeployMap            []*Machine
	JobList              []*Job
	JobMap               []*Job
	JobDeployMap         []*Machine
}

func NewResourceManagement(
	appResourcesConfigMap []*AppResourcesConfig,
	appInterferenceConfigMap [][]int,
	machineConfigList []*MachineResourcesConfig,
	instanceDeployConfigList []*InstanceDeployConfig,
	jobConfigMap []*JobConfig,
	jobConfigDAG []*JobConfig,
	dataset string,
	outputDir string) *ResourceManagement {

	r := &ResourceManagement{}
	r.AppResourcesConfigMap = appResourcesConfigMap
	r.AppInterferenceConfigMap = appInterferenceConfigMap
	r.MachineConfigList = machineConfigList
	r.InstanceDeployConfigList = instanceDeployConfigList
	r.JobConfigMap = jobConfigMap
	r.JobConfigDAG = jobConfigDAG
	r.Dataset = dataset
	r.OutputDir = outputDir
	r.Rand = rand.New(rand.NewSource(0))
	r.MachineConfigPool = NewMachineConfigPool()

	for _, config := range r.MachineConfigList {
		if config != nil && config.MachineId > r.MaxMachineId {
			r.MaxMachineId = config.MachineId
		}
	}
	r.MachineConfigMap = make([]*MachineResourcesConfig, r.MaxMachineId+1)
	for _, config := range r.MachineConfigList {
		r.MachineConfigMap[config.MachineId] = config
	}

	for _, config := range r.InstanceDeployConfigList {
		if config != nil && config.InstanceId > r.MaxInstanceId {
			r.MaxInstanceId = config.InstanceId
		}
	}
	r.InstanceDeployConfigMap = make([]*InstanceDeployConfig, r.MaxInstanceId+1)
	for _, config := range r.InstanceDeployConfigList {
		r.InstanceDeployConfigMap[config.InstanceId] = config
	}

	return r
}

func (r *ResourceManagement) createMachines() (machineList []*Machine, machineMap []*Machine) {
	machineMap = make([]*Machine, r.MaxMachineId+1)
	for _, config := range r.MachineConfigList {
		m := NewMachine(r, config.MachineId, r.MachineConfigPool.GetConfig(&config.MachineConfig))
		machineMap[m.MachineId] = m
		machineList = append(machineList, m)
	}

	sort.Slice(machineList, func(i, j int) bool {
		return machineList[i].Config.Cpu > machineList[j].Config.Cpu
	})

	return machineList, machineMap
}

func (r *ResourceManagement) createInstances() (instanceList []*Instance, instanceMap []*Instance) {
	instanceMap = make([]*Instance, r.MaxInstanceId+1)
	for _, config := range r.InstanceDeployConfigList {
		if config == nil {
			continue
		}

		instance := NewInstance(r, config.InstanceId, r.AppResourcesConfigMap[config.AppId])
		instanceMap[instance.InstanceId] = instance
		instanceList = append(instanceList, instance)
	}

	return instanceList, instanceMap
}

func (r *ResourceManagement) createJobs() {
	r.JobMap = make([]*Job, 0)
	r.JobMap = append(r.JobMap, nil)
	currentJobInstanceId := 0
	//打包创建实例
	for _, config := range r.JobConfigMap {
		if config == nil {
			continue
		}

		config.State = &JobCommonState{}
		config.State.UpdateTime()

		rest := config.InstanceCount
		packCount := config.getPackCount()
		for {
			count := packCount
			if rest < packCount {
				count = rest
			}

			currentJobInstanceId++
			job := NewJob(r, currentJobInstanceId, config, count)
			r.JobList = append(r.JobList, job)
			r.JobMap = append(r.JobMap, job)
			r.MaxJobInstanceId = currentJobInstanceId
			config.State.Jobs = append(config.State.Jobs, job)

			rest -= packCount
			if rest <= 0 {
				break
			}
		}
	}

	r.log("createJobs MaxJobInstanceId=%d\n", r.MaxJobInstanceId)
}

func (r *ResourceManagement) init() (err error) {
	err = MakeDirIfNotExists(r.OutputDir + "/")
	if err != nil {
		return err
	}

	//创建机器
	r.MachineList, r.MachineMap = r.createMachines()

	//创建实例
	r.InstanceList, r.InstanceMap = r.createInstances()

	//创建任务
	r.initJobConfigs()
	r.createJobs()

	//数据简单分析
	r.analysis()

	return nil
}

func (r *ResourceManagement) CalcTotalScore() float64 {
	score := float64(0)
	for _, m := range r.MachineList[:r.DeployedMachineCount] {
		score += m.GetCpuCost()
	}

	return score
}

func (r *ResourceManagement) CalcTotalScoreReal() float64 {
	score := float64(0)
	for _, m := range r.MachineList[:r.DeployedMachineCount] {
		score += m.GetCpuCostReal()
	}

	return score
}

func (r *ResourceManagement) initE() (err error) {
	r.DeployedMachineCount = 8000
	for _, config := range r.InstanceDeployConfigList {
		m := r.MachineMap[config.MachineId]
		m.AddInstance(r.InstanceMap[config.InstanceId])
		r.DeployMap[config.InstanceId] = m
	}

	return nil
}

func (r *ResourceManagement) beginOffline() {
	//将计算点从实例的98点提升到98*15点
	for _, m := range r.MachineList {
		m.beginOffline()
	}
}

func (r *ResourceManagement) Run() (err error) {
	err = r.init()
	if err != nil {
		return err
	}

	r.DeployMap = make([]*Machine, r.MaxInstanceId+1)
	r.JobDeployMap = make([]*Machine, r.MaxJobInstanceId+1)

	//初始化部署实例
	if r.Dataset == "e" {
		err = r.initE()
	} else {
		err = r.firstFitInstances()
	}
	if err != nil {
		return err
	}

	//优化迭代
	r.scheduleLoop()

	return nil
}
