package main

import (
	"fmt"
	"sort"
)

func (c *AnalysisContext) AnalysisAppInference() {
	sort.Slice(c.appResourcesList, func(i, j int) bool {
		return c.appResourcesList[i].InferenceAppCount > c.appResourcesList[j].InferenceAppCount
	})

	n := 0
	for _, resource := range c.appResourcesList {
		if resource.InferenceAppCount < 2 {
			continue
		}

		fmt.Printf("interence count=%d app=%d  ", resource.InferenceAppCount, resource.AppId)
		resource.DebugPrint()

		instanceCount := 0
		for _, instance := range c.instanceDeployList {
			if instance.AppId == resource.AppId {
				instanceCount++
				n++
			}
		}
		fmt.Println("instanceCount", instanceCount)

		for _, inference := range c.appInterferenceList {
			if inference.AppId1 == resource.AppId && inference.AppId2 == resource.AppId {
				fmt.Println("inf self", inference.Interference)
				break
			}
		}

		for _, inference := range c.appInterferenceList {
			if inference.AppId1 == resource.AppId || inference.AppId2 == resource.AppId {
				if inference.AppId1 == resource.AppId && (inference.AppId2 == 798 ||
					inference.AppId2 == 6659 ||
					inference.AppId2 == 4825 ||
					inference.AppId2 == 536 ||
					inference.AppId2 == 1900 ||
					inference.AppId2 == 4026 ||
					inference.AppId2 == 6381) || inference.AppId2 == resource.AppId && (inference.AppId1 == 798 ||
					inference.AppId1 == 6659 ||
					inference.AppId1 == 4825 ||
					inference.AppId1 == 536 ||
					inference.AppId1 == 1900 ||
					inference.AppId1 == 4026 ||
					inference.AppId1 == 6381) {
					fmt.Println("inf ", inference.AppId1, inference.AppId2, inference.Interference)
				}
			}
		}
	}

	fmt.Println("AnalysisAppInference total instance count by inference limit", n)
}
