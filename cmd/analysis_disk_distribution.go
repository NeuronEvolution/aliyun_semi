package main

import (
	"fmt"
	"sort"
)

type diskCount struct {
	disk  int
	count int
}

type appCount struct {
	appId int
	count int
}

type diskAppCount struct {
	disk         int
	appCountList []*appCount
}

func (c *AnalysisContext) AnalysisDiskDistributionByApp() {
	fmt.Printf("AnalysisDiskDistributionByApp\n")
	appDiskDist := make(map[int]int)
	for _, a := range c.appResourcesList {
		_, has := appDiskDist[a.Disk]
		if !has {
			appDiskDist[a.Disk] = 1
		} else {
			appDiskDist[a.Disk]++
		}
	}
	appDiskCountList := make([]*diskCount, 0)
	for i, v := range appDiskDist {
		appDiskCountList = append(appDiskCountList, &diskCount{disk: i, count: v})
	}
	sort.Slice(appDiskCountList, func(i, j int) bool {
		return appDiskCountList[i].disk > appDiskCountList[j].disk
	})
	for _, v := range appDiskCountList {
		fmt.Printf("    disk=%5d,count=%7d\n", v.disk, v.count)
	}
}

func (c *AnalysisContext) AnalysisDiskDistributionByInstance() {
	fmt.Printf("AnalysisDiskDistributionByInstance\n")
	instanceDiskDist := make(map[int]int)
	for _, v := range c.instanceDeployList {
		appResource := c.appResourcesMap[v.AppId]
		_, has := instanceDiskDist[appResource.Disk]
		if !has {
			instanceDiskDist[appResource.Disk] = 1
		} else {
			instanceDiskDist[appResource.Disk]++
		}
	}
	instanceDiskCountList := make([]*diskCount, 0)
	for i, v := range instanceDiskDist {
		instanceDiskCountList = append(instanceDiskCountList, &diskCount{disk: i, count: v})
	}
	sort.Slice(instanceDiskCountList, func(i, j int) bool {
		return instanceDiskCountList[i].disk > instanceDiskCountList[j].disk
	})
	for _, v := range instanceDiskCountList {
		fmt.Printf("    disk=%5d,count=%7d\n", v.disk, v.count)
	}
}

func (c *AnalysisContext) AnalysisDiskDistributionByAppInstance() {
	fmt.Printf("AnalysisDiskDistributionByAppInstance\n")
	diskAppInstanceCountMap := make(map[int]map[int]int)
	for _, v := range c.instanceDeployList {
		appResource := c.appResourcesMap[v.AppId]
		_, hasDisK := diskAppInstanceCountMap[appResource.Disk]
		if !hasDisK {
			diskAppInstanceCountMap[appResource.Disk] = make(map[int]int)
		}
		_, hasApp := diskAppInstanceCountMap[appResource.Disk][v.AppId]
		if !hasApp {
			diskAppInstanceCountMap[appResource.Disk][v.AppId] = 1
		} else {
			diskAppInstanceCountMap[appResource.Disk][v.AppId]++
		}
	}

	diskAppCountList := make([]*diskAppCount, 0)
	for disk, v := range diskAppInstanceCountMap {
		diskAppCount := &diskAppCount{}
		diskAppCount.disk = disk
		diskAppCount.appCountList = make([]*appCount, 0)
		for appId, count := range v {
			diskAppCount.appCountList = append(diskAppCount.appCountList, &appCount{appId: appId, count: count})
		}
		diskAppCountList = append(diskAppCountList, diskAppCount)
	}
	sort.Slice(diskAppCountList, func(i, j int) bool {
		return diskAppCountList[i].disk > diskAppCountList[j].disk
	})
	for _, v := range diskAppCountList {
		sort.Slice(v.appCountList, func(i, j int) bool {
			return v.appCountList[i].count > v.appCountList[j].count
		})
	}
	for _, v := range diskAppCountList {
		totalCount := 0
		for _, appCount := range v.appCountList {
			totalCount += appCount.count
		}
		fmt.Printf("    disk=%d,count=%d\n", v.disk, totalCount)

		for _, appCount := range v.appCountList {
			a := c.appResourcesMap[appCount.appId]
			if v.disk <= 100 && a.CpuAvg < 16 && a.MemAvg < 32 {
				continue
			}

			fmt.Printf("    CpuAvg=%4.1f,CpuDev=%4.1f,CpuMin=%4.1f,CpuMax=%4.1f,",
				a.CpuAvg, a.CpuDeviation, a.CpuMin, a.CpuMax)
			fmt.Printf("MemAvg=%5.1f,MemDev=%5.1f,MemMin=%5.1f,MemMax=%5.1f,",
				a.MemAvg, a.MemDeviation, a.MemMin, a.MemMax)
			fmt.Printf("appId=%5d,count=%7d\n", appCount.appId, appCount.count)
		}
	}
}

func (c *AnalysisContext) AnalysisMemAvgDistributionByInstance() {
	fmt.Printf("AnalysisMemAvgDistributionByInstance\n")
	instanceMemDist := make(map[int]int)
	for _, v := range c.instanceDeployList {
		appResource := c.appResourcesMap[v.AppId]
		_, has := instanceMemDist[int(appResource.MemAvg)]
		if !has {
			instanceMemDist[int(appResource.MemAvg)] = 1
		} else {
			instanceMemDist[int(appResource.MemAvg)]++
		}
	}
	instanceDiskCountList := make([]*diskCount, 0)
	for i, v := range instanceMemDist {
		instanceDiskCountList = append(instanceDiskCountList, &diskCount{disk: i, count: v})
	}
	sort.Slice(instanceDiskCountList, func(i, j int) bool {
		return instanceDiskCountList[i].disk > instanceDiskCountList[j].disk
	})
	for _, v := range instanceDiskCountList {
		fmt.Printf("    memAvg=%5d,count=%7d\n", v.disk, v.count)
	}
}

func (c *AnalysisContext) AnalysisMemMaxDistributionByInstance() {
	fmt.Printf("AnalysisMemMaxDistributionByInstance\n")
	instanceMemDist := make(map[int]int)
	for _, v := range c.instanceDeployList {
		appResource := c.appResourcesMap[v.AppId]
		_, has := instanceMemDist[int(appResource.MemMax)]
		if !has {
			instanceMemDist[int(appResource.MemMax)] = 1
		} else {
			instanceMemDist[int(appResource.MemMax)]++
		}
	}
	instanceDiskCountList := make([]*diskCount, 0)
	for i, v := range instanceMemDist {
		instanceDiskCountList = append(instanceDiskCountList, &diskCount{disk: i, count: v})
	}
	sort.Slice(instanceDiskCountList, func(i, j int) bool {
		return instanceDiskCountList[i].disk > instanceDiskCountList[j].disk
	})
	for _, v := range instanceDiskCountList {
		fmt.Printf("    memAvg=%5d,count=%7d\n", v.disk, v.count)
	}
}

func (c *AnalysisContext) AnalysisMemSameDistributionByInstance() {
	fmt.Printf("AnalysisMemSameDistributionByInstance\n")
	instanceMemDist := make(map[int]int)
	for _, v := range c.instanceDeployList {
		appResource := c.appResourcesMap[v.AppId]
		if (appResource.MemMax - appResource.MemAvg) > 0.001 {
			continue
		}

		_, has := instanceMemDist[int(appResource.MemMax)]
		if !has {
			instanceMemDist[int(appResource.MemMax)] = 1
		} else {
			instanceMemDist[int(appResource.MemMax)]++
		}
	}
	instanceDiskCountList := make([]*diskCount, 0)
	for i, v := range instanceMemDist {
		instanceDiskCountList = append(instanceDiskCountList, &diskCount{disk: i, count: v})
	}
	sort.Slice(instanceDiskCountList, func(i, j int) bool {
		return instanceDiskCountList[i].disk > instanceDiskCountList[j].disk
	})
	for _, v := range instanceDiskCountList {
		fmt.Printf("    memAvg=%5d,count=%7d\n", v.disk, v.count)
	}
}

func (c *AnalysisContext) AnalysisDiskDistribution() {
	c.AnalysisDiskDistributionByApp()
	c.AnalysisDiskDistributionByInstance()
	c.AnalysisDiskDistributionByAppInstance()

	c.AnalysisMemAvgDistributionByInstance()
	c.AnalysisMemMaxDistributionByInstance()
	c.AnalysisMemSameDistributionByInstance()
}
