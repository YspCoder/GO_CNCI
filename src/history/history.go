package history

/*
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
*/
/*
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
*/

//score, err := os.OpenFile(TempScore, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
//if err != nil {
//	Error("OpenFile : [%s] Err : [%s]", TempScore, err.Error())
//	return
//}
//_, err = score.WriteString(PROPERTY_STR)
//if err != nil {
//	Error("WriteString Err : [%s]", err.Error())
//}
//detil, err := os.OpenFile(TempDetil, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
//if err != nil {
//	Error("OpenFile : [%s] Err : [%s]", TempDetil, err.Error())
//	return
//}
//_, err = detil.WriteString(DETIL_STR)
//if err != nil {
//	Error("WriteString Err : [%s]", err.Error())
//}

//MaxValue := make([]float32, 0)
//MaxString := make([]string, 0)
//LengthStoreArray := make([]int, 0)
//Pos := make([]string, 0)
//OtherCdsArray := make([]float32, 0)

//var err error
//score, err := os.Create(this.TempScore)
//if err != nil {
//	Error("Create tmp error![%v]\n", err.Error())
//	return
//}
//_ = score.Close()
//detil, err := os.Create(this.TempDetil)
//if err != nil {
//	Error("Create tmp error![%v]\n", err.Error())
//	return
//}
//_ = detil.Close()
