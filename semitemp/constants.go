package semitemp

const TimeSampleCount = 98

const MaxAppId = 10000
const MaxInstanceId = 100000
const MaxMachineId = 6000 + 1
const MachineCpuMax = 92
const MachineMemMax = 288
const MachineDiskMax = 1024
const MachinePMax = 7
const MachineMMax = 7
const MachinePMMax = 9

const MaxDeployCommandCount = 1024 * 1024
const MaxInstancePerMachine = 1024
const MaxAppPerMachine = 64

const ParamMachineCostMultiply = 4 //影响大
const ParamDeviationMultiply = 1   //影响大
const ParamAppInferenceMultiply = 1
const ParamAppCostMultiply = 5 //影响大
const ParamCpuCostMultiply = 1

const MaxCpuRatio = float64(0.5)

const ConstraintE = float64(0.0000001)

//92,288,2457,7,7,9
const HighCpu = 92
const HighMem = 288
const HighDisk = 2457

//32,64,1440,7,3,7
const LowCpu = 32
const LowMem = 64
const LowDisk = 1440
