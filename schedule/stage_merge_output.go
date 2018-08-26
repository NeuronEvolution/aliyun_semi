package schedule

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

func (r *ResourceManagement) buildJobDeployCommands() (commands []*JobDeployCommand) {
	commands = make([]*JobDeployCommand, 0)
	for _, m := range r.MachineList {
		if m.JobListCount == 0 {
			continue
		}

		for _, job := range m.JobList[:m.JobListCount] {
			commands = append(commands, &JobDeployCommand{
				JobId:        job.Config.RealJobId,
				MachineId:    m.MachineId,
				Count:        job.InstanceCount,
				StartMinutes: job.StartMinutes,
			})
		}
	}

	return commands
}

func (r *ResourceManagement) output(instanceMoveCommands []*InstanceMoveCommand, jobDeployCommands []*JobDeployCommand) (err error) {
	outputFile := fmt.Sprintf(r.OutputDir+"/%s", time.Now().Format("20060102_150405"))
	latestFile := fmt.Sprintf(r.OutputDir + "/latest.csv")

	buf := bytes.NewBufferString("")
	if instanceMoveCommands != nil {
		for _, v := range instanceMoveCommands {
			buf.WriteString(fmt.Sprintf("%d,inst_%d,machine_%d\n", v.Round, v.InstanceId, v.MachineId))
		}
	}
	if jobDeployCommands != nil {
		for _, v := range jobDeployCommands {
			buf.WriteString(fmt.Sprintf("%s,machine_%d,%d,%d\n", v.JobId, v.MachineId, v.Count, v.StartMinutes))
		}
	}

	err = ioutil.WriteFile(latestFile, buf.Bytes(), os.ModePerm)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(outputFile+".csv", buf.Bytes(), os.ModePerm)
	if err != nil {
		return err
	}

	summaryBuf := bytes.NewBufferString("")
	summaryBuf.WriteString(fmt.Sprintf("machineCount=%d\n", r.DeployedMachineCount))
	summaryBuf.WriteString(fmt.Sprintf("instanceMoveCommand=%d\n", len(instanceMoveCommands)))
	summaryBuf.WriteString(fmt.Sprintf("cost=%f,realCost=%f\n", r.CalcTotalScore(), r.CalcTotalScoreReal()))

	err = ioutil.WriteFile(fmt.Sprintf(outputFile+"_summary.csv"), summaryBuf.Bytes(), os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func (r *ResourceManagement) mergeOutput() (err error) {
	r.log("mergeOutput start\n")

	r.beginOffline()

	err = r.firstFitJobs()
	if err != nil {
		return err
	}

	instanceMoveCommands, err := NewOnlineMerge(r).Run()
	if err != nil {
		return err
	}

	jobDeployCommands := r.buildJobDeployCommands()

	err = NewReplay(r, instanceMoveCommands, jobDeployCommands).Run()
	if err != nil {
		return err
	}

	return r.output(instanceMoveCommands, jobDeployCommands)
}
