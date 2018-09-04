package schedule

//todo base从0.45开始，step再降低
const JobScheduleCpuLimitStep = 0.005 //任务调度水平线上移步长
const JobPackCpu = 2                  //任务打包Cpu上限
const JobPackMem = 4                  //任务打包Mem上限
const JobPackLimit = 320000           //任务打包限制，总数低于此不打包
const ScaleLimitH1 = 1.01             //实例最优机器数量增长-统计因子Cpu
const ScaleLimitH2 = 1.1              //实例最优机器数量增长-统计因子Cpu
const ScaleLimitH3 = 1.2              //实例最优机器数量增长-统计因子Cpu
const ScaleRatioH1 = 400              //实例最优机器数量增长-统计因子数量
const ScaleRatioH2 = 40               //实例最优机器数量增长-统计因子数量
const ScaleRatioH3 = 20               //实例最优机器数量增长-统计因子数量
const ScaleBase = 128                 //实例最优机器数量增长-每次增长后迭代次数基数
const ScaleRatio = 1.414              //实例最优机器数量增长-每次增长后迭代次数增长指数

//4600|0.01-4978.402645
//4700|0.01-4868.195688                 |0.005-4861.169780
//4725|                                 |0.005-4851.899044
//4750|0.01-4850.778961|0.02-4866.668925|0.005-4839.375240,(1,2)4843.341124
//4775|
//4800|0.01-4859.439903
const MachineA = 4775 //机器A数量

//4600-0.01-5216.240981
//4700-0.01-5014.719198
//4750-                                  0.005-4930.245791
//4775-                                  0.005-
//4800-0.01-4913.862059 0.02-4928.812726 0.005-4909.216453
//4850-0.01-4915.017962
const MachineB = 4775 //机器B数量

//6000|                                 |0.005-7264.363959
//6200|
//6300|0.01-7129.214311|0.02-7137.326849|0.005-7127.638963
//6400|0.01-7134.548322
//6500|0.01-7137.349613
//6600|0.01-7142.123558
//6700|0.01-7167.983697
//6800|0.01-7184.881388
const MachineC = 6200 //机器C数量

//5900|7011.106186                                  |0.005-7067.784797
//6000|                                             |0.005-7023.596781
//6400|           |0.01-7037.418658|0.02-7037.368647|0.005-7038.260330
//6500|           |0.01-7068.338875|
//6600||0.01-7089.672663|
//6700||0.01-7120.390462|
//6800||0.01-7154.584273|
const MachineD = 5900 //机器D数量

const MachineE = 8000 //机器E数量

const MachineALoop = 4096 //实例调度迭代次数A，可以统一为8192或更高
const MachineBLoop = 4096 //实例调度迭代次数B，可以统一为8192或更高
const MachineCLoop = 8192 //实例调度迭代次数C，可以统一为8192或更高
const MachineDLoop = 8192 //实例调度迭代次数D，可以统一为8192或更高

const TimeSampleCount = 98             //实例时间点数量
const MaxInstancePerMachine = 256      //单每机器最大实例数
const MaxAppPerMachine = 256           //单机器最大应用数
const MaxJobPerMachine = 4096          //单机器最大任务数
const ConstraintE = float64(0.0000001) //浮点判断浮动
const ParallelCpuCount = 6             //并行计算Cpu数量

//92,288,2457,7,7,9
const HighCpu = 92
const HighMem = 288
const HighDisk = 2457

//32,64,1440,7,3,7
const LowCpu = 32
const LowMem = 64
const LowDisk = 1440
