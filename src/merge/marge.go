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
	"GO_CNCI/src/utils"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func Merge(inputDir, outfile string, number int) ([]string, error) {
	scoreArray := make([]string, 0)
	detilArray := make([]string, 0)
	for i := 1; i <= number; i++ {
		Score_string := fmt.Sprintf("%s/GO_CNCI_file_score%v", inputDir, i)
		_score := utils.ReadFileArray(Score_string)
		scoreArray = append(scoreArray, _score...)
	}
	for i := 1; i <= number; i++ {
		Detil_string := fmt.Sprintf("%s/GO_CNCI_file_detil%v", inputDir, i)
		_detil := utils.ReadFileArray(Detil_string)
		detilArray = append(detilArray, _detil...)
	}
	AddSvmLabel(scoreArray, outfile)
	return detilArray, nil
}

func AddSvmLabel(rec []string, FileName string) {
	SVM_arr_store := []string{}
	SVM_FILE_ONE, err := os.Create(FileName)
	if err != nil {
		Error("Create error![%v]\n", err.Error())
		return
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
}
