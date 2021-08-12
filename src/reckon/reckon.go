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
	OS_MAX_VALUE    = sync.Map{} // map[i]map[d][]float32
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
	i := 1
	for k, v := range sequenceArr {
		sm.Add(1)
		go compare(sm, i, v, k, HashMatrix)
		i++
	}
	sm.Wait()
}

func compare(sa *gsema.Semaphore, idx int, Seq, Label string, HashMatrix map[string]string) {
	defer sa.Done()
	DetilLen := len(Seq)
	tran_fir_seq := strings.ToLower(Seq)
	tran_sec_seq := strings.ReplaceAll(tran_fir_seq, "u", "t")
	sequenceProcessArr := StringToArray(tran_sec_seq)
	sequenceProcessArrR := Reverse(sequenceProcessArr)
	slen := len(sequenceProcessArr) - 1
	var wg sync.WaitGroup
	for o := 0; o < 6; o++ {
		wg.Add(1)
		go multilayerComparison(&wg, o, idx, slen, sequenceProcessArr, sequenceProcessArrR, HashMatrix)
	}
	wg.Wait()
	rmv, _ := OS_MAX_VALUE.Load(idx)
	rv := rmv.([]float32)
	rMaxValue := rv[:]
	sort.Sort(Float32List(rMaxValue))
	rMaxValue = ReverseFloats32(rMaxValue)
	M := rMaxValue[0]
	orf_index := 0
	for k, v := range rv {
		if v == M {
			orf_index = k
		}
	}
	omx, _ := OS_MAX.Load(idx)
	mx := omx.([]string)
	o_tmp_arr := strings.Split(mx[orf_index], " ")
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
			}
		}
	}
	SequenceLen = SequenceLen + 2
	mScore = mScore / float32(SequenceLen)
	mlcds := strings.Join(o_arr, " ")
	mlcdsSequence := StringToArray(mlcds)
	mlcdsSequenceR := mlcdsSequence
	mlcdsSequenceR = Reverse(mlcdsSequenceR)
	mlen := len(mlcdsSequence) - 1
	for o := 1; o < 6; o++ {
		wg.Add(1)
		go multilayerComparisonTwo(&wg, o, idx, mlen, mlcdsSequence, mlcdsSequenceR, HashMatrix)
	}
	wg.Wait()
	score_distance := float32(0)
	ooc, _ := OS_OTHER_CDS.Load(idx)
	oc := ooc.([]float32)

	for _, v := range oc {
		score_distance += mScore - v
	}
	score_distance = score_distance / 5
	op, _ := OS_POS.Load(idx)
	p := op.([]string)
	out_pos := p[orf_index]
	ols, _ := OS_LENGTH_STORE.Load(idx)
	ls := ols.([]int)
	mlength := ls[orf_index]
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
	PROPERTY_STR := fmt.Sprintf("%v %v %v %v %v %v", M, mlength, mScore, length_precent, score_distance, Array_Str)
	OS_PROPERTY = append(OS_PROPERTY, PROPERTY_STR)
	DETIL_STR := fmt.Sprintf("%v;;;;; %v %v %v", Label, out_pos, mScore, DetilLen)
	OS_DETIL = append(OS_DETIL, DETIL_STR)
}

func multilayerComparison(wg *sync.WaitGroup, o, idx, slen int, sequenceProcessArr, sequenceProcessArrR []string, HashMatrix map[string]string) {
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
			num := float32(0)
			SingleArray := make([]string, 0)
			for t := EachCodon; t < WindowStep+EachCodon; t++ {
				SingleArray = append(SingleArray, TempArray[t])
			}
			for w := range XRangeInt(0, len(SingleArray), 2) {
				v := SingleArray[w] + SingleArray[w+1]
				byteMatchResult, _ := regexp.Match(`[atcg]{6}`, []byte(v))
				if byteMatchResult {
					v1, _ := strconv.ParseFloat(HashMatrix[v], 64)
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
		outStr := ""
		for out := Start; out < End+1; out++ {
			outStr = outStr + TempArray[out] + " "
		}
		Start = Start * 3
		End = End * 3
		Position := fmt.Sprintf("%v %v", Start, End)
		p := make([]string, 0)
		op, _ := OS_POS.Load(idx)
		if op != nil {
			p = op.([]string)
		}
		p = append(p, Position)
		OS_POS.Delete(idx)
		OS_POS.Store(idx, p)

		mv := make([]float32, 0)
		omv, _ := OS_MAX_VALUE.Load(idx)
		if omv != nil {
			mv = omv.([]float32)
		}
		mv = append(mv, Max)
		OS_MAX_VALUE.Delete(idx)
		OS_MAX_VALUE.Store(idx, mv)

		mx := make([]string, 0)
		omx, _ := OS_MAX.Load(idx)
		if omx != nil {
			mx = omx.([]string)
		}
		mx = append(mx, outStr)
		OS_MAX.Delete(idx)
		OS_MAX.Store(idx, mx)

		OutParray := strings.Split(outStr, " ")
		max_length := len(OutParray) - 1

		ls := make([]int, 0)
		ols, _ := OS_LENGTH_STORE.Load(idx)
		if ols != nil {
			ls = ols.([]int)
		}
		ls = append(ls, max_length)
		OS_LENGTH_STORE.Delete(idx)
		OS_LENGTH_STORE.Store(idx, ls)

	} else {
		num := float32(0)
		for w := range XRangeInt(0, seqLength, 2) {
			v := TempArray[w] + TempArray[w+1]
			byteMatchResult, _ := regexp.Match(`[atcg]{6}`, []byte(v))
			if byteMatchResult {
				if v9, ok := HashMatrix[v]; ok {
					v1, _ := strconv.ParseFloat(v9, 64)
					num = num + float32(v1)
				}
			}
		}
		outStr := strings.Join(TempArray, " ")
		p := make([]string, 0)
		op, _ := OS_POS.Load(idx)
		if op != nil {
			p = op.([]string)
		}
		p = append(p, "Full Length")
		OS_POS.Delete(idx)
		OS_POS.Store(idx, p)

		mv := make([]float32, 0)
		omv, _ := OS_MAX_VALUE.Load(idx)
		if omv != nil {
			mv = omv.([]float32)
		}
		mv = append(mv, num)
		OS_MAX_VALUE.Delete(idx)
		OS_MAX_VALUE.Store(idx, mv)

		mx := make([]string, 0)
		omx, _ := OS_MAX.Load(idx)
		if omx != nil {
			mx = omx.([]string)
		}
		mx = append(mx, outStr)
		OS_MAX.Delete(idx)
		OS_MAX.Store(idx, mx)

		ls := make([]int, 0)
		ols, _ := OS_LENGTH_STORE.Load(idx)
		if ols != nil {
			ls = ols.([]int)
		}
		ls = append(ls, seqLength)
		OS_LENGTH_STORE.Delete(idx)
		OS_LENGTH_STORE.Store(idx, ls)

	}
}

func multilayerComparisonTwo(wg *sync.WaitGroup, o, idx, mlen int, mlcdsSequence, mlcdsSequenceR []string, HashMatrix map[string]string) {
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
	otherNum := float32(0)
	mlcdsArrayLen := len(MLCDS_array) - 1
	for w := range XRangeInt(0, mlcdsArrayLen, 2) {
		v := MLCDS_array[w] + MLCDS_array[w+1]
		byteMatchResult, _ := regexp.Match(`[atcg]{6}`, []byte(v))
		if byteMatchResult {
			if v9, ok := HashMatrix[v]; ok {
				v1, _ := strconv.ParseFloat(v9, 64)
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
