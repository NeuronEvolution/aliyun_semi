package schedule

const JobScheduleCpuLimitStep = 0.01

//4600-0.01-4978.402645
//4700-0.01-4868.195688
//4750-0.01-4850.778961 0.02-4866.668925 0.005-4839.375240
//4800-0.01-4859.439903
const MachineA = 4750

//分开5129.502338
//4600-0.01-5216.240981
//4700-0.01-5014.719198
//4800-0.01-4913.862059 0.02-4928.812726 0.005-4909.216453
//4850-0.01-4915.017962
const MachineB = 4800

//6300-0.01-7129.214311 0.02-7137.326849 0.005-7127.638963
//6400-0.01-7134.548322
//6500-0.01-7137.349613
//6600-0.01-7142.123558
//6700-0.01-7167.983697
//6800-0.01-7184.881388
const MachineC = 6300

//6400-0.01-7037.418658 0.02-7037.368647 0.005-7038.260330
//6500-0.01-7068.338875
//6600-0.01-7089.672663
//6700-0.01-7120.390462
//6800-0.01-7154.584273
const MachineD = 6400

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
const MaxInstancePerMachine = 256
const MaxAppPerMachine = 256
const MaxJobPerMachine = 4096
const MaxCpuRatio = float64(0.5)
const ConstraintE = float64(0.0000001)
const ParallelCpuCount = 6
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
