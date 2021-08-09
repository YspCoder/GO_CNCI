/**
 * @Author: lipengfei
 * @Description:
 * @File:  reckon
 * @Version: 1.0.0
 * @Date: 2021/08/06 14:51
 */
package reckon

import (
	. "GO_CNCI/src/base"
	. "GO_CNCI/src/utils"
	"fmt"
	"github.com/EDDYCJY/gsema"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
)

var (
	MaxValue         = make([]float64, 0)
	MaxString        = make([]string, 0)
	LengthStoreArray = make([]int, 0)
	Pos              = make([]string, 0)
	OtherCdsArray    = make([]float64, 0)
	sema             = gsema.NewSemaphore(12)
)

type Reckon struct {
	TempScore           string
	TempDetil           string
	TempInput           string
	SequenceArr         []string
	Label               string
	Seq                 string
	SeqLen              int
	Rounds              int
	SequenceProcessArr  []string
	SequenceProcessArrR []string
	MLCDS_sequence      []string
	MLCDS_sequenceR     []string
	MLCDS_seq_length    int
	HashMatrix          *sync.Map
}

func New() *Reckon {
	return &Reckon{
		TempScore:           "",
		TempDetil:           "",
		TempInput:           "",
		SequenceArr:         nil,
		Label:               "",
		Seq:                 "",
		SeqLen:              0,
		Rounds:              0,
		SequenceProcessArr:  nil,
		SequenceProcessArrR: nil,
		MLCDS_sequence:      nil,
		MLCDS_sequenceR:     nil,
		MLCDS_seq_length:    0,
		HashMatrix:          nil,
	}
}

func (this *Reckon) Init(ws *sync.WaitGroup) {
	defer ws.Done()
	var err error
	score, err := os.Create(this.TempScore)
	if err != nil {
		Error("Create tmp error![%v]\n", err)
		return
	}
	_ = score.Close()
	detil, err := os.Create(this.TempDetil)
	if err != nil {
		Error("Create tmp error![%v]\n", err)
		return
	}
	_ = detil.Close()
	sequence_Arr := ReadFileArray(this.TempInput)
	sLen := len(sequence_Arr) - 1
	sequence_Arr = sequence_Arr[:sLen]
	label_Arr_tmp := make([]string, 0)
	fast_seq_Arr_tmp := make([]string, 0)
	for n := 0; n < len(sequence_Arr); n++ {
		if n == 0 || n%2 == 0 {
			label_Arr_tmp = append(label_Arr_tmp, sequence_Arr[n])
		} else {
			fast_seq_Arr_tmp = append(fast_seq_Arr_tmp, sequence_Arr[n])
		}
	}
	for d := 0; d < len(label_Arr_tmp); d++ {
		sema.Add(1)
		go this.Compare()
	}
	sema.Wait()

}

func (this *Reckon) Compare() {
	defer sema.Done()
	DetilLen := len(this.Seq)
	tran_fir_seq := strings.ToLower(this.Seq)
	tran_sec_seq := strings.ReplaceAll(tran_fir_seq, "u", "t")
	this.SequenceProcessArr = StringToArray(tran_sec_seq)
	this.SequenceProcessArrR = append(this.SequenceProcessArrR, this.SequenceProcessArr...)
	this.SequenceProcessArrR = Reverse(this.SequenceProcessArrR)
	this.SeqLen = len(this.SequenceProcessArr) - 1

	var wgs sync.WaitGroup
	for o1 := 0; o1 < 6; o1++ {
		wgs.Add(1)
		this.Rounds = o1
		go this.multilayerComparison(&wgs)
	}
	wgs.Wait()
	r_max_Value := MaxValue[:]
	sort.Float64s(r_max_Value)
	r_max_Value = ReverseFloats(r_max_Value)
	M := r_max_Value[0]
	orf_index := 0
	for o := 0; o < len(MaxValue); o++ {
		temp := MaxValue[o]
		if temp == M {
			orf_index = o
		}
	}
	o_tmp_arr := strings.Split(MaxString[orf_index], " ")
	var o_arr = make([]string, 0)
	o_arr = o_tmp_arr[:len(o_tmp_arr)-1]
	SequenceLen := len(o_arr) - 1
	M_score := 0.0
	for j := 0; j < SequenceLen; j++ {
		temp_trip := o_arr[j] + o_arr[j+1]
		byteMatchResult, _ := regexp.Match(`[atcg]{6}`, []byte(temp_trip))
		if byteMatchResult {
			if v9, ok := this.HashMatrix.Load(temp_trip); ok {
				v8 := v9.(float64)
				M_score = M_score + v8
			} else {
				continue
			}
		}
	}
	SequenceLen = SequenceLen + 2
	M_score = M_score / float64(SequenceLen)
	MLCDS_str := strings.Join(o_arr, "")
	this.MLCDS_sequence = StringToArray(MLCDS_str)
	this.MLCDS_sequenceR = append(this.MLCDS_sequenceR, this.MLCDS_sequence...)
	this.MLCDS_sequenceR = Reverse(this.MLCDS_sequenceR)
	this.MLCDS_seq_length = len(this.MLCDS_sequence) - 1
	for o1 := 1; o1 < 6; o1++ {
		wgs.Add(1)
		this.Rounds = o1
		go this.multilayerComparisonTwo(&wgs)
	}
	wgs.Wait()
	score_distance := 0.0
	for m := 0; m < len(OtherCdsArray); m++ {
		score_distance += M_score - OtherCdsArray[m]
	}
	score_distance = score_distance / 5
	if len(Pos) <= orf_index {
		Warn("Insufficient length")
		return
	}
	out_pos := Pos[orf_index]
	M_length := LengthStoreArray[orf_index]
	length_total_score := 0.0
	for p := 0; p < len(LengthStoreArray); p++ {
		length_total_score = length_total_score + float64(LengthStoreArray[p])
	}
	length_precent := float64(M_length) / length_total_score
	Coding_Array_one := make([]string, 0)
	for n := 0; n < len(o_arr); n++ {
		temp1 := o_arr[n]
		byteMatchResult, _ := regexp.Match(`[atcg{3}]`, []byte(temp1))
		if byteMatchResult && temp1 != "taa" && temp1 != "tag" && temp1 != "tga" {
			if v, ok := this.HashMatrix.Load(temp1); ok {
				v1 := v.(int)
				v2 := v1 + 1
				this.HashMatrix.LoadOrStore(temp1, v2)
			}
		}
	}

	this.HashMatrix.Range(func(key, value interface{}) bool {
		v := value.(int)
		Coding_Array_one = append(Coding_Array_one, strconv.Itoa(v))
		return true
	})
	C_num1 := 0.0
	for i := 0; i < len(Coding_Array_one); i++ {
		v1, _ := strconv.ParseFloat(Coding_Array_one[i], 64)
		C_num1 = C_num1 + v1
	}
	if C_num1 == 0.0 {
		C_num1 = 1
	}
	for i := 0; i < len(Coding_Array_one); i++ {
		v1, _ := strconv.ParseFloat(Coding_Array_one[i], 64)
		v2 := v1 / C_num1
		Coding_Array_one[i] = fmt.Sprintf("%v", v2)
	}
	Array_Str := strings.Join(Coding_Array_one, " ")
	PROPERTY_STR := fmt.Sprintf("%v %v %v %v %v %v\n", M, M_length, M_score, length_precent, score_distance, Array_Str)
	DETIL_STR := fmt.Sprintf("%v;;;;;%v %v %v\n", this.Label, out_pos, M_score, DetilLen)
	score, err := os.OpenFile(this.TempScore, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		Error("OpenFile : [%s] Err : [%s]", this.TempScore, err.Error())
		return
	}
	_, err = score.WriteString(PROPERTY_STR)
	if err != nil {
		Error("WriteString Err : [%s]", err.Error())
	}
	detil, err := os.OpenFile(this.TempDetil, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		Error("OpenFile : [%s] Err : [%s]", this.TempDetil, err.Error())
		return
	}
	_, err = detil.WriteString(DETIL_STR)
	if err != nil {
		Error("WriteString Err : [%s]", err.Error())
	}
}

func (this *Reckon) multilayerComparison(wg *sync.WaitGroup) {
	defer wg.Done()
	CodonScore := make([]float64, 0)
	TempStr := ""
	if this.Rounds < 3 {
		TempStr = InitCodonSeq(this.Rounds, this.SeqLen-1, 3, this.SequenceProcessArr)
	}
	if 2 < this.Rounds && this.Rounds < 6 {
		TempStr = InitCodonSeq(this.Rounds-3, this.SeqLen-1, 3, this.SequenceProcessArrR)
	}
	TempArray := strings.Split(TempStr, " ")
	TempArray = TempArray[:len(TempArray)-1]
	seqLength := len(TempArray)
	WindowStep := 50
	WinLen := seqLength - WindowStep
	if seqLength > WindowStep {
		for EachCodon := 0; EachCodon < WinLen; EachCodon++ {
			num := 0.0
			SingleArray := make([]string, 0)
			for t := range XRangeInt(EachCodon, WindowStep+EachCodon) {
				SingleArray = append(SingleArray, TempArray[t])
			}
			SinLen := len(SingleArray) - 1
			for n := 0; n < SinLen; n++ {
				temp1 := SingleArray[n] + SingleArray[n+1]
				byteMatchResult, _ := regexp.Match(`[atcg]{6}`, []byte(temp1))
				if byteMatchResult {
					if v9, ok := this.HashMatrix.Load(temp1); ok {
						v8 := v9.(float64)
						num = num + v8
					} else {
						continue
					}
				}
			}
			num = num / 50
			CodonScore = append(CodonScore, num)
		}
		Start := 0
		End := 0
		Max := 0.0
		for r := 0; r < len(CodonScore); r++ {
			sum := 0.0
			CodonLength := len(CodonScore)
			for e := range XRangeInt(r, CodonLength) {
				sum = sum + CodonScore[e]
				if sum > Max {
					Start = r
					End = e
					Max = sum
				}
			}
		}
		OutStr := ""
		for out := Start; out < End+1; out++ {
			OutStr = OutStr + TempArray[out] + " "
		}
		Start = Start * 3
		End = End * 3
		Position := fmt.Sprintf("%v %v", Start, End)
		Pos = append(Pos, Position)
		MaxValue = append(MaxValue, Max)
		MaxString = append(MaxString, OutStr)
		OutParray := strings.Split(OutStr, " ")
		max_length := len(OutParray) - 1
		Onum := 0.0
		for n := 0; n < max_length; n++ {
			temp1 := OutParray[n] + OutParray[n+1]
			byteMatchResult, _ := regexp.Match(`[atcg]{6}`, []byte(temp1))
			if byteMatchResult {
				if v9, ok := this.HashMatrix.Load(temp1); ok {
					v8 := v9.(float64)
					Onum = Onum + v8
				} else {
					continue
				}
			}
		}
		LengthStoreArray = append(LengthStoreArray, max_length)
	} else {
		num := 0.0
		for n := 0; n < seqLength-1; n++ {
			temp1 := TempArray[n] + TempArray[n+1]
			byteMatchResult, _ := regexp.Match(`[atcg]{6}`, []byte(temp1))
			if byteMatchResult {
				if v9, ok := this.HashMatrix.Load(temp1); ok {
					v8 := v9.(float64)
					//v1, _ := strconv.ParseFloat(v8, 64)
					num = num + v8
				} else {
					continue
				}
			}
		}
		OutStr := strings.Join(TempArray, " ")
		Pos = append(Pos, "Full Length")
		MaxValue = append(MaxValue, num)
		MaxString = append(MaxString, OutStr)
		LengthStoreArray = append(LengthStoreArray, seqLength)
	}
}

func (this *Reckon) multilayerComparisonTwo(wg *sync.WaitGroup) {
	defer wg.Done()
	MLCDS_TempStr := ""
	if this.Rounds < 3 {
		MLCDS_TempStr = InitCodonSeq(this.Rounds, this.MLCDS_seq_length-1, 3, this.MLCDS_sequence)
	}
	if 2 < this.Rounds && this.Rounds < 6 {
		MLCDS_TempStr = InitCodonSeq(this.Rounds, this.MLCDS_seq_length-1, 3, this.MLCDS_sequenceR)
	}
	MLCDS_array := strings.Split(MLCDS_TempStr, " ")
	MLCDS_array = MLCDS_array[:len(MLCDS_array)-1]
	other_num := 0.0
	MLCDS_array_Len := len(MLCDS_array) - 1
	for j := 0; j < MLCDS_array_Len; j++ {
		temp2 := MLCDS_array[j] + MLCDS_array[j+1]
		byteMatchResult, _ := regexp.Match(`[atcg]{6}`, []byte(temp2))
		if byteMatchResult {
			if v9, ok := this.HashMatrix.Load(temp2); ok {
				v8 := v9.(float64)
				other_num = other_num + v8
			} else {
				continue
			}
		}
	}
	MLCDS_array_Len = MLCDS_array_Len + 2
	other_num = other_num / float64(MLCDS_array_Len)
	OtherCdsArray = append(OtherCdsArray, other_num)
}
