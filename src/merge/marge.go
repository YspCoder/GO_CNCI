/**
 * @Author: lipengfei
 * @Description:
 * @File:  marge
 * @Version: 1.0.0
 * @Date: 2021/08/06 16:39
 */
package merge

import (
	. "GO_CNCI/src/base"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func Merge(inputDir, score_path, detil_path string, number int) error {

	score, err := os.Create(score_path)
	if err != nil {
		Error("Create GO_CNCI_score Err : [%s]", err.Error())
		return err
	}

	detil, err := os.Create(detil_path)
	if err != nil {
		Error("Create GO_CNCI_detil Err : [%s]", err.Error())
		return err
	}
	for i := 1; i <= number; i++ {
		Score_string := fmt.Sprintf("%s/GO_CNCI_file_score%v", inputDir, i)
		f, err := os.OpenFile(Score_string, os.O_RDONLY, os.ModePerm)
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
		Detil_string := fmt.Sprintf("%s/GO_CNCI_file_detil%v", inputDir, i)
		f, err := os.OpenFile(Detil_string, os.O_RDONLY, os.ModePerm)
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
	SVM_arr_store := []string{}
	SVM_FILE_ONE, err := os.Create(FileName)
	if err != nil {
		Error("Create error![%v]\n", err.Error())
		return err
	}
	for i := 0; i < len(rec); i++ {
		temp_str := rec[i]
		temp_arr := strings.Split(temp_str, " ")
		for j := 0; j < len(temp_arr); j++ {
			index := j + 1
			temp_arr[j] = strconv.Itoa(index) + ":" + temp_arr[j]
		}
		str_temp := strings.Join(temp_arr, " ")
		SVM_arr_store = append(SVM_arr_store, str_temp)
		_, err = SVM_FILE_ONE.WriteString(str_temp + "\n")
		if err != nil {
			Error("WriteString error![%v]\n", err.Error())
			continue
		}
	}
	defer SVM_FILE_ONE.Close()
	return nil
}
