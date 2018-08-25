package semitemp

import (
	"fmt"
	"sort"
)

type DiskCount struct {
	Disk  int
	Count int
}

func AnalysisDiskDistributionByInstance(instanceList []*Instance) {
	fmt.Printf("AnalysisDiskDistributionByInstance instanceCount=%d\n", len(instanceList))
	instanceDiskDist := make(map[int]int)
	for _, v := range instanceList {
		_, has := instanceDiskDist[v.Config.Disk]
		if !has {
			instanceDiskDist[v.Config.Disk] = 1
		} else {
			instanceDiskDist[v.Config.Disk]++
		}
	}
	instanceDiskCountList := make([]*DiskCount, 0)
	for i, v := range instanceDiskDist {
		instanceDiskCountList = append(instanceDiskCountList, &DiskCount{Disk: i, Count: v})
	}
	sort.Slice(instanceDiskCountList, func(i, j int) bool {
		return instanceDiskCountList[i].Disk > instanceDiskCountList[j].Disk
	})
	for _, v := range instanceDiskCountList {
		fmt.Printf("    disk=%5d,count=%7d\n", v.Disk, v.Count)
	}
}

func AnalysisMemMaxDistributionByInstance(instanceList []*Instance) {
	fmt.Printf("AnalysisMemDistributionByInstance instanceCount=%d\n", len(instanceList))
	instanceMemDist := make(map[int]int)
	for _, v := range instanceList {
		_, has := instanceMemDist[int(v.Config.MemMax)]
		if !has {
			instanceMemDist[int(v.Config.MemMax)] = 1
		} else {
			instanceMemDist[int(v.Config.MemMax)]++
		}
	}
	instanceMemCountList := make([]*DiskCount, 0)
	for i, v := range instanceMemDist {
		instanceMemCountList = append(instanceMemCountList, &DiskCount{Disk: i, Count: v})
	}
	sort.Slice(instanceMemCountList, func(i, j int) bool {
		return instanceMemCountList[i].Disk > instanceMemCountList[j].Disk
	})
	for _, v := range instanceMemCountList {
		fmt.Printf("    memMax=%5d,count=%7d\n", v.Disk, v.Count)
	}
}

func AnalysisMemAvgDistributionByInstance(instanceList []*Instance) {
	fmt.Printf("AnalysisMemAvgDistributionByInstance instanceCount=%d\n", len(instanceList))
	instanceMemDist := make(map[int]int)
	for _, v := range instanceList {
		_, has := instanceMemDist[int(v.Config.MemAvg)]
		if !has {
			instanceMemDist[int(v.Config.MemAvg)] = 1
		} else {
			instanceMemDist[int(v.Config.MemAvg)]++
		}
	}
	instanceMemCountList := make([]*DiskCount, 0)
	for i, v := range instanceMemDist {
		instanceMemCountList = append(instanceMemCountList, &DiskCount{Disk: i, Count: v})
	}
	sort.Slice(instanceMemCountList, func(i, j int) bool {
		return instanceMemCountList[i].Disk > instanceMemCountList[j].Disk
	})
	for _, v := range instanceMemCountList {
		fmt.Printf("    memMax=%5d,count=%7d\n", v.Disk, v.Count)
	}
}
