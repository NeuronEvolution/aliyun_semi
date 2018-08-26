package schedule

func (r *ResourceManagement) analysis() {
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
