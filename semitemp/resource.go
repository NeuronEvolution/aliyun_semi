package semitemp

import (
	"fmt"
	"math"
)

type Resource struct {
	Cpu  [TimeSampleCount]float64
	Mem  [TimeSampleCount]float64
	Disk int
	P    int
	M    int
	PM   int

	ResourceCost          float64
	ResourceCostDeviation float64

	CpuMax       float64
	CpuMin       float64
	CpuAvg       float64
	CpuDeviation float64
	MemMax       float64
	MemMin       float64
	MemAvg       float64
	MemDeviation float64

	cpuCostValid bool
	CpuCost      float64
}

func (r *Resource) Reset() {
	for i := 0; i < len(r.Cpu); i++ {
		r.Cpu[i] = 0
	}

	for i := 0; i < len(r.Mem); i++ {
		r.Mem[i] = 0
	}

	r.Disk = 0
	r.P = 0
	r.M = 0
	r.PM = 0

	r.CpuCost = 0
	r.cpuCostValid = true
}

func (r *Resource) DebugPrint() {
	r.CalcTimedResourceStatistics()

	fmt.Printf("Resource.Print ")
	fmt.Printf("Disk=%4d,", r.Disk)
	fmt.Printf("P=%d,M=%d,PM=%d,", r.P, r.M, r.PM)
	fmt.Printf("CpuAvg=%4.1f,CpuDev=%4.1f,CpuMin=%4.1f,CpuMax=%4.1f,",
		r.CpuAvg, r.CpuDeviation, r.CpuMin, r.CpuMax)
	fmt.Printf("MemAvg=%5.1f,MemDev=%5.1f,MemMin=%5.1f,MemMax=%5.1f\n",
		r.MemAvg, r.MemDeviation, r.MemMin, r.MemMax)
}

func (r *Resource) AddResource(p *Resource) {
	for i, v := range p.Cpu {
		r.Cpu[i] += v
	}
	for i, v := range p.Mem {
		r.Mem[i] += v
	}
	r.Disk += p.Disk
	r.M += p.M
	r.P += p.P
	r.PM += p.PM

	r.cpuCostValid = false
}

func (r *Resource) RemoveResource(p *Resource) {
	for i, v := range p.Cpu {
		r.Cpu[i] -= v
	}
	for i, v := range p.Mem {
		r.Mem[i] -= v
	}
	r.Disk -= p.Disk
	r.M -= p.M
	r.P -= p.P
	r.PM -= p.PM

	r.cpuCostValid = false
}

func (r *Resource) CalcTimedResourceStatistics() {
	r.CpuAvg, r.CpuDeviation, r.CpuMin, r.CpuMax = Statistics(r.Cpu)
	r.MemAvg, r.MemDeviation, r.MemMin, r.MemMax = Statistics(r.Mem)
}

func (r *Resource) calcCostEval(config *MachineLevelConfig) {
	avgCpu := float64(0)
	for _, v := range r.Cpu {
		avgCpu += v
	}
	avgCpu = avgCpu / float64(len(r.Cpu))

	avgMem := float64(0)
	for _, v := range r.Mem {
		avgMem += v
	}
	avgMem = avgMem / float64(len(r.Mem))

	cpu := avgCpu * ParamCpuCostMultiply / config.Cpu
	mem := avgMem / config.Mem
	disk := float64(r.Disk) / float64(config.Disk)
	p := float64(r.P) / float64(config.P)
	m := float64(r.M) / float64(config.M)
	pm := float64(r.PM) / float64(config.PM)

	r.ResourceCost = scaleCost(cpu) +
		scaleCost(mem) +
		scaleCost(disk) +
		scaleCost(p) +
		scaleCost(m) +
		scaleCost(pm)

	r.ResourceCostDeviation = calcResourceCostDeviation(cpu, mem, disk, p, m, pm)
}

func (r *Resource) GetCpuDerivation() float64 {
	avg := float64(0)
	for i := 0; i < TimeSampleCount; i++ {
		ratio := r.Cpu[i] / HighCpu
		avg += ratio
	}

	avg = avg / float64(TimeSampleCount)
	d := float64(0)
	for i := 0; i < TimeSampleCount; i++ {
		ratio := r.Cpu[i] / HighCpu
		d += (ratio - avg) * (ratio - avg)
	}
	d = math.Sqrt(d / TimeSampleCount)

	return d
}

func (r *Resource) GetCpuCost(cpuLimit float64) float64 {
	if r.cpuCostValid {
		return r.CpuCost
	}

	totalCost := float64(0)
	for i := 0; i < TimeSampleCount; i++ {
		r := r.Cpu[i] / cpuLimit
		if r > 0.5 {
			totalCost += 1 + 10*(Exp(r-0.5)-1)
		} else {
			totalCost += 1
		}
	}

	r.CpuCost = totalCost / TimeSampleCount
	r.cpuCostValid = true

	return r.CpuCost
}

func (r *Resource) GetLinearCpuCost(cpuLimit float64) float64 {
	totalCost := float64(0)
	for i := 0; i < TimeSampleCount; i++ {
		r := r.Cpu[i] / cpuLimit
		if r > 0.5 {
			totalCost += 1 + 10*(Exp(r-0.5)-1)
		} else {
			totalCost += r * 2
		}
	}

	return totalCost / TimeSampleCount
}

func (r *Resource) GetCostWithInstance(instance *Instance, cpuLimit float64) float64 {
	totalCost := float64(0)
	for i := 0; i < TimeSampleCount; i++ {
		r := (r.Cpu[i] + instance.Config.Cpu[i]) / cpuLimit
		if r > 0.5 {
			totalCost += 1 + 10*(Exp(r-0.5)-1)
		} else {
			totalCost += r * 2
		}
	}

	return totalCost / TimeSampleCount
}

func calcResourceCostDeviation(cpu float64, mem float64, disk float64, p float64, m float64, pm float64) float64 {
	avg := (cpu + mem + disk + p + m + pm) / 6
	return math.Sqrt(((cpu-avg)*(cpu-avg) + (mem-avg)*(mem-avg) + (disk-avg)*(disk-avg) +
		(p-avg)*(p-avg) + (m-avg)*(m-avg) + (pm-avg)*(pm-avg)) / float64(6))
}

func scaleCost(f float64) float64 {
	return Exp(ParamAppCostMultiply * f)
}
