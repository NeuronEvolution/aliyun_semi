package main

import "sort"

func SortMachineByCpuCost(p []*Machine) {
	sort.Slice(p, func(i, j int) bool {
		return p[i].GetCpuCost() > p[j].GetCpuCost()
	})
}
