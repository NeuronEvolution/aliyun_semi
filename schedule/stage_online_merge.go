package schedule

func (r *ResourceManagement) onlineMerge() (err error) {
	r.log("onlineMerge start\n")

	return NewOnlineMerge(r).Run()
}
