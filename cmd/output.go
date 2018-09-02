package main

import (
	"bytes"
	"fmt"
	"github.com/NeuronEvolution/aliyun_semi/schedule"
	"io/ioutil"
	"os"
	"time"
	"strings"
	"strconv"
)

func outputSummary(buf *bytes.Buffer, dataset string) (totalScore float64 ){
	buf.WriteString(dataset + "\n")
	summary, err := ioutil.ReadFile("./_output/" + dataset + "/best_summary.csv")
	if err == nil {
		buf.WriteString(string(summary))
		lines:= strings.Split(string(summary),"\n")
		if len(lines)>0{
			score,err:=strconv.ParseFloat(lines[0],64)
			if err!=nil{
				totalScore+=score
			}
		}
	} else {
		fmt.Println("outputSummary failed", dataset, err)
	}
	buf.WriteString("\n")

	return totalScore
}

func output() (err error) {
	fmt.Println("output")

	a, err := ioutil.ReadFile("./_output/a/best.csv")
	if err != nil {
		return err
	}

	b, err := ioutil.ReadFile("./_output/b/best.csv")
	if err != nil {
		return err
	}

	c, err := ioutil.ReadFile("./_output/c/best.csv")
	if err != nil {
		return err
	}

	d, err := ioutil.ReadFile("./_output/d/best.csv")
	if err != nil {
		return err
	}

	e, err := ioutil.ReadFile("./_output/e/best.csv")
	if err != nil {
		return err
	}

	buf := bytes.NewBufferString("")
	buf.WriteString(string(a))
	buf.WriteString("#\n")
	buf.WriteString(string(b))
	buf.WriteString("#\n")
	buf.WriteString(string(c))
	buf.WriteString("#\n")
	buf.WriteString(string(d))
	buf.WriteString("#\n")
	buf.WriteString(string(e))

	err = schedule.MakeDirIfNotExists("./_output/submit/")
	if err != nil {
		return err
	}

	outputFile := fmt.Sprintf("./_output/submit/submit_%s", time.Now().Format("20060102_150405"))
	err = ioutil.WriteFile(outputFile+".csv", buf.Bytes(), os.ModePerm)
	if err != nil {
		return err
	}

	summaryBuf := bytes.NewBufferString("")
	outputSummary(summaryBuf, "a")
	outputSummary(summaryBuf, "b")
	outputSummary(summaryBuf, "c")
	outputSummary(summaryBuf, "d")
	outputSummary(summaryBuf, "e")
	summaryFile := outputFile + "_summary.csv"
	err = ioutil.WriteFile(summaryFile, summaryBuf.Bytes(), os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
