package main

import (
	. "GO_CNCI/src/base"
	"GO_CNCI/src/merge"
	"GO_CNCI/src/reckon"
	. "GO_CNCI/src/utils"
	"fmt"
	"github.com/EDDYCJY/gsema"
	"os"
	"strconv"
	"time"
)

func main() {
	start := time.Now()
	if len(os.Args) < 6 {
		Info("Insufficient required parameters")
		Info("./GO_CNCI reference_folder inputFile number_of_file_partitions outDir libsvmpath thread")
		Info("./GO_CNCI ./CNCI_Parameters ./94d6346_candidate.fa 10 ./test ./libsvm 8")
		return
	}
	CNCI_Parameters := os.Args[1]
	inputFile := os.Args[2]
	number, err := strconv.Atoi(os.Args[3])
	if err != nil {
		Info("Please enter a positive integer")
		return
	}
	outDir := os.Args[4]
	libsvm_path := os.Args[5]
	thread, err := strconv.Atoi(os.Args[6])
	if err != nil {
		Info("Please enter a positive integer")
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
	sema := gsema.NewSemaphore(thread)
	hashMatrix := ReadFileMatrix(CNCI_Parameters + "/CNCI_matrix")
	sequence_Arr := ReadFileArray(inputFile)
	sLen := len(sequence_Arr) - 1
	sequence_Arr = sequence_Arr[:sLen]
	fastArray := TwoLineFasta(sequence_Arr)
	Label_Array, Fasta_Seq_Array := Tran_checkSeq(fastArray)

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
	outfile := fmt.Sprintf("%s/pro", outDir)
	Info("---------Start merging files--------")
	detilArray, err := merge.Merge(out_temp, outfile, number)
	if err != nil {
		Error("Merge err : [%s]", err.Error())
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
	FirResult := PutResult(detilArray, SvmFile)
	SvmFinalResult := fmt.Sprintf("%s/GO_CNCI.index", outDir)
	PrintResult(FirResult, SvmFinalResult)
	Info("---------End of output file----------")
	cost := time.Since(start)
	Info("Time use [%s]", cost)
}
