package main

import (
	. "GO_CNCI/src/base"
	"GO_CNCI/src/merge"
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
	if len(os.Args) < 5 {
		Info("Insufficient required parameters")
		Info("./GO_CNCI reference_folder inputFile outDir libsvmpath thread")
		Info("./GO_CNCI ./CNCI_Parameters ./94d6346_candidate.fa ./test ./libsvm 8")
		return
	}
	CNCI_Parameters := os.Args[1]
	inputFile := os.Args[2]
	outDir := os.Args[3]
	libsvm_path := os.Args[4]
	thread, err := strconv.Atoi(os.Args[5])
	if err != nil {
		Info("Please enter a positive integer -- thread")
		return
	}
	out_temp := fmt.Sprintf("%s/temp", outDir)
	if !PathExists(out_temp) {
		err := os.MkdirAll(out_temp, os.ModePerm)
		if err != nil {
			Error("Create Temp Err : [%s]", err.Error())
			return
		}
	}
	hashMatrix := ReadFileMatrix(CNCI_Parameters + "/CNCI_matrix")
	sequence_Arr := ReadFileArray(inputFile)
	sLen := len(sequence_Arr) - 1
	sequence_Arr = sequence_Arr[:sLen]
	fastArray := TwoLineFasta(sequence_Arr)
	Label_Array, Fasta_Seq_Array := Tran_checkSeq(fastArray)

	TOT_STRING := GetLabelArray(Label_Array, Fasta_Seq_Array)
	Info("-------Start splitting file------")
	SplitFile(TOT_STRING, thread, out_temp)
	Info("--------End of split file-------")
	Info("--------Start calculation-------")
	var wgs sync.WaitGroup
	for i := 1; i <= thread; i++ {
		wgs.Add(1)
		rk := reckon.New()
		rk.HashMatrix = hashMatrix
		rk.TempInput = fmt.Sprintf("%s/GO_CNCI_file%v", out_temp, i)
		rk.TempScore = fmt.Sprintf("%s/GO_CNCI_file_score%v", out_temp, i)
		rk.TempDetil = fmt.Sprintf("%s/GO_CNCI_file_detil%v", out_temp, i)
		rk.Thread = thread
		go rk.Init(&wgs)
	}
	wgs.Wait()
	Info("--------End of calculation-------")
	outfile := fmt.Sprintf("%s/pro", outDir)
	Info("---------Start merging files--------")
	score_path := fmt.Sprintf("%s/GO_CNCI_score", outDir)
	detil_path := fmt.Sprintf("%s/GO_CNCI_detil", outDir)
	err = merge.Merge(out_temp, score_path, detil_path, thread)
	if err != nil {
		Error("Merge err : [%s]", err.Error())
		return
	}
	score_array := ReadFileArray(score_path)
	scoreSLength := len(score_array) - 1
	score_array = score_array[:scoreSLength]
	detil_array := ReadFileArray(detil_path)
	detilSLength := len(detil_array) - 1
	detil_array = detil_array[:detilSLength]
	err = merge.AddSvmLabel(score_array, outfile)
	if err != nil {
		Error("AddSvmLabel err : [%s]", err.Error())
		return
	}
	Info("---------End of merge file-------")
	SvmPutFileName := fmt.Sprintf("%s/svm", outDir)
	SvmFile := fmt.Sprintf("%s/file", outDir)
	SvmTmp := fmt.Sprintf("%s/tmp", outDir)
	Info("-------Start vector calculation------")
	err = Libsvm(outfile, SvmPutFileName, SvmFile, SvmTmp, libsvm_path, CNCI_Parameters)
	if err != nil {
		Error("Libsvm err : [%s]", err.Error())
		return
	}
	Info("----------End of vector calculation--------")
	Info("Start output file")
	FirResult := PutResult(detil_array, SvmFile)
	SvmFinalResult := fmt.Sprintf("%s/GO_CNCI.index", outDir)
	PrintResult(FirResult, SvmFinalResult)
	Info("---------End of output file----------")
	cost := time.Since(start)
	Info("Time use [%s]", cost)
}
