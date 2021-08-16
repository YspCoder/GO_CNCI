package utils

import (
	. "GO_CNCI/src/base"
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// XRangeInt step-by-step
func XRangeInt(args ...int) chan int {
	if l := len(args); l < 1 || l > 3 {
		Error("Error args length, xRangeInt requires 1-3 int arguments")
	}
	var start, stop int
	var step = 1
	switch len(args) {
	case 1:
		stop = args[0]
		start = 0
	case 2:
		start, stop = args[0], args[1]
	case 3:
		start, stop, step = args[0], args[1], args[2]
	}
	ch := make(chan int)
	go func() {
		if step > 0 {
			for start < stop {
				ch <- start
				start = start + step
			}
		} else {
			for start > stop {
				ch <- start
				start = start + step
			}
		}
		close(ch)
	}()
	return ch
}

// ReadFileArray Convert each line of the file to an array parameter
func ReadFileArray(path string) []string {
	f, err := os.Open(path)
	if err != nil {
		Error("Read file fail : %v", err.Error())
		return nil
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	fileArray := make([]string, 0)
	for scanner.Scan() {
		fileArray = append(fileArray, scanner.Text())
	}
	Info("Read file success filename : [%s]", path)
	return fileArray
}

func TwoLineFasta(sequenceArr []string) []string {
	TmpSequenceArr := make([]string, 0)
	TmpTransStr := ""
	for i := 0; i < len(sequenceArr); i++ {
		if strings.Contains(sequenceArr[i], ">") {
			if i == 0 {
				TmpSequenceArr = append(TmpSequenceArr, sequenceArr[i])
			} else {
				TmpSequenceArr = append(TmpSequenceArr, TmpTransStr)
				TmpSequenceArr = append(TmpSequenceArr, sequenceArr[i])
				TmpTransStr = ""
			}
		} else {
			if i == len(sequenceArr)-1 {
				TmpTransStr = fmt.Sprintf("%v%v", TmpTransStr, sequenceArr[i])
				TmpSequenceArr = append(TmpSequenceArr, TmpTransStr)
			} else {
				TmpTransStr = fmt.Sprintf("%v%v", TmpTransStr, sequenceArr[i])
			}
		}
	}
	return TmpSequenceArr
}

func SplitFile(files []string, thread int) *sync.Map {
	fileNum := len(files) / 2
	splitStep := fileNum / thread
	splitStep = splitStep * 2
	start := 0
	end := splitStep
	in := sync.Map{}
	for i := 1; i <= thread; i++ {
		mp := make(map[string]string)
		for v := range XRangeInt(start, end, 2) {
			key := strings.ReplaceAll(files[v], "\r", "")
			v1 := strings.ReplaceAll(files[v+1], "\r", "")
			value := strings.ReplaceAll(v1, "u", "t")
			mp[key] = value
		}
		in.Store(i, mp)
		start += splitStep
		end += splitStep
	}
	return &in
}

func Libsvm(filepath, outSvm, outfile, outTmp, libsvmPath, CnciParameters, classModel string) error {

	var scale, model string
	if classModel == "ve" {
		scale = "/go_scale"
		model = "/go_model"
	} else if classModel == "pl" {
		scale = "/plant_scale"
		model = "/plant_model"
	}

	err := CmdBash("bash", "-c", libsvmPath+"/svm-scale -r "+CnciParameters+scale+" "+filepath+" > "+outSvm)
	if err != nil {
		Error("svm-scale err [%s]", err.Error())
		return err
	}
	err = CmdBash("bash", "-c", libsvmPath+"/svm-predict "+outSvm+" "+CnciParameters+model+" "+outfile+" > "+outTmp)
	if err != nil {
		Error("svm-predict err [%s]", err.Error())
		return err
	}
	return nil
}

func PutResult(detilArray []string, filepath string) []string {
	fileArr := ReadFileArray(filepath)
	classifyIndex := 0
	indexCoding := "1"
	TempResultArr := make([]string, 0)
	sort.Strings(detilArray)
	for i := 0; i < len(detilArray); i++ {
		tempLabelArrLabel := strings.Split(detilArray[i], ";;;;;")
		Label := tempLabelArrLabel[0]
		tempLabelArr := strings.Split(tempLabelArrLabel[1], " ")
		subTempLabelArr := tempLabelArr[1:]
		subTempLabelStr := strings.Join(subTempLabelArr, " ")
		if fileArr[classifyIndex] == indexCoding {
			Label = fmt.Sprintf("%s;;;;; coding", Label)
		} else {
			Label = fmt.Sprintf("%s;;;;; noncoding", Label)
		}
		classifyIndex = classifyIndex + 1
		TempResultStr := fmt.Sprintf("%s %s", Label, subTempLabelStr)
		TempResultArr = append(TempResultArr, TempResultStr)
	}
	return TempResultArr
}

func PrintResult(result []string, outDetil string) {
	OutFileResult, err := os.Create(outDetil)
	if err != nil {
		Error("PrintResult Err : [%s]", err.Error())
		return
	}
	Tabel := "TranscriptId" + "\t" + "index" + "\t" + "score" + "\t" + "start" + "\t" + "end" + "\t" + "length" + "\n"
	_, _ = OutFileResult.WriteString(Tabel)
	for _, v := range result {
		outLabelArr := strings.Split(v, ";;;;; ")
		labelArr := strings.Split(outLabelArr[1], " ")
		if len(labelArr) < 5 {
			continue
		}
		TLabel := outLabelArr[0]
		TableLabel := TLabel[1:]
		property := labelArr[0]
		startPosition := labelArr[1]
		stopPosition := labelArr[2]
		value := labelArr[3]
		v1, _ := strconv.ParseFloat(substring(value), 64)
		tlen := labelArr[4]
		if v1 == 0 {
			v1 = v1 + 0.001
		}
		tempOutStr := ""
		if property == "noncoding" {
			v3 := (0.64 * v1) * 0.64
			if v3 > 0 {
				if v3 > 1 {
					v4 := -1 / v3
					tempOutStr = fmt.Sprintf("%v\t%v\t%v\t%v\t%v\t%v\n", TableLabel, property, v4, startPosition, stopPosition, tlen)
				} else {
					v4 := -1 * v3
					tempOutStr = fmt.Sprintf("%v\t%v\t%v\t%v\t%v\t%v\n", TableLabel, property, v4, startPosition, stopPosition, tlen)
				}
			} else {
				tempOutStr = fmt.Sprintf("%v\t%v\t%v\t%v\t%v\t%v\n", TableLabel, property, v3, startPosition, stopPosition, tlen)
			}
		} else if property == "coding" {
			if v1 <= 0.0 {
				if v1 <= -1 {
					v3 := -1 / v1
					tempOutStr = fmt.Sprintf("%v\t%v\t%v\t%v\t%v\t%v\n", TableLabel, property, v3, startPosition, stopPosition, tlen)
				} else {
					v3 := -1 * v1
					tempOutStr = fmt.Sprintf("%v\t%v\t%v\t%v\t%v\t%v\n", TableLabel, property, v3, startPosition, stopPosition, tlen)
				}
			} else {
				tempOutStr = fmt.Sprintf("%v\t%v\t%v\t%v\t%v\t%v\n", TableLabel, property, v1, startPosition, stopPosition, tlen)
			}
		}
		_, _ = OutFileResult.WriteString(tempOutStr)
	}
	defer OutFileResult.Close()
}

func substring(param string) string {
	if len(param) <= 5 {
		return param
	} else {
		return param[:5]
	}
}

func CmdBash(commandName string, p1 string, p2 string) error {
	cmd := exec.Command(commandName, p1, p2)
	Info("cmd : %v", cmd)
	cmd.Stderr = cmd.Stdout
	err := cmd.Start()
	if err != nil {
		_ = cmd.Process.Kill()
		return err
	}
	err = cmd.Wait()
	if err != nil {
		Error("Wait Err : %v", err.Error())
		return err
	}
	return nil

}

func Reverse(params []string) []string {
	for i, j := 0, len(params)-1; i < j; i, j = i+1, j-1 {
		params[i], params[j] = params[j], params[i]
	}
	return params
}

func ReverseFloats64(params []float64) []float64 {
	for i, j := 0, len(params)-1; i < j; i, j = i+1, j-1 {
		params[i], params[j] = params[j], params[i]
	}
	return params
}

func StringToArray(params string) []string {
	paramsCharAr := []byte(params)
	var paramArray = make([]string, 0)
	for i := 0; i < len(paramsCharAr); i++ {
		paramArray = append(paramArray, string(paramsCharAr[i]))
	}
	return paramArray
}

func InitCodonSeq(num, length, step int, Arr []string) string {
	TempStrPar := ""
	for w := range XRangeInt(num, length, step) {
		index := w
		code1 := Arr[index]
		index += 1
		code2 := Arr[index]
		index += 1
		code3 := Arr[index]
		Temp := code1 + code2 + code3
		TempStrPar = TempStrPar + Temp + " "
	}
	return TempStrPar
}

func ReadFileMatrix(path string) map[string]string {
	matrix := make(map[string]string, 0)
	f, err := os.Open(path)
	if err != nil {
		Error("Read file fail : %s", err.Error())
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		params := strings.Split(line, "\t")
		matrix[params[0]] = params[1]
	}
	return matrix
}
