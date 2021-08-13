/**
 * @Author: lipengfei
 * @Description:
 * @File:  common
 * @Version: 1.0.0
 * @Date: 2021/08/06 14:20
 */
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

var (
	Alphabet = []string{"ttt", "ttc", "tta", "ttg", "tct", "tcc", "tca", "tcg", "tat", "tac", "tgt", "tgc", "tgg", "ctt", "ctc", "cta", "ctg", "cct", "ccc", "cca", "ccg", "cat", "cac", "caa", "cag", "cgt", "cgc", "cga", "cgg", "att", "atc", "ata", "atg", "act", "acc", "aca", "acg", "aat", "aac", "aaa", "aag", "agt", "agc", "aga", "agg", "gtt", "gtc", "gta", "gtg", "gct", "gcc", "gca", "gcg", "gat", "gac", "gaa", "gag", "ggt", "ggc", "gga", "ggg"}
)

//step-by-step
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

//Convert each line of the file to an array parameter
func ReadFileArray(path string) []string {
	f, err := os.Open(path)
	if err != nil {
		Error("Read file fail : %v", err.Error())
		return nil
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	fileArray := make([]string, 0)
	if strings.Contains(path, "GO_CNCI_file_score6") {
		Warn("*******************************************************")
		for scanner.Scan() {
			Warn(scanner.Text())
			fileArray = append(fileArray, scanner.Text())
		}
	} else {
		for scanner.Scan() {
			fileArray = append(fileArray, scanner.Text())
		}
	}
	Info("Read file success filename : [%s]", path)
	return fileArray
}

func TwoLineFasta(sequence_Arr []string) []string {
	Tmp_sequence_Arr := make([]string, 0)
	Tmp_trans_str := ""
	for i := 0; i < len(sequence_Arr); i++ {
		if strings.Contains(sequence_Arr[i], ">") {
			if i == 0 {
				Tmp_sequence_Arr = append(Tmp_sequence_Arr, sequence_Arr[i])
			} else {
				Tmp_sequence_Arr = append(Tmp_sequence_Arr, Tmp_trans_str)
				Tmp_sequence_Arr = append(Tmp_sequence_Arr, sequence_Arr[i])
				Tmp_trans_str = ""
			}
		} else {
			if i == len(sequence_Arr)-1 {
				Tmp_trans_str = fmt.Sprintf("%v%v", Tmp_trans_str, sequence_Arr[i])
				Tmp_sequence_Arr = append(Tmp_sequence_Arr, Tmp_trans_str)
			} else {
				Tmp_trans_str = fmt.Sprintf("%v%v", Tmp_trans_str, sequence_Arr[i])
			}
		}
	}
	return Tmp_sequence_Arr
}

func GetLabelArray(labelArray, fastaSeqArray []string) []string {
	TOT_STRING := make([]string, 0)
	for i := 0; i < len(labelArray); i++ {
		tmp_label := strings.ReplaceAll(labelArray[i], "\r", "")
		Temp_Seq := strings.ReplaceAll(fastaSeqArray[i], "\r", "")
		TOT_STRING = append(TOT_STRING, tmp_label)
		TOT_STRING = append(TOT_STRING, Temp_Seq)
	}
	return TOT_STRING
}

func SplitFile(files []string, thread int) *sync.Map {
	file_num := len(files) / 2
	split_step := file_num / thread
	split_step = split_step * 2
	start := 0
	end := split_step
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
		start += split_step
		end += split_step
	}
	return &in
}

func Libsvm(filepath, outSvm, outfile, outTmp, libsvm_path, CNCI_Parameters string) error {
	err := CmdBash("bash", "-c", libsvm_path+"/svm-scale -r "+CNCI_Parameters+"/python_scale "+filepath+" > "+outSvm)
	if err != nil {
		Error("svm-scale err [%s]", err.Error())
		return err
	}
	err = CmdBash("bash", "-c", libsvm_path+"/svm-predict "+outSvm+" "+CNCI_Parameters+"/python_model "+outfile+" > "+outTmp)
	if err != nil {
		Error("svm-predict err [%s]", err.Error())
		return err
	}
	return nil
}

func PutResult(detil_array []string, filepath string) []string {
	file_Arr := ReadFileArray(filepath)
	classify_index := 0
	index_coding := "1"
	Temp_Result_Arr := make([]string, 0)
	sort.Strings(detil_array)
	for _, v := range detil_array {
		temp_label_arr_label := strings.Split(v, ";;;;;")
		Label := temp_label_arr_label[0]
		temp_label_arr := strings.Split(temp_label_arr_label[1], " ")
		sub_temp_label_arr := temp_label_arr[1:]
		sub_temp_label_str := strings.Join(sub_temp_label_arr, " ")
		if file_Arr[classify_index] == index_coding {
			Label = fmt.Sprintf("%s;;;;; coding", Label)
		} else {
			Label = fmt.Sprintf("%s;;;;; noncoding", Label)
		}
		classify_index = classify_index + 1
		Temp_Result_str := fmt.Sprintf("%s %s", Label, sub_temp_label_str)
		Temp_Result_Arr = append(Temp_Result_Arr, Temp_Result_str)
	}
	return Temp_Result_Arr
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
		T_label := outLabelArr[0]
		Tabel_label := T_label[1:]
		property := labelArr[0]
		start_position := labelArr[1]
		stop_position := labelArr[2]
		value := labelArr[3]
		v1, _ := strconv.ParseFloat(substring(value), 32)
		tlen := labelArr[4]
		if v1 == 0 {
			v1 = v1 + 0.001
		}
		temp_out_str := ""
		if property == "noncoding" {
			v3 := (0.64 * v1) * 0.64
			if v3 > 0 {
				if v3 > 1 {
					v4 := -1 / v3
					temp_out_str = fmt.Sprintf("%v\t%v\t%.8f\t%v\t%v\t%v\n", Tabel_label, property, v4, start_position, stop_position, tlen)
				} else {
					v4 := -1 * v3
					temp_out_str = fmt.Sprintf("%v\t%v\t%.8f\t%v\t%v\t%v\n", Tabel_label, property, v4, start_position, stop_position, tlen)
				}
			} else {
				temp_out_str = fmt.Sprintf("%v\t%v\t%.8f\t%v\t%v\t%v\n", Tabel_label, property, v3, start_position, stop_position, tlen)
			}
		} else if property == "coding" {
			if v1 <= 0 {
				if v1 <= -1 {
					v3 := -1 / v1
					temp_out_str = fmt.Sprintf("%v\t%v\t%.8f\t%v\t%v\t%v\n", Tabel_label, property, v3, start_position, stop_position, tlen)
				} else {
					v3 := -1 * v1
					temp_out_str = fmt.Sprintf("%v\t%v\t%.8f\t%v\t%v\t%v\n", Tabel_label, property, v3, start_position, stop_position, tlen)
				}
			} else {
				temp_out_str = fmt.Sprintf("%v\t%v\t%v\t%.8f\t%v\t%v\n", Tabel_label, property, v1, start_position, stop_position, tlen)
			}
		}
		_, _ = OutFileResult.WriteString(temp_out_str)
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

func ReverseFloats32(params []float32) []float32 {
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

func Tran_checkSeq(input_arr []string) ([]string, []string) {
	label_Arr := make([]string, 0)
	FastA_seq_Arr := make([]string, 0)
	for n := 0; n < len(input_arr); n++ {
		if n == 0 || n%2 == 0 {
			label_Arr = append(label_Arr, input_arr[n])
		} else {
			FastA_seq_Arr = append(FastA_seq_Arr, input_arr[n])
		}
	}
	num := 0
	for i := 0; i < len(label_Arr); i++ {
		Seq := FastA_seq_Arr[num]
		tran_fir_seq := strings.ToLower(Seq)
		tran_sec_seq_one := strings.ReplaceAll(tran_fir_seq, "u", "t")
		tran_sec_seq := strings.ReplaceAll(tran_sec_seq_one, "\r", "")
		if strings.Contains(tran_sec_seq, "n") {
			label_Arr = append(label_Arr[:num], label_Arr[num+1:]...)
			FastA_seq_Arr = append(FastA_seq_Arr[:num], FastA_seq_Arr[num+1:]...)
			continue
		}
		if strings.Contains(tran_sec_seq, "n") {
			label_Arr = append(label_Arr[:num], label_Arr[num+1:]...)
			FastA_seq_Arr = append(FastA_seq_Arr[:num], FastA_seq_Arr[num+1:]...)
			continue
		}
		if strings.Contains(tran_sec_seq, "w") {
			label_Arr = append(label_Arr[:num], label_Arr[num+1:]...)
			FastA_seq_Arr = append(FastA_seq_Arr[:num], FastA_seq_Arr[num+1:]...)
			continue
		}
		if strings.Contains(tran_sec_seq, "d") {
			label_Arr = append(label_Arr[:num], label_Arr[num+1:]...)
			FastA_seq_Arr = append(FastA_seq_Arr[:num], FastA_seq_Arr[num+1:]...)
			continue
		}
		if strings.Contains(tran_sec_seq, "r") {
			label_Arr = append(label_Arr[:num], label_Arr[num+1:]...)
			FastA_seq_Arr = append(FastA_seq_Arr[:num], FastA_seq_Arr[num+1:]...)
			continue
		}
		if strings.Contains(tran_sec_seq, "s") {
			label_Arr = append(label_Arr[:num], label_Arr[num+1:]...)
			FastA_seq_Arr = append(FastA_seq_Arr[:num], FastA_seq_Arr[num+1:]...)
			continue
		}
		if strings.Contains(tran_sec_seq, "y") {
			label_Arr = append(label_Arr[:num], label_Arr[num+1:]...)
			FastA_seq_Arr = append(FastA_seq_Arr[:num], FastA_seq_Arr[num+1:]...)
			continue
		}
		if strings.Contains(tran_sec_seq, "m") {
			label_Arr = append(label_Arr[:num], label_Arr[num+1:]...)
			FastA_seq_Arr = append(FastA_seq_Arr[:num], FastA_seq_Arr[num+1:]...)
			continue
		}
		num = num + 1
	}
	return label_Arr, FastA_seq_Arr
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

//func ReadFileMatrix(path string) map[string]string {
//	matrix := make(map[string]string, 0)
//	f, err := os.Open(path)
//	if err != nil {
//		Error("Read file fail : %s", err.Error())
//	}
//	defer f.Close()
//	scanner := bufio.NewScanner(f)
//	for scanner.Scan() {
//		line := scanner.Text()
//		params := strings.Split(line, "\t")
//		idx := strings.LastIndex(params[1], ".")
//		v1 := params[1][idx+1:]
//		v2 := params[1][:idx+1]
//		v3 := v1[:len(v1)-5]
//		v4 := fmt.Sprintf("%v%v", v2, v3)
//		matrix[params[0]] = v4
//	}
//	return matrix
//}

//init sync.Map
func GetAlphabetMap() *sync.Map {
	var ab = &sync.Map{}
	for _, v := range Alphabet {
		ab.Store(v, 0)
	}
	return ab
}
