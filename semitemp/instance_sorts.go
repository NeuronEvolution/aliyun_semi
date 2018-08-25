package semitemp

import "sort"

//92,288,1024,7,7,9
func SortInstanceByTotalMaxHigh(p []*Instance) {
	sort.Slice(p, func(i, j int) bool {
		c1 := p[i].Config
		c2 := p[j].Config
		a1 := float64(c1.Disk)/float64(HighDisk) + c1.CpuMax/float64(HighCpu*MaxCpuRatio) + c1.MemMax/float64(HighMem) +
			float64(c1.P)/float64(7) + float64(c1.M)/float64(7) + float64(c1.PM)/float64(9)
		a2 := float64(c2.Disk)/float64(HighDisk) + c2.CpuMax/float64(HighCpu*MaxCpuRatio) + c2.MemMax/float64(HighMem) +
			float64(c2.P)/float64(7) + float64(c2.M)/float64(7) + float64(c2.PM)/float64(9)
		return a1 > a2
	})
}

//32,64,1440,7,3,7
func SortInstanceByTotalMaxLowWithInference(p []*Instance, inferenceLimit int) {
	sort.Slice(p, func(i, j int) bool {
		c1 := p[i].Config
		c2 := p[j].Config
		a1 := float64(c1.Disk)/float64(LowDisk) + c1.CpuMax/float64(LowCpu*MaxCpuRatio) + c1.MemMax/float64(LowMem) +
			float64(c1.P)/float64(7) + float64(c1.M)/float64(3) + float64(c1.PM)/float64(7)
		a2 := float64(c2.Disk)/float64(LowDisk) + c2.CpuMax/float64(LowCpu*MaxCpuRatio) + c2.MemMax/float64(LowMem) +
			float64(c2.P)/float64(7) + float64(c2.M)/float64(3) + float64(c2.PM)/float64(7)

		if c1.InferenceAppCount < inferenceLimit && c2.InferenceAppCount < inferenceLimit {
			return a1 > a2
		}

		if c1.InferenceAppCount > c2.InferenceAppCount {
			return true
		} else if c1.InferenceAppCount == c2.InferenceAppCount {
			return a1 > a2
		} else {
			return false
		}
	})
}

func SortInstanceByTotalMaxLow(p []*Instance) {
	sort.Slice(p, func(i, j int) bool {
		c1 := p[i].Config
		c2 := p[j].Config
		a1 := float64(c1.Disk)/float64(LowDisk) + c1.CpuMax/float64(LowCpu*MaxCpuRatio) + c1.MemMax/float64(LowMem) +
			float64(c1.P)/float64(7) + float64(c1.M)/float64(3) + float64(c1.PM)/float64(7)
		a2 := float64(c2.Disk)/float64(LowDisk) + c2.CpuMax/float64(LowCpu*MaxCpuRatio) + c2.MemMax/float64(LowMem) +
			float64(c2.P)/float64(7) + float64(c2.M)/float64(3) + float64(c2.PM)/float64(7)
		return a1 > a2
	})
}

func SortInstanceByDisk(p []*Instance) {
	sort.Slice(p, func(i, j int) bool {
		if p[i].Config.Disk > p[j].Config.Disk {
			return true
		} else if p[i].Config.Disk == p[j].Config.Disk {
			a1 := p[i].Config.CpuMax*float64(HighMem) + p[i].Config.MemMax*float64(HighCpu*MaxCpuRatio)
			a2 := p[j].Config.CpuMax*float64(HighMem) + p[j].Config.MemMax*float64(HighCpu*MaxCpuRatio)
			if a1 > a2 {
				return true
			} else if a1 == a2 {
				if p[i].Config.AppId < p[j].Config.AppId {
					return true
				} else {
					return false
				}
			} else {
				return false
			}
		} else {
			return false
		}
	})
}

func SortInstanceByCpu(p []*Instance) {
	sort.Slice(p, func(i, j int) bool {
		if p[i].Config.CpuMax > p[j].Config.CpuMax {
			return true
		} else if p[i].Config.CpuMax > p[j].Config.CpuMax {
			a1 := p[i].Config.MemMax*float64(HighDisk) + float64(p[i].Config.Disk)*float64(HighMem)
			a2 := p[j].Config.MemMax*float64(HighDisk) + float64(p[j].Config.Disk)*float64(HighMem)
			if a1 > a2 {
				return true
			} else if a1 == a2 {
				if p[i].Config.AppId < p[j].Config.AppId {
					return true
				} else {
					return false
				}
			} else {
				return false
			}
		} else {
			return false
		}
	})
}

func SortInstanceByMem(p []*Instance) {
	sort.Slice(p, func(i, j int) bool {
		if p[i].Config.MemMax > p[j].Config.MemMax {
			return true
		} else if p[i].Config.MemMax > p[j].Config.MemMax {
			a1 := p[i].Config.CpuMax*float64(HighDisk) + float64(p[i].Config.Disk)*float64(HighCpu*MaxCpuRatio)
			a2 := p[j].Config.CpuMax*float64(HighDisk) + float64(p[j].Config.Disk)*float64(HighCpu*MaxCpuRatio)
			if a1 > a2 {
				return true
			} else if a1 == a2 {
				if p[i].Config.AppId < p[j].Config.AppId {
					return true
				} else {
					return false
				}
			} else {
				return false
			}
		} else {
			return false
		}
	})
}
