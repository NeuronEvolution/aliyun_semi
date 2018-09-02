package schedule

func (r *ResourceManagement) newLoop() {
	machineCount := 0
	for _, m := range r.MachineList {
		if m.InstanceListCount > 0 || m.JobListCount > 0 {
			machineCount++
		}
	}

	r.log("newLoop machineCount=%d\n", machineCount)

	return
}
