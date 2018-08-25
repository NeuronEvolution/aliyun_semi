package schedule

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

func (r *ResourceManagement) buildJobDeployCommands() (commands []*JobDeployCommand) {
	return nil
}

func (r *ResourceManagement) output(instanceMoveCommands []*InstanceMoveCommand, jobDeployCommands []*JobDeployCommand) (err error) {
	outputFile := fmt.Sprintf(r.OutputDir+"/%s", time.Now().Format("20060102_150405"))

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

	err = ioutil.WriteFile(outputFile+".csv", buf.Bytes(), os.ModePerm)
	if err != nil {
		return err
	}

	summaryBuf := bytes.NewBufferString("")
	summaryBuf.WriteString(fmt.Sprintf("machineCount=%d\n", r.DeployedMachineCount))
	summaryBuf.WriteString(fmt.Sprintf("instanceMoveCommand=%d\n", len(instanceMoveCommands)))
	summaryBuf.WriteString(fmt.Sprintf("cost=%f\n", r.CalcTotalScore()))

	err = ioutil.WriteFile(fmt.Sprintf(outputFile+"_summary.csv"), summaryBuf.Bytes(), os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func (r *ResourceManagement) mergeOutput() (err error) {
	r.log("mergeOutput start\n")

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
