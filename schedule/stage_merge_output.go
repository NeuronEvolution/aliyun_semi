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

func (r *ResourceManagement) buildJobDeployCommands(machines []*Machine) (commands []*JobDeployCommand) {
	commands = make([]*JobDeployCommand, 0)
	for _, m := range machines {
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

	//输出结果说明
	costReal := r.CalcTotalScoreReal()
	summaryBuf := bytes.NewBufferString("")
	summaryBuf.WriteString(fmt.Sprintf("%f\n", costReal))
	summaryBuf.WriteString(fmt.Sprintf("file=%s\n", outputFile))
	summaryBuf.WriteString(fmt.Sprintf("machineCount=%d\n", r.DeployedMachineCount))
	summaryBuf.WriteString(fmt.Sprintf("instanceMoveCommand=%d\n", len(instanceMoveCommands)))
	summaryBuf.WriteString(fmt.Sprintf("cost=%f,realCost=%f\n", r.CalcTotalScore(), costReal))
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

func (r *ResourceManagement) mergeOutput() (err error) {
	r.log("mergeOutput start\n")

	r.beginOffline()

	//todo 这里需要考虑在线迁移时的实例交换,改为从初始状态迁移后再部署任务
	machines := MachinesClone(r.MachineList)

	for _, m := range machines {
		m.beginOffline()
	}

	totalScore := float64(0)
	for _, m := range machines {
		totalScore += m.GetCpuCost()
	}
	r.log("mergeOutput init totalScore=%f\n", totalScore)

	err = r.firstFitJobs(machines)
	if err != nil {
		return err
	}

	totalScore = float64(0)
	for _, m := range machines {
		totalScore += m.GetCpuCost()
	}
	r.log("mergeOutput firstFitJobs totalScore=%f\n", totalScore)

	instanceMoveCommands, err := NewOnlineMerge(r).Run()
	if err != nil {
		return err
	}

	jobDeployCommands := r.buildJobDeployCommands(machines)

	err = NewReplay(r, instanceMoveCommands, jobDeployCommands).Run(machines)
	if err != nil {
		return nil
	}

	return r.output(instanceMoveCommands, jobDeployCommands)
}
