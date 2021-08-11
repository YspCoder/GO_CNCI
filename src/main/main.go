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
	if len(os.Args) < 5 {
		Info("Insufficient required parameters")
		Info("./GO_CNCI reference_folder inputFile outDir libsvmpath thread")
		Info("./GO_CNCI ./CNCI_Parameters ./94d6346_candidate.fa ./test ./libsvm 8")
		return
	}
	cnciParameters := os.Args[1]
	inputFile := os.Args[2]
	outDir := os.Args[3]
	libsvmPath := os.Args[4]
	thread, err := strconv.Atoi(os.Args[5])
	if err != nil {
		Info("Please enter a positive integer -- thread")
		return
	}
	outTemp := fmt.Sprintf("%s/temp", outDir)
	if !PathExists(outTemp) {
		err := os.MkdirAll(outTemp, os.ModePerm)
		if err != nil {
			Error("Create Temp Err : [%s]", err.Error())
			return
		}
	}
	hashMatrix := ReadFileMatrix(cnciParameters + "/CNCI_matrix")
	sequenceArr := ReadFileArray(inputFile)
	sLen := len(sequenceArr) - 1
	sequenceArr = sequenceArr[:sLen]
	fastArray := TwoLineFasta(sequenceArr)
	labelArray, FastqSeqArray := Tran_checkSeq(fastArray)

	tot := GetLabelArray(labelArray, FastqSeqArray)
	Info("-------Start splitting file------")
	SplitFile(tot, thread, outTemp)
	Info("--------End of split file-------")
	Info("--------Start calculation-------")
	seam := gsema.NewSemaphore(1)
	for i := 1; i <= thread; i++ {
		seam.Add(1)
		rk := reckon.New()
		rk.HashMatrix = hashMatrix
		rk.TempInput = fmt.Sprintf("%s/GO_CNCI_file%v", outTemp, i)
		rk.TempScore = fmt.Sprintf("%s/GO_CNCI_file_score%v", outTemp, i)
		rk.TempDetil = fmt.Sprintf("%s/GO_CNCI_file_detil%v", outTemp, i)
		rk.Thread = thread
		go rk.Init(seam)
	}
	seam.Wait()
	Info("--------End of calculation-------")
	outfile := fmt.Sprintf("%s/pro", outDir)
	Info("---------Start merging files--------")
	scorePath := fmt.Sprintf("%s/GO_CNCI_score", outDir)
	detilPath := fmt.Sprintf("%s/GO_CNCI_detil", outDir)
	err = merge.Merge(outTemp, scorePath, detilPath, thread)
	if err != nil {
		Error("Merge err : [%s]", err.Error())
		return
	}
	scoreArray := ReadFileArray(scorePath)
	//scoreSLength := len(scoreArray) - 1
	//scoreArray = scoreArray[:scoreSLength]
	detilArray := ReadFileArray(detilPath)
	//detilSLength := len(detilArray) - 1
	//detilArray = detilArray[:detilSLength]
	err = merge.AddSvmLabel(scoreArray, outfile)
	if err != nil {
		Error("AddSvmLabel err : [%s]", err.Error())
		return
	}
	Info("---------End of merge file-------")
	SvmPutFileName := fmt.Sprintf("%s/svm", outDir)
	SvmFile := fmt.Sprintf("%s/file", outDir)
	SvmTmp := fmt.Sprintf("%s/tmp", outDir)
	Info("-------Start vector calculation------")
	err = Libsvm(outfile, SvmPutFileName, SvmFile, SvmTmp, libsvmPath, cnciParameters)
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
