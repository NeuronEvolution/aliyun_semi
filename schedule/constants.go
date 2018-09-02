package schedule

//4600-0.01-4978.402645
//4700-0.01-4868.195688
const MachineA = 4800

//分开5129.502338
//4600-0.01-5216.240981
//4700-0.01-5014.719198
const MachineB = 4800

//6500-0.01-7137.349613
//6600-0.01-7142.123558
//6700-0.01-7167.983697
const MachineC = 6800

//6600-0.01-7089.672663
//6700-0.01-7120.390462
const MachineD = 6800

const MachineE = 8000
const MachineALoop = 4096
const MachineBLoop = 4096
const MachineCLoop = 8192
const MachineDLoop = 8192

const TimeSampleCount = 98
const ScaleLimitH1 = 1.01
const ScaleLimitH2 = 1.1
const ScaleLimitH3 = 1.2
const ScaleRatioH1 = 400 //400
const ScaleRatioH2 = 40  //40
const ScaleRatioH3 = 20  //20
const ScaleBase = 128
const ScaleRatio = 1.414
const MaxAppId = 10000
const MaxInstancePerMachine = 256
const MaxAppPerMachine = 256
const MaxJobPerMachine = 65536
const MaxCpuRatio = float64(0.5)
const ConstraintE = float64(0.0000001)
const ParallelCpuCount = 6
const JobScheduleCpuLimitStep = 0.01
const JobPackCpu = 2
const JobPackMem = 4

//92,288,2457,7,7,9
const HighCpu = 92
const HighMem = 288
const HighDisk = 2457

//32,64,1440,7,3,7
const LowCpu = 32
const LowMem = 64
const LowDisk = 1440
