package merge

import (
	. "GO_CNCI/src/base"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func Merge(inputDir, scorePath, detilPath string, number int) error {
	score, err := os.Create(scorePath)
	if err != nil {
		Error("Create GO_CNCI_score Err : [%s]", err.Error())
		return err
	}

	detil, err := os.Create(detilPath)
	if err != nil {
		Error("Create GO_CNCI_detil Err : [%s]", err.Error())
		return err
	}
	for i := 1; i <= number; i++ {
		scoreString := fmt.Sprintf("%s/GO_CNCI_file_score%v", inputDir, i)
		f, err := os.OpenFile(scoreString, os.O_RDONLY, os.ModePerm)
		if err != nil {
			Error("OpenFile GO_CNCI_file_score Err : [%s]", err.Error())
			return err
		}
		b, err := ioutil.ReadAll(f)
		if err != nil {
			Error("ReadAll GO_CNCI_file_score Err : [%s]", err.Error())
			return err
		}
		_, _ = score.Write(b)
		_ = f.Close()
	}
	for i := 1; i <= number; i++ {
		detilString := fmt.Sprintf("%s/GO_CNCI_file_detil%v", inputDir, i)
		f, err := os.OpenFile(detilString, os.O_RDONLY, os.ModePerm)
		if err != nil {
			Error("OpenFile GO_CNCI_file_detil Err : [%s]", err.Error())
			return err
		}
		b, err := ioutil.ReadAll(f)
		if err != nil {
			Error("ReadAll GO_CNCI_file_detil Err : [%s]", err.Error())
			return err
		}
		_, _ = detil.Write(b)
		_ = f.Close()
	}
	return nil
}

func AddSvmLabel(rec []string, FileName string) error {
	SvmArrScore := make([]string, 0)
	SvmFileOne, err := os.Create(FileName)
	if err != nil {
		Error("Create error![%v]\n", err.Error())
		return err
	}
	for i := 0; i < len(rec); i++ {
		tempStr := rec[i]
		tempArr := strings.Split(tempStr, " ")
		for j := 0; j < len(tempArr); j++ {
			index := j + 1
			tempArr[j] = strconv.Itoa(index) + ":" + tempArr[j]
		}
		strTemp := strings.Join(tempArr, " ")
		SvmArrScore = append(SvmArrScore, strTemp)
		_, err = SvmFileOne.WriteString(strTemp + "\n")
		if err != nil {
			Error("WriteString error![%v]\n", err.Error())
			continue
		}
	}
	defer SvmFileOne.Close()
	return nil
}
