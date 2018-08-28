package schedule

const TimeSampleCount = 98
const ScaleLimitH1 = 1.01
const ScaleLimitH2 = 1.1
const ScaleLimitH3 = 1.2
const ScaleRatioH1 = 200 //400
const ScaleRatioH2 = 20  //40
const ScaleRatioH3 = 10  //20
const ScaleBase = 128
const ScaleRatio = 1.414
const MaxAppId = 10000
const MaxInstancePerMachine = 256
const MaxAppPerMachine = 256
const MaxJobPerMachine = 4096
const MaxCpuRatio = float64(0.5)
const ConstraintE = float64(0.0000001)
const MaxJobExecMinutes = 144

//92,288,2457,7,7,9
const HighCpu = 92
const HighMem = 288
const HighDisk = 2457

//32,64,1440,7,3,7
const LowCpu = 32
const LowMem = 64
const LowDisk = 1440
