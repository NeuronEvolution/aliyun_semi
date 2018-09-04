package schedule

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

func (r *ResourceManagement) output(
	machines []*Machine, instanceMoveCommands []*InstanceMoveCommand, jobDeployCommands []*JobDeployCommand) (err error) {
	//输出结果
	outputFile := fmt.Sprintf(r.OutputDir+"/%s", time.Now().Format("20060102_150405"))
	buf := bytes.NewBufferString("")
	if instanceMoveCommands != nil {
		for _, v := range instanceMoveCommands {
			buf.WriteString(fmt.Sprintf("%d,inst_%d,machine_%d\n", v.Round, v.InstanceId, v.MachineId))
		}
	}
	if jobDeployCommands != nil {
		for _, v := range jobDeployCommands {
			buf.WriteString(fmt.Sprintf("%s,machine_%d,%d,%d\n", v.JobId, v.MachineId, v.StartMinutes, v.Count))
		}
	}
	err = ioutil.WriteFile(outputFile+".csv", buf.Bytes(), os.ModePerm)
	if err != nil {
		return err
	}

	totalMachineCount := 0
	for _, m := range machines {
		if m.InstanceListCount > 0 || m.JobListCount > 0 {
			totalMachineCount++
		}
	}

	totalJobCount := 0
	for _, config := range r.JobConfigMap {
		if config != nil {
			totalJobCount += config.InstanceCount
		}
	}

	totalJobWithInstance := 0
	for _, m := range r.MachineList {
		if m.InstanceListCount > 0 {
			totalJobWithInstance += m.JobListCount
		}
	}

	//输出结果说明
	timeCost := time.Now().Sub(r.StartTime).Seconds()
	costReal := MachinesGetScoreReal(machines)
	summaryBuf := bytes.NewBufferString("")
	summaryBuf.WriteString(fmt.Sprintf("%f\n", costReal))
	summaryBuf.WriteString(fmt.Sprintf("timeCost=%f\n", timeCost))
	summaryBuf.WriteString(fmt.Sprintf("file=%s\n", outputFile))
	summaryBuf.WriteString(fmt.Sprintf("JobScheduleCpuLimitStep=%f\n", JobScheduleCpuLimitStep))
	summaryBuf.WriteString(fmt.Sprintf("instanceMachineCount=%d,totalMachineCount=%d\n",
		r.DeployedMachineCount, totalMachineCount))
	summaryBuf.WriteString(fmt.Sprintf("instanceMoveCommand=%d\n", len(instanceMoveCommands)))
	summaryBuf.WriteString(fmt.Sprintf("jobDeployCommands=%d,totalJobCount=%d,jobWithInstance=%d\n",
		len(jobDeployCommands), totalJobCount, totalJobWithInstance))
	summaryBuf.WriteString(fmt.Sprintf("cost=%f,realCost=%f\n", MachinesGetScore(machines), costReal))
	err = ioutil.WriteFile(outputFile+"_summary.csv", summaryBuf.Bytes(), os.ModePerm)
	if err != nil {
		return err
	}

	//更新最佳结果
	update := false
	bestSummaryFile := r.OutputDir + "/best_summary.csv"
	bestSummary, err := ioutil.ReadFile(bestSummaryFile)
	if err != nil {
		update = true
	} else {
		tokens := strings.Split(string(bestSummary), "\n")
		if len(tokens) == 0 {
			update = true
		} else {
			bestScore, err := strconv.ParseFloat(tokens[0], 64)
			if err != nil {
				update = true
			} else {
				if costReal < bestScore {
					update = true
				} else {
					update = false
				}
			}
		}
	}
	if update {
		r.log("output update best cost=%f\n", costReal)
		bestFile := r.OutputDir + "/best.csv"
		err = ioutil.WriteFile(bestFile, buf.Bytes(), os.ModePerm)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(bestSummaryFile, summaryBuf.Bytes(), os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}
