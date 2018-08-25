package semitemp

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

//todo 重构代码
func (r *ResourceManagement) MergeAndOutput() (err error) {
	merge := NewResourceManagement()
	err = merge.Init(
		r,
		r.machineResourcesConfig,
		r.appResourcesConfig,
		r.appInterferenceConfig,
		r.tempInstanceDeployConfig)
	if err != nil {
		return err
	}

	err = merge.MergeTo(r)
	if err != nil {
		return err
	}

	err = merge.Output(time.Duration(0))
	if err != nil {
		return err
	}

	playback := NewResourceManagement()
	err = playback.Init(
		r,
		r.machineResourcesConfig,
		r.appResourcesConfig,
		r.appInterferenceConfig,
		r.tempInstanceDeployConfig)
	if err != nil {
		return err
	}

	err = playback.Play(merge.DeployCommandHistory)
	if err != nil {
		return err
	}

	return nil
}

func (r *ResourceManagement) Output(duration time.Duration) (err error) {
	fmt.Printf("time=%f\n", duration.Seconds())

	a5506, err := ioutil.ReadFile("./_output/submit_20180710_031628.csv")
	if err != nil {
		return err
	}

	outputFile := fmt.Sprintf("_output/submit_%s", time.Now().Format("20060102_150405"))

	buf := bytes.NewBuffer(a5506)
	buf.WriteString("#\n")
	buf.Write(r.DeployCommandHistory.OutputCSV())

	err = ioutil.WriteFile(outputFile+".csv", buf.Bytes(), os.ModePerm)
	if err != nil {
		return err
	}

	summaryBuf := bytes.NewBufferString("")
	r.DebugStatus(summaryBuf)
	summaryBuf.WriteString(fmt.Sprintf("time=%f秒\n", duration.Seconds()))
	err = ioutil.WriteFile(fmt.Sprintf(outputFile+"_summary.csv"),
		summaryBuf.Bytes(), os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
