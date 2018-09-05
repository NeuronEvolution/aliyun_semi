package schedule

//todo base从0.45开始，step再降低
const JobScheduleCpuLimitStep = 0.005 //任务调度水平线上移步长
const JobPackCpu = 1                  //任务打包Cpu上限
const JobPackMem = 2                  //任务打包Mem上限
const JobPackLimit = 640000           //任务打包限制，总数低于此不打包
const ScaleLimitH1 = 1.01             //实例最优机器数量增长-统计因子Cpu
const ScaleLimitH2 = 1.1              //实例最优机器数量增长-统计因子Cpu
const ScaleLimitH3 = 1.2              //实例最优机器数量增长-统计因子Cpu
const ScaleRatioH1 = 400              //实例最优机器数量增长-统计因子数量
const ScaleRatioH2 = 40               //实例最优机器数量增长-统计因子数量
const ScaleRatioH3 = 20               //实例最优机器数量增长-统计因子数量
const ScaleBase = 128                 //实例最优机器数量增长-每次增长后迭代次数基数
const ScaleRatio = 1.414              //实例最优机器数量增长-每次增长后迭代次数增长指数

//4400|4784.511507|                                 |0.005-5619.464343|1-5619.464343
//4600|4716.993993|0.01-4978.402645
//4700|4745.163044|0.01-4868.195688                 |0.005-4861.169780
//4725|4757.255831                                  |0.005-4851.899044
//4750|4772.426304|0.01-4850.778961|0.02-4866.668925|0.005-4839.375240|1-4818.295249|2-4807.713511
//4775|4789.694502                      |0.005-4841.233585
//4800|4809.162456|0.01-4859.439903
const MachineA = 4200 //机器A数量

//4400|4783.167046|                                 |0.005-5964.154670|1-5772.231651
//4600|4716.936160|0.01-5216.240981|1-5073.877344
//4700|4743.572720|0.01-5014.719198                 |0.005-4993.132891|8-4803.619554
//4750|4772.426304|                                 |0.005-4930.245791|
//4775|4788.754447                                  |0.005-4920.215002
//4800|4809.488216|0.01-4913.862059|0.02-4928.812726|0.005-4909.216453|1-4860.152612|2-4847.329039
//4850|4852.117385|0.01-4915.017962
const MachineB = 4200 //机器B数量

//5900|7002.015487|                                 |0.005-7323.966549|2-7233.322033
//6000|7004.522488                                  |0.005-7264.363959
//6200|7006.543782|                                 |0.005-7155.616830
//6250|7022.131489|                                 |0.005-7153.724555|4-7055.944702
//6300|7013.023175|0.01-7129.214311|0.02-7137.326849|0.005-7127.638963|1-7096.347922
//6400|7038.654348|0.01-7134.548322
//6500|7058.643933|0.01-7137.349613
//6600|7088.231469|0.01-7142.123558
//6700|7123.407200|0.01-7167.983697
//6800|7158.672580|0.01-7184.881388
const MachineC = 5700 //机器C数量

//5900|7011.106186                                  |0.005-7067.784797
//6000|6982.099893|                                 |0.005-7023.596781|1-7000.260552|22-6988.707484
//6100|6995.932736|0.005-7024.917345|20-7000.368497
//6400|7029.473588|0.01-7037.418658|0.02-7037.368647|0.005-7038.260330
//6500|7064.360492|0.01-7068.338875|
//6600|7086.582644|0.01-7089.672663|
//6700|7117.632487|0.01-7120.390462|
//6800|7152.028460|0.01-7154.584273|
const MachineD = 5700 //机器D数量

//8000|8448.281015
const MachineE = 8000 //机器E数量

const MachineALoop = 8192  //实例调度迭代次数A，可以统一为8192或更高
const MachineBLoop = 8192  //实例调度迭代次数B，可以统一为8192或更高
const MachineCLoop = 32768 //实例调度迭代次数C，可以统一为8192或更高
const MachineDLoop = 32768 //实例调度迭代次数D，可以统一为8192或更高
const MachineELoop = 65536 * 32

const TimeSampleCount = 98             //实例时间点数量
const MaxInstancePerMachine = 256      //单每机器最大实例数
const MaxAppPerMachine = 256           //单机器最大应用数
const MaxJobPerMachine = 4096          //单机器最大任务数
const ConstraintE = float64(0.0000001) //浮点判断浮动
const ParallelCpuCount = 12            //并行计算Cpu数量

//92,288,2457,7,7,9
const HighCpu = 92
const HighMem = 288
const HighDisk = 2457

//32,64,1440,7,3,7
const LowCpu = 32
const LowMem = 64
const LowDisk = 1440
