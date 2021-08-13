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
	hashMatrix := ReadFileMatrix(cnciParameters + "/CNCI_matrix")
	sequenceArr := ReadFileArray(inputFile)
	sLen := len(sequenceArr) - 1
	sequenceArr = sequenceArr[:sLen]

	fastArray := TwoLineFasta(sequenceArr)
	if len(fastArray) < thread {
		Info("Please set a smaller number of threads")
		return
	}

	Info("-------Start splitting file------")
	in := SplitFile(fastArray, thread)
	Info("--------End of split file-------")
	Info("--------Start calculation-------")
	seam := gsema.NewSemaphore(thread)
	for i := 1; i <= thread; i++ {
		seam.Add(1)
		rk := reckon.New()
		rk.HashMatrix = hashMatrix
		rk.FileInput, _ = in.Load(i)
		rk.Thread = thread
		go rk.Init(seam)
	}
	seam.Wait()
	Info("--------End of calculation-------")
	outfile := fmt.Sprintf("%s/pro", outDir)
	Info("---------Start merging--------")
	err = merge.AddSvmLabel(reckon.OS_PROPERTY, outfile)
	if err != nil {
		Error("AddSvmLabel err : [%s]", err.Error())
		return
	}
	Info("---------End of merge-------")
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
	FirResult := PutResult(reckon.OS_DETIL, SvmFile)
	SvmFinalResult := fmt.Sprintf("%s/GO_CNCI.index", outDir)
	PrintResult(FirResult, SvmFinalResult)
	Info("---------End of output file----------")
	cost := time.Since(start)
	Info("Time use [%s]", cost)
	_ = os.Remove(SvmPutFileName)
	_ = os.Remove(SvmFile)
	_ = os.Remove(SvmTmp)
	_ = os.Remove(outfile)
}
