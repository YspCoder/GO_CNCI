package main

import (
	. "GO_CNCI/src/base"
	"GO_CNCI/src/reckon"
	. "GO_CNCI/src/utils"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

func main() {
	start := time.Now()
	CNCI_matrix := os.Args[1]
	inputFile := os.Args[2]
	number, _ := strconv.Atoi(os.Args[3])
	out_temp := os.Args[4]
	logFile := os.Args[5]
	_ = os.Args[5]

	hashMatrix := ReadFileMatrix(CNCI_matrix)
	sequence_Arr := ReadFileArray(inputFile)
	sLen := len(sequence_Arr) - 1
	sequence_Arr = sequence_Arr[:sLen]
	fastArray := TwoLineFasta(sequence_Arr)
	Label_Array, Fasta_Seq_Array := Tran_checkSeq(fastArray, logFile)

	TOT_STRING := GetLabelArray(Label_Array, Fasta_Seq_Array)
	SplitFile(TOT_STRING, number, out_temp)
	var wg sync.WaitGroup
	for i := 1; i < number; i++ {
		wg.Add(1)
		rk := reckon.New()
		rk.HashMatrix = hashMatrix
		rk.TempInput = fmt.Sprintf("%s/GO_CNCI_file%v", out_temp, i)
		rk.TempScore = fmt.Sprintf("%s/GO_CNCI_file_score%v", out_temp, number)
		rk.TempDetil = fmt.Sprintf("%s/GO_CNCI_file_detil%v", out_temp, number)
		go rk.Init(&wg)
	}
	wg.Wait()

	cost := time.Since(start)
	Info("Cost=[%s]", cost)
}
