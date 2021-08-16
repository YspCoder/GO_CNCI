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

func (this *Reckon) Init(wg *sync.WaitGroup) {
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
	sequenceProcessArrR := StringToArray(tran_fir_seq)
	sequenceProcessArrR = Reverse(sequenceProcessArrR)
	slen := len(sequenceProcessArr) - 1
	var wg sync.WaitGroup
	for o := 0; o < 6; o++ {
		wg.Add(1)
		go multilayerComparison(&wg, o, slen, Label, sequenceProcessArr, sequenceProcessArrR, HashMatrix)
	}
	wg.Wait()
	rmv, _ := OS_MAX_VALUE.Load(Label)
	rv := rmv.([]map[float64]string)
	omx, _ := OS_MAX.Load(Label)
	mx := omx.([]float64)
	rMaxValue := mx[:]
	sort.Float64s(rMaxValue)
	rMaxValue = ReverseFloats64(rMaxValue)
	M := rMaxValue[0]
	tmpStr := ""
	for _, v := range rv {
		tmpStr = v[M]
		if tmpStr != "" {
			break
		}
	}
	o_tmp_arr := strings.Split(tmpStr, " ")
	var o_arr = make([]string, 0)
	o_arr = o_tmp_arr[:len(o_tmp_arr)-1]
	SequenceLen := len(o_arr) - 1
	var mScore float64
	for j := 0; j < SequenceLen; j++ {
		v := o_arr[j] + o_arr[j+1]
		byteMatchResult, _ := regexp.Match(`[atcg]{6}`, []byte(v))
		if byteMatchResult {
			if v9, ok := HashMatrix[v]; ok {
				v1, _ := strconv.ParseFloat(v9, 64)
				mScore = mScore + float64(v1)
			}
		}
	}
	SequenceLen = SequenceLen + 2
	mScore = mScore / float64(SequenceLen)
	mlcds := strings.Join(o_arr, "")
	mlcdsSequence := StringToArray(mlcds)
	mlcdsSequenceR := StringToArray(mlcds)
	mlcdsSequenceR = Reverse(mlcdsSequenceR)
	mlen := len(mlcdsSequence) - 1
	for h := 1; h < 6; h++ {
		wg.Add(1)
		go multilayerComparisonTwo(&wg, h, mlen, Label, mlcdsSequence, mlcdsSequenceR, HashMatrix)
	}
	wg.Wait()
	var score_distance float64
	ooc, _ := OS_OTHER_CDS.Load(Label)
	oc := ooc.([]float64)

	for _, v := range oc {
		score_distance += mScore - v
	}
	score_distance = score_distance / 5
	op, _ := OS_POS.Load(Label)
	p := op.([]map[float64]string)
	out_pos := ""
	for _, v := range p {
		out_pos = v[M]
		if out_pos != "" {
			break
		}
	}
	ols, _ := OS_LENGTH_STORE.Load(Label)
	ls := ols.([]map[float64]int)
	mlength := 0
	for _, v := range ls {
		mlength = v[M]
		if mlength != 0 {
			break
		}
	}
	length_total_score := float64(0)
	for _, v := range ls {
		for _, t := range v {
			length_total_score = length_total_score + float64(t)
		}
	}
	length_precent := float64(mlength) / length_total_score
	codingArray := make([]string, 0)
	codonArr := map[string]float64{"ttt": 0, "ttc": 0, "tta": 0, "ttg": 0, "tct": 0, "tcc": 0, "tca": 0, "tcg": 0, "tat": 0, "tac": 0, "tgt": 0, "tgc": 0, "tgg": 0, "ctt": 0, "ctc": 0, "cta": 0, "ctg": 0, "cct": 0, "ccc": 0, "cca": 0, "ccg": 0, "cat": 0, "cac": 0, "caa": 0, "cag": 0, "cgt": 0, "cgc": 0, "cga": 0, "cgg": 0, "att": 0, "atc": 0, "ata": 0, "atg": 0, "act": 0, "acc": 0, "aca": 0, "acg": 0, "aat": 0, "aac": 0, "aaa": 0, "aag": 0, "agt": 0, "agc": 0, "aga": 0, "agg": 0, "gtt": 0, "gtc": 0, "gta": 0, "gtg": 0, "gct": 0, "gcc": 0, "gca": 0, "gcg": 0, "gat": 0, "gac": 0, "gaa": 0, "gag": 0, "ggt": 0, "ggc": 0, "gga": 0, "ggg": 0}
	ca := []string{"ttt", "ttc", "tta", "ttg", "tct", "tcc", "tca", "tcg", "tat", "tac", "tgt", "tgc", "tgg", "ctt", "ctc", "cta", "ctg", "cct", "ccc", "cca", "ccg", "cat", "cac", "caa", "cag", "cgt", "cgc", "cga", "cgg", "att", "atc", "ata", "atg", "act", "acc", "aca", "acg", "aat", "aac", "aaa", "aag", "agt", "agc", "aga", "agg", "gtt", "gtc", "gta", "gtg", "gct", "gcc", "gca", "gcg", "gat", "gac", "gaa", "gag", "ggt", "ggc", "gga", "ggg"}

	for _, v := range o_arr {
		byteMatchResult, _ := regexp.Match(`[atcg{3}]`, []byte(v))
		if byteMatchResult && v != "taa" && v != "tag" && v != "tga" {
			if v1, ok := codonArr[v]; ok {
				v2 := v1 + 1
				codonArr[v] = v2
			}
		}
	}
	var CNum1 float64

	for _, v := range codonArr {
		CNum1 = CNum1 + v
	}
	if CNum1 == 0 {
		CNum1 = 1
	}
	for i := 0; i < len(ca); i++ {
		v1 := codonArr[ca[i]] / CNum1
		codingArray = append(codingArray, fmt.Sprintf("%v", v1))
	}

	Array_Str := strings.Join(codingArray, " ")
	PROPERTY_STR := fmt.Sprintf("%v %v %v %v %v %v %v", Label, M, mlength, mScore, length_precent, score_distance, Array_Str)
	OS_PROPERTY = append(OS_PROPERTY, PROPERTY_STR)
	DETIL_STR := fmt.Sprintf("%v;;;;; %v %v %v", Label, out_pos, mScore, DetilLen)
	OS_DETIL = append(OS_DETIL, DETIL_STR)
}

func multilayerComparison(wg *sync.WaitGroup, o, slen int, idx string, sequenceProcessArr, sequenceProcessArrR []string, HashMatrix map[string]string) {
	defer wg.Done()
	CodonScore := make([]float64, 0)
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
			var num float64
			SingleArray := make([]string, 0)
			for t := EachCodon; t < WindowStep+EachCodon; t++ {
				SingleArray = append(SingleArray, TempArray[t])
			}
			SinLen := len(SingleArray) - 1
			for n := 0; n < SinLen; n++ {
				v := SingleArray[n] + SingleArray[n+1]
				byteMatchResult, _ := regexp.Match(`[atcg]{6}`, []byte(v))
				if byteMatchResult {
					if v9, ok := HashMatrix[v]; ok {
						v1, _ := strconv.ParseFloat(v9, 64)
						num = num + v1
					}
				}
			}
			num = num / 50
			CodonScore = append(CodonScore, num)
		}
		Start := 0
		End := 0
		var Max float64

		for r := 0; r < len(CodonScore); r++ {
			var sum float64
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
		p := make([]map[float64]string, 0)
		p1 := make(map[float64]string, 0)
		op, _ := OS_POS.Load(idx)
		if op != nil {
			p = op.([]map[float64]string)
		}
		p1[Max] = Position
		p = append(p, p1)
		OS_POS.Delete(idx)
		OS_POS.Store(idx, p)

		mv := make([]map[float64]string, 0)
		mv1 := make(map[float64]string, 0)
		omv, _ := OS_MAX_VALUE.Load(idx)
		if omv != nil {
			mv = omv.([]map[float64]string)
		}
		mv1[Max] = outStr
		mv = append(mv, mv1)
		OS_MAX_VALUE.Delete(idx)
		OS_MAX_VALUE.Store(idx, mv)

		mx := make([]float64, 0)
		omx, _ := OS_MAX.Load(idx)
		if omx != nil {
			mx = omx.([]float64)
		}
		mx = append(mx, Max)
		OS_MAX.Delete(idx)
		OS_MAX.Store(idx, mx)

		OutParray := strings.Split(outStr, " ")
		max_length := len(OutParray) - 1

		ls := make([]map[float64]int, 0)
		ls1 := make(map[float64]int, 0)
		ols, _ := OS_LENGTH_STORE.Load(idx)
		if ols != nil {
			ls = ols.([]map[float64]int)
		}
		ls1[Max] = max_length
		ls = append(ls, ls1)
		OS_LENGTH_STORE.Delete(idx)
		OS_LENGTH_STORE.Store(idx, ls)

	} else {
		var num float64
		for w := range XRangeInt(0, seqLength, 2) {
			v := TempArray[w] + TempArray[w+1]
			byteMatchResult, _ := regexp.Match(`[atcg]{6}`, []byte(v))
			if byteMatchResult {
				if v9, ok := HashMatrix[v]; ok {
					v1, _ := strconv.ParseFloat(v9, 64)
					num = num + float64(v1)
				}
			}
		}
		outStr := strings.Join(TempArray, " ")
		p := make([]map[float64]string, 0)
		p1 := make(map[float64]string, 0)
		op, _ := OS_POS.Load(idx)
		if op != nil {
			p = op.([]map[float64]string)
		}
		p1[num] = "Full Length"
		p = append(p, p1)
		OS_POS.Delete(idx)
		OS_POS.Store(idx, p)

		mv := make([]map[float64]string, 0)
		mv1 := make(map[float64]string, 0)
		omv, _ := OS_MAX_VALUE.Load(idx)
		if omv != nil {
			mv = omv.([]map[float64]string)
		}
		mv1[num] = outStr
		mv = append(mv, mv1)
		OS_MAX_VALUE.Delete(idx)
		OS_MAX_VALUE.Store(idx, mv)

		mx := make([]float64, 0)
		omx, _ := OS_MAX.Load(idx)
		if omx != nil {
			mx = omx.([]float64)
		}
		mx = append(mx, num)
		OS_MAX.Delete(idx)
		OS_MAX.Store(idx, mx)

		ls := make([]map[float64]int, 0)
		ls1 := make(map[float64]int, 0)
		ols, _ := OS_LENGTH_STORE.Load(idx)
		if ols != nil {
			ls = ols.([]map[float64]int)
		}
		ls1[num] = seqLength
		ls = append(ls, ls1)
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
	var otherNum float64
	mlcdsArrayLen := len(MLCDS_array) - 1
	for i := 0; i < mlcdsArrayLen; i++ {
		v := MLCDS_array[i] + MLCDS_array[i+1]
		byteMatchResult, _ := regexp.Match(`[atcg]{6}`, []byte(v))
		if byteMatchResult {
			if v9, ok := HashMatrix[v]; ok {
				v1, _ := strconv.ParseFloat(v9, 64)
				otherNum = otherNum + float64(v1)
			}
		}
	}

	mlcdsArrayLen = mlcdsArrayLen + 2
	otherNum = otherNum / float64(mlcdsArrayLen)
	oc := make([]float64, 0)
	ooc, _ := OS_OTHER_CDS.Load(idx)
	if ooc != nil {
		oc = ooc.([]float64)
	}
	oc = append(oc, otherNum)
	OS_OTHER_CDS.Delete(idx)
	OS_OTHER_CDS.Store(idx, oc)
}
