package reckon

import (
	. "GO_CNCI/src/utils"
	"fmt"
	"github.com/EDDYCJY/gsema"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
)

var (
	OS_MAX_VALUE    = sync.Map{}
	OS_MAX          = sync.Map{}
	OS_LENGTH_STORE = sync.Map{}
	OS_POS          = sync.Map{}
	OS_PROPERTY     = make([]string, 0)
	OS_DETIL        = make([]string, 0)
	OS_OTHER_CDS    = sync.Map{}
)

type Float32List []float32

func (s Float32List) Len() int           { return len(s) }
func (s Float32List) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s Float32List) Less(i, j int) bool { return s[i] < s[j] }

type Reckon struct {
	FileInput  interface{}
	HashMatrix map[string]string
	Thread     int
}

func New() *Reckon {
	return &Reckon{
		FileInput:  nil,
		HashMatrix: nil,
		Thread:     0,
	}
}

func (this *Reckon) Init(wg *gsema.Semaphore) {
	defer wg.Done()

	var sm = gsema.NewSemaphore(this.Thread)
	HashMatrix := this.HashMatrix
	sequenceArr := this.FileInput.(map[string]string)
	for k, v := range sequenceArr {
		sm.Add(1)
		go compare(sm, v, k, HashMatrix)
	}
	sm.Wait()
}

func compare(sa *gsema.Semaphore, Seq, Label string, HashMatrix map[string]string) {
	defer sa.Done()
	DetilLen := len(Seq)
	tran_fir_seq := strings.ToLower(Seq)
	sequenceProcessArr := StringToArray(tran_fir_seq)
	sequenceProcessArrR := Reverse(sequenceProcessArr)
	slen := len(sequenceProcessArr) - 1
	var wg sync.WaitGroup
	for o := 0; o < 6; o++ {
		wg.Add(1)
		go multilayerComparison(&wg, o, slen, Label, sequenceProcessArr, sequenceProcessArrR, HashMatrix)
	}
	wg.Wait()
	rmv, _ := OS_MAX_VALUE.Load(Label)
	rv := rmv.(map[float32]string)
	omx, _ := OS_MAX.Load(Label)
	mx := omx.([]float32)
	rMaxValue := mx[:]
	sort.Sort(Float32List(rMaxValue))
	rMaxValue = ReverseFloats32(rMaxValue)
	M := rMaxValue[0]
	o_tmp_arr := strings.Split(rv[M], " ")
	var o_arr = make([]string, 0)
	o_arr = o_tmp_arr[:len(o_tmp_arr)-1]
	SequenceLen := len(o_arr) - 1
	var mScore float32
	for j := 0; j < SequenceLen; j++ {
		v := o_arr[j] + o_arr[j+1]
		byteMatchResult, _ := regexp.Match(`[atcg]{6}`, []byte(v))
		if byteMatchResult {
			if v9, ok := HashMatrix[v]; ok {
				v1, _ := strconv.ParseFloat(v9, 32)
				mScore = mScore + float32(v1)
			}
		}
	}
	SequenceLen = SequenceLen + 2
	mScore = mScore / float32(SequenceLen)
	mlcds := strings.Join(o_arr, "")
	mlcdsSequence := StringToArray(mlcds)
	mlcdsSequenceR := mlcdsSequence
	mlcdsSequenceR = Reverse(mlcdsSequenceR)
	mlen := len(mlcdsSequence) - 1
	for h := 1; h < 6; h++ {
		wg.Add(1)
		go multilayerComparisonTwo(&wg, h, mlen, Label, mlcdsSequence, mlcdsSequenceR, HashMatrix)
	}
	wg.Wait()
	var score_distance float32
	ooc, _ := OS_OTHER_CDS.Load(Label)
	oc := ooc.([]float32)

	for _, v := range oc {
		score_distance += mScore - v
	}
	score_distance = score_distance / 5
	op, _ := OS_POS.Load(Label)
	p := op.(map[float32]string)
	out_pos := p[M]
	ols, _ := OS_LENGTH_STORE.Load(Label)
	ls := ols.(map[float32]int)
	mlength := ls[M]
	length_total_score := float32(0)
	for _, v := range ls {
		length_total_score = length_total_score + float32(v)
	}
	length_precent := float32(mlength) / length_total_score
	codingArray := make([]string, 0)
	codonArr := GetAlphabetMap()
	for _, v := range o_arr {
		byteMatchResult, _ := regexp.Match(`[atcg{3}]`, []byte(v))
		if byteMatchResult && v != "taa" && v != "tag" && v != "tga" {
			if v, ok := codonArr.Load(v); ok {
				v1 := v.(int)
				v2 := v1 + 1
				codonArr.Delete(v)
				codonArr.Store(v, v2)
			}
		}
	}
	codonArr.Range(func(key, value interface{}) bool {
		v := fmt.Sprintf("%v", value)
		codingArray = append(codingArray, v)
		return true
	})
	C_num1 := 0
	for _, v := range codingArray {
		v1, _ := strconv.Atoi(v)
		C_num1 = C_num1 + v1
	}
	if C_num1 == 0 {
		C_num1 = 1
	}
	for k, v := range codingArray {
		v1, _ := strconv.Atoi(v)
		v2 := v1 / C_num1
		codingArray[k] = fmt.Sprintf("%v", v2)
	}
	Array_Str := strings.Join(codingArray, " ")
	PROPERTY_STR := fmt.Sprintf("%v %v %v %v %v %v %v", Label, M, mlength, mScore, length_precent, score_distance, Array_Str)
	OS_PROPERTY = append(OS_PROPERTY, PROPERTY_STR)
	DETIL_STR := fmt.Sprintf("%v;;;;; %v %v %v", Label, out_pos, mScore, DetilLen)
	OS_DETIL = append(OS_DETIL, DETIL_STR)
}

func multilayerComparison(wg *sync.WaitGroup, o, slen int, idx string, sequenceProcessArr, sequenceProcessArrR []string, HashMatrix map[string]string) {
	defer wg.Done()
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
			var num float32
			SingleArray := make([]string, 0)
			for t := EachCodon; t < WindowStep+EachCodon; t++ {
				SingleArray = append(SingleArray, TempArray[t])
			}
			for w := range XRangeInt(0, len(SingleArray), 2) {
				v := SingleArray[w] + SingleArray[w+1]
				byteMatchResult, _ := regexp.Match(`[atcg]{6}`, []byte(v))
				if byteMatchResult {
					if v9, ok := HashMatrix[v]; ok {
						v1, _ := strconv.ParseFloat(v9, 32)
						num = num + float32(v1)
					}
				}
			}
			num = num / float32(WindowStep)
			CodonScore = append(CodonScore, num)
		}
		Start := 0
		End := 0
		var Max float32

		for r := 0; r < len(CodonScore); r++ {
			var sum float32
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
		outStr := ""
		for out := Start; out < End+1; out++ {
			outStr = outStr + TempArray[out] + " "
		}
		Start = Start * 3
		End = End * 3
		Position := fmt.Sprintf("%v %v", Start, End)
		p := make(map[float32]string, 0)
		op, _ := OS_POS.Load(idx)
		if op != nil {
			p = op.(map[float32]string)
		}
		p[Max] = Position
		OS_POS.Delete(idx)
		OS_POS.Store(idx, p)

		mv := make(map[float32]string, 0)
		omv, _ := OS_MAX_VALUE.Load(idx)
		if omv != nil {
			mv = omv.(map[float32]string)
		}
		mv[Max] = outStr
		OS_MAX_VALUE.Delete(idx)
		OS_MAX_VALUE.Store(idx, mv)

		mx := make([]float32, 0)
		omx, _ := OS_MAX.Load(idx)
		if omx != nil {
			mx = omx.([]float32)
		}
		mx = append(mx, Max)
		OS_MAX.Delete(idx)
		OS_MAX.Store(idx, mx)

		OutParray := strings.Split(outStr, " ")
		max_length := len(OutParray) - 1

		ls := make(map[float32]int, 0)
		ols, _ := OS_LENGTH_STORE.Load(idx)
		if ols != nil {
			ls = ols.(map[float32]int)
		}
		ls[Max] = max_length
		OS_LENGTH_STORE.Delete(idx)
		OS_LENGTH_STORE.Store(idx, ls)

	} else {
		var num float32
		for w := range XRangeInt(0, seqLength, 2) {
			v := TempArray[w] + TempArray[w+1]
			byteMatchResult, _ := regexp.Match(`[atcg]{6}`, []byte(v))
			if byteMatchResult {
				if v9, ok := HashMatrix[v]; ok {
					v1, _ := strconv.ParseFloat(v9, 32)
					num = num + float32(v1)
				}
			}
		}
		outStr := strings.Join(TempArray, " ")
		p := make(map[float32]string, 0)
		op, _ := OS_POS.Load(idx)
		if op != nil {
			p = op.(map[float32]string)
		}
		p[num] = "Full Length"
		OS_POS.Delete(idx)
		OS_POS.Store(idx, p)

		mv := make(map[float32]string, 0)
		omv, _ := OS_MAX_VALUE.Load(idx)
		if omv != nil {
			mv = omv.(map[float32]string)
		}
		mv[num] = outStr
		OS_MAX_VALUE.Delete(idx)
		OS_MAX_VALUE.Store(idx, mv)

		mx := make([]float32, 0)
		omx, _ := OS_MAX.Load(idx)
		if omx != nil {
			mx = omx.([]float32)
		}
		mx = append(mx, num)
		OS_MAX.Delete(idx)
		OS_MAX.Store(idx, mx)

		ls := make(map[float32]int, 0)
		ols, _ := OS_LENGTH_STORE.Load(idx)
		if ols != nil {
			ls = ols.(map[float32]int)
		}
		ls[num] = seqLength
		OS_LENGTH_STORE.Delete(idx)
		OS_LENGTH_STORE.Store(idx, ls)

	}
}

func multilayerComparisonTwo(wg *sync.WaitGroup, o, mlen int, idx string, mlcdsSequence, mlcdsSequenceR []string, HashMatrix map[string]string) {
	defer wg.Done()
	MLCDS_TempStr := ""
	if o < 3 {
		MLCDS_TempStr = InitCodonSeq(o, mlen-1, 3, mlcdsSequence)
	}
	if 2 < o && o < 6 {
		MLCDS_TempStr = InitCodonSeq(o, mlen-1, 3, mlcdsSequenceR)
	}
	MLCDS_array := strings.Split(MLCDS_TempStr, " ")
	MLCDS_array = MLCDS_array[:len(MLCDS_array)-1]
	var otherNum float32
	mlcdsArrayLen := len(MLCDS_array) - 1
	for i := 0; i < mlcdsArrayLen; i++ {
		v := MLCDS_array[i] + MLCDS_array[i+1]
		byteMatchResult, _ := regexp.Match(`[atcg]{6}`, []byte(v))
		if byteMatchResult {
			if v9, ok := HashMatrix[v]; ok {
				v1, _ := strconv.ParseFloat(v9, 32)
				otherNum = otherNum + float32(v1)
			}
		}
	}

	mlcdsArrayLen = mlcdsArrayLen + 2
	otherNum = otherNum / float32(mlcdsArrayLen)
	oc := make([]float32, 0)
	ooc, _ := OS_OTHER_CDS.Load(idx)
	if ooc != nil {
		oc = ooc.([]float32)
	}
	oc = append(oc, otherNum)
	OS_OTHER_CDS.Delete(idx)
	OS_OTHER_CDS.Store(idx, oc)
}
