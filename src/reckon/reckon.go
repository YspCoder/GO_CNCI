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
)

type Reckon struct {
	TempScore  string
	TempDetil  string
	TempInput  string
	HashMatrix map[string]string
	Thread     int
}

func New() *Reckon {
	return &Reckon{
		TempScore:  "",
		TempDetil:  "",
		TempInput:  "",
		HashMatrix: nil,
		Thread:     0,
	}
}

func (this *Reckon) Init(wg *gsema.Semaphore) {
	defer wg.Done()
	var err error
	score, err := os.Create(this.TempScore)
	if err != nil {
		Error("Create tmp error![%v]\n", err.Error())
		return
	}
	_ = score.Close()
	detil, err := os.Create(this.TempDetil)
	if err != nil {
		Error("Create tmp error![%v]\n", err.Error())
		return
	}
	_ = detil.Close()
	sequenceArr := ReadFileArray(this.TempInput)
	labelArrTmp := make([]string, 0)
	fastqSeqArrTmp := make([]string, 0)
	for n := 0; n < len(sequenceArr); n++ {
		if n == 0 || n%2 == 0 {
			labelArrTmp = append(labelArrTmp, sequenceArr[n])
		} else {
			fastqSeqArrTmp = append(fastqSeqArrTmp, sequenceArr[n])
		}
	}
	var sm = gsema.NewSemaphore(this.Thread)
	HashMatrix := this.HashMatrix
	for d := 0; d < len(labelArrTmp); d++ {
		if d >= len(labelArrTmp) {
			continue
		}
		if d >= len(fastqSeqArrTmp) {
			continue
		}
		sm.Add(1)
		go func(sa *gsema.Semaphore, Seq, Label, TempScore, TempDetil string, thread int) {
			defer sa.Done()
			MaxValue := make([]float32, 0)
			MaxString := make([]string, 0)
			LengthStoreArray := make([]int, 0)
			Pos := make([]string, 0)
			OtherCdsArray := make([]float32, 0)
			DetilLen := len(Seq)
			tran_fir_seq := strings.ToLower(Seq)
			tran_sec_seq := strings.ReplaceAll(tran_fir_seq, "u", "t")
			sequenceProcessArr := StringToArray(tran_sec_seq)
			sequenceProcessArrR := Reverse(sequenceProcessArr)
			slen := len(sequenceProcessArr) - 1
			for o := 0; o < 6; o++ {
				CodonScore := make([]float32, 0)
				TempStr := ""
				if o < 3 {
					TempStr = InitCodonSeq(o, slen-1, 3, sequenceProcessArr)
				}
				if 2 < o && o < 6 {
					TempStr = InitCodonSeq(o-3, slen-1, 3, sequenceProcessArrR)
				}
				TempArray := strings.Split(TempStr, " ")
				TempArray = TempArray[:len(TempArray)-1]
				seqLength := len(TempArray)
				WindowStep := 50
				WinLen := seqLength - WindowStep
				if seqLength > WindowStep {
					for EachCodon := 0; EachCodon < WinLen; EachCodon++ {
						num := float32(0)
						SingleArray := make([]string, 0)
						for t := EachCodon; t < WindowStep+EachCodon; t++ {
							SingleArray = append(SingleArray, TempArray[t])
						}
						SinLen := len(SingleArray) - 1
						for n := 0; n < SinLen; n++ {
							temp1 := SingleArray[n] + SingleArray[n+1]
							byteMatchResult, _ := regexp.Match(`[atcg]{6}`, []byte(temp1))
							if byteMatchResult {
								v1, _ := strconv.ParseFloat(HashMatrix[temp1], 64)
								num = num + float32(v1)
							}
						}
						num = num / float32(WindowStep)
						CodonScore = append(CodonScore, num)
					}

					Start := 0
					End := 0
					Max := float32(0)
					for r := 0; r < len(CodonScore); r++ {
						sum := float32(0)
						CodonLength := len(CodonScore)
						for e := r; e < CodonLength; e++ {
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
					LengthStoreArray = append(LengthStoreArray, max_length)
				} else {
					num := float32(0)
					for n := 0; n < seqLength-1; n++ {
						temp1 := TempArray[n] + TempArray[n+1]
						byteMatchResult, _ := regexp.Match(`[atcg]{6}`, []byte(temp1))
						if byteMatchResult {
							if v9, ok := HashMatrix[temp1]; ok {
								v1, _ := strconv.ParseFloat(v9, 64)
								num = num + float32(v1)
							} else {
								num = num + 0
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
			r_max_Value := MaxValue[:]

			sort.Sort(Float32List(r_max_Value))
			r_max_Value = ReverseFloats32(r_max_Value)
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
			mScore := float32(0)
			for j := 0; j < SequenceLen; j++ {
				temp_trip := o_arr[j] + o_arr[j+1]
				byteMatchResult, _ := regexp.Match(`[atcg]{6}`, []byte(temp_trip))
				if byteMatchResult {
					if v9, ok := HashMatrix[temp_trip]; ok {
						v1, _ := strconv.ParseFloat(v9, 64)
						mScore = mScore + float32(v1)
					} else {
						mScore = mScore + 0
					}
				}
			}
			SequenceLen = SequenceLen + 2
			mScore = mScore / float32(SequenceLen)
			mlcds := strings.Join(o_arr, " ")
			MlcdsSequence := StringToArray(mlcds)
			MlcdsSequenceR := MlcdsSequence
			MlcdsSequenceR = Reverse(MlcdsSequenceR)
			MlcdsSeqLength := len(MlcdsSequence) - 1
			for o := 1; o < 6; o++ {
				MLCDS_TempStr := ""
				if o < 3 {
					MLCDS_TempStr = InitCodonSeq(o, MlcdsSeqLength-1, 3, MlcdsSequence)
				}
				if 2 < o && o < 6 {
					MLCDS_TempStr = InitCodonSeq(o, MlcdsSeqLength-1, 3, MlcdsSequenceR)
				}
				MLCDS_array := strings.Split(MLCDS_TempStr, " ")
				MLCDS_array = MLCDS_array[:len(MLCDS_array)-1]
				otherNum := float32(0)
				mlcdsArrayLen := len(MLCDS_array) - 1
				for j := 0; j < mlcdsArrayLen; j++ {
					temp2 := MLCDS_array[j] + MLCDS_array[j+1]
					byteMatchResult, _ := regexp.Match(`[atcg]{6}`, []byte(temp2))
					if byteMatchResult {
						if v9, ok := HashMatrix[temp2]; ok {
							v1, _ := strconv.ParseFloat(v9, 64)
							otherNum = otherNum + float32(v1)
						}
					}
				}
				mlcdsArrayLen = mlcdsArrayLen + 2
				otherNum = otherNum / float32(mlcdsArrayLen)
				OtherCdsArray = append(OtherCdsArray, otherNum)
			}
			score_distance := float32(0)
			for m := 0; m < len(OtherCdsArray); m++ {
				score_distance += mScore - OtherCdsArray[m]
			}
			score_distance = score_distance / 5

			out_pos := Pos[orf_index]
			M_length := LengthStoreArray[orf_index]
			length_total_score := float32(0)
			for p := 0; p < len(LengthStoreArray); p++ {
				length_total_score = length_total_score + float32(LengthStoreArray[p])
			}
			length_precent := float32(M_length) / length_total_score
			Coding_Array_one := make([]string, 0)
			codonArr := GetAlphabetMap()
			for n := 0; n < len(o_arr); n++ {
				temp1 := o_arr[n]
				byteMatchResult, _ := regexp.Match(`[atcg{3}]`, []byte(temp1))
				if byteMatchResult && temp1 != "taa" && temp1 != "tag" && temp1 != "tga" {
					if v, ok := codonArr.Load(temp1); ok {
						v1 := v.(int)
						v2 := v1 + 1
						codonArr.LoadOrStore(temp1, v2)
					}
				}
			}

			codonArr.Range(func(key, value interface{}) bool {
				v := fmt.Sprintf("%v", value)
				Coding_Array_one = append(Coding_Array_one, v)
				return true
			})
			C_num1 := 0
			for i := 0; i < len(Coding_Array_one); i++ {
				v1, _ := strconv.Atoi(Coding_Array_one[i])
				C_num1 = C_num1 + v1
			}
			if C_num1 == 0 {
				C_num1 = 1
			}
			for i := 0; i < len(Coding_Array_one); i++ {
				v1, _ := strconv.Atoi(Coding_Array_one[i])
				v2 := v1 / C_num1
				Coding_Array_one[i] = fmt.Sprintf("%v", v2)
			}
			Array_Str := strings.Join(Coding_Array_one, " ")
			PROPERTY_STR := fmt.Sprintf("%v %v %v %v %v %v\n", M, M_length, mScore, length_precent, score_distance, Array_Str)
			DETIL_STR := fmt.Sprintf("%v;;;;; %v %v %v\n", Label, out_pos, mScore, DetilLen)
			score, err := os.OpenFile(TempScore, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
			if err != nil {
				Error("OpenFile : [%s] Err : [%s]", TempScore, err.Error())
				return
			}
			_, err = score.WriteString(PROPERTY_STR)
			if err != nil {
				Error("WriteString Err : [%s]", err.Error())
			}
			detil, err := os.OpenFile(TempDetil, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
			if err != nil {
				Error("OpenFile : [%s] Err : [%s]", TempDetil, err.Error())
				return
			}
			_, err = detil.WriteString(DETIL_STR)
			if err != nil {
				Error("WriteString Err : [%s]", err.Error())
			}
		}(sm, fastqSeqArrTmp[d], labelArrTmp[d], this.TempScore, this.TempDetil, this.Thread)
	}
	sm.Wait()
}

type Float32List []float32

func (s Float32List) Len() int           { return len(s) }
func (s Float32List) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s Float32List) Less(i, j int) bool { return s[i] < s[j] }
