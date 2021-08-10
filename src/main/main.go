package main

import (
	. "GO_CNCI/src/base"
	"GO_CNCI/src/merge"
	"GO_CNCI/src/reckon"
	. "GO_CNCI/src/utils"
	"fmt"
	"github.com/EDDYCJY/gsema"
	"time"
)

func main() {
	start := time.Now()
	CNCI_matrix := "./CNCI_Parameters/CNCI_matrix"
	inputFile := "./94d6346_candidate.fa"
	number := 6
	out_temp := "./test"
	logFile := "./"
	sema := gsema.NewSemaphore(12)
	hashMatrix := ReadFileMatrix(CNCI_matrix)
	sequence_Arr := ReadFileArray(inputFile)
	sLen := len(sequence_Arr) - 1
	sequence_Arr = sequence_Arr[:sLen]
	fastArray := TwoLineFasta(sequence_Arr)
	Label_Array, Fasta_Seq_Array := Tran_checkSeq(fastArray, logFile)

	TOT_STRING := GetLabelArray(Label_Array, Fasta_Seq_Array)
	Info("-------Start splitting file------")
	SplitFile(TOT_STRING, number, out_temp)
	Info("--------End of split file-------")
	Info("--------Start calculation-------")
	for i := 1; i <= number; i++ {
		sema.Add(1)
		rk := reckon.New()
		rk.HashMatrix = hashMatrix
		rk.TempInput = fmt.Sprintf("%s/GO_CNCI_file%v", out_temp, i)
		rk.TempScore = fmt.Sprintf("%s/GO_CNCI_file_score%v", out_temp, i)
		rk.TempDetil = fmt.Sprintf("%s/GO_CNCI_file_detil%v", out_temp, i)
		go rk.Init(sema)
	}
	sema.Wait()
	Info("--------End of calculation-------")
	outfile := fmt.Sprintf("%s/pro", out_temp)
	Info("---------Start merging files--------")
	detilArray, err := merge.Merge(out_temp, outfile, number)
	if err != nil {
		Error("Merge err : [%s]", err.Error())
		return
	}
	Info("---------End of merge file-------")
	SvmPutFileName := fmt.Sprintf("%s/svm", out_temp)
	SvmFile := fmt.Sprintf("%s/file", out_temp)
	SvmTmp := fmt.Sprintf("%s/tmp", out_temp)
	Info("-------Start vector calculation------")
	err = Libsvm(outfile, SvmPutFileName, SvmFile, SvmTmp)
	if err != nil {
		Error("Libsvm err : [%s]", err.Error())
		return
	}
	Info("----------End of vector calculation--------")
	Info("Start output file")
	FirResult := PutResult(detilArray, SvmFile)
	SvmFinalResutl := fmt.Sprintf("%s/GO_CNCI.index", out_temp)
	PrintResult(FirResult, SvmFinalResutl)
	Info("---------End of output file----------")
	cost := time.Since(start)
	Info("Time use [%s]", cost)
}
