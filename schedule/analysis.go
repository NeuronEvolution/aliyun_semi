package schedule

import "fmt"

func (r *ResourceManagement) analysis() {
	n := 0
	for _, v := range r.JobConfigMap {
		if v != nil {
			if v.Mem > 16 {
				fmt.Println(v.InstanceCount, v.Cpu, v.Mem, v.ExecMinutes)
				n += v.InstanceCount
			}
		}
	}
	r.log("%d\n", n)

	return

	jobTrees := 0
	for _, v := range r.JobConfigDAG {
		if len(v.Children) > 0 {
			jobTrees++
		}
	}
	r.log("instances=%d\n", len(r.InstanceDeployConfigList))
	r.log("jobs=%d,rootJobs=%d,trees=%d\n", len(r.JobConfigMap), len(r.JobConfigDAG), jobTrees)
	r.log("maxMachineId=%d,maxInstanceId=%d,maxJobInstanceId=%d\n",
		r.MaxMachineId, r.MaxInstanceId, r.MaxJobInstanceId)
}
