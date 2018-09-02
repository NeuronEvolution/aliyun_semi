package schedule

import (
	"math/rand"
	"sort"
	"time"
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
	Rand      *rand.Rand
	StartTime time.Time

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
	r.StartTime = time.Now()
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

	totalJobCount := 0
	for _, config := range r.JobConfigMap {
		if config == nil {
			continue
		}

		totalJobCount += config.InstanceCount
	}

	//打包创建实例
	currentJobInstanceId := 0
	for _, config := range r.JobConfigMap {
		if config == nil {
			continue
		}

		rest := config.InstanceCount
		packCount := config.getPackCount(totalJobCount)
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
	//todo 为节约时间，这里不再自动探测最佳机器数量
	if r.Dataset == "a" {
		r.DeployedMachineCount = MachineA
	} else if r.Dataset == "b" {
		r.DeployedMachineCount = MachineB
	} else if r.Dataset == "c" {
		r.DeployedMachineCount = MachineC
	} else if r.Dataset == "d" {
		r.DeployedMachineCount = MachineD
	} else if r.Dataset == "e" {
		r.DeployedMachineCount = MachineE
	}

	//初始化
	err = r.init()
	if err != nil {
		return err
	}

	r.DeployMap = make([]*Machine, r.MaxInstanceId+1)
	r.JobDeployMap = make([]*Machine, r.MaxJobInstanceId+1)

	//return NewJobScheduler(r, r.MachineList).Run()

	//加载预先计算的实例部署，节约时间
	instanceMoveCommands, err := r.loadInstanceMoveCommands()
	if err != nil {
		//初始化部署实例
		if r.Dataset == "e" {
			err = r.initE()
		} else {
			err = r.firstFitInstances()
		}
		if err != nil {
			return err
		}

		//实例调度
		err = r.instanceSchedule()
		if err != nil {
			return err
		}

		//之后实例不再调度，先计算出实例迁移指令
		instanceMoveCommands, err := NewOnlineMerge(r).Run()
		if err != nil {
			return err
		}

		r.saveInstanceMoveCommands(instanceMoveCommands)
	}

	//todo 这里需要考虑在线迁移时的实例交换,改为从初始状态迁移后再部署任务,暂时不需要优化，除了e数据不需要固定实例
	//重新插入实例，避免浮点精度问题
	machines := MachinesCloneWithInstances(r.MachineList)
	r.log("jobSchedule init totalScore=%f\n", MachinesGetScore(machines))

	//任务调度
	//err = r.firstFitJobs(machines)
	err = NewJobScheduler(r, machines[:r.DeployedMachineCount]).Run()
	if err != nil {
		return err
	}
	r.log("jobSchedule totalScore=%f\n", MachinesGetScore(machines))

	//构造任务调度指令
	jobDeployCommands := r.buildJobDeployCommands(machines)

	//验证结果
	err = NewReplay(r, instanceMoveCommands, jobDeployCommands).Run()
	if err != nil {
		return err
	}

	//输出结果
	return r.output(machines, instanceMoveCommands, jobDeployCommands)
}
