package merge

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func AddSvmLabel(rec []string, FileName string) error {
	SvmArrScore := make([]string, 0)
	SvmFileOne, err := os.Create(FileName)
	if err != nil {
		fmt.Printf("Create error![%v]\n", err.Error())
		return err
	}
	sort.Strings(rec)
	for i := 0; i < len(rec); i++ {
		tempArr := strings.Split(rec[i], " ")
		tempArr = tempArr[1:]
		for j := 0; j < len(tempArr); j++ {
			index := j + 1
			tempArr[j] = strconv.Itoa(index) + ":" + tempArr[j]
		}
		strTemp := strings.Join(tempArr, " ")
		SvmArrScore = append(SvmArrScore, strTemp)
		_, err = SvmFileOne.WriteString(strTemp + "\n")
		if err != nil {
			fmt.Printf("WriteString error![%v]\n", err.Error())
			continue
		}
	}
	defer SvmFileOne.Close()
	return nil
}
