package main

import (
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
	if len(os.Args) < 6 {
		fmt.Println("Insufficient required parameters")
		fmt.Println("./GO_CNCI reference_folder inputFile outDir libsvmpath model thread")
		fmt.Println("./GO_CNCI ./CNCI_Parameters ./94d6346_candidate.fa ./test ./libsvm ve 8")
		return
	}
	cnciParameters := os.Args[1]
	inputFile := os.Args[2]
	outDir := os.Args[3]
	libsvmPath := os.Args[4]
	classModel := os.Args[5]
	thread, err := strconv.Atoi(os.Args[6])
	if err != nil {
		fmt.Println("Please enter a positive integer -- thread")
		return
	}
	hashMatrix := ReadFileMatrix(cnciParameters + "/GO_CNCI_matrix")
	sequenceArr := ReadFileArray(inputFile)
	sLen := len(sequenceArr)
	if sequenceArr[sLen-1] == "" {
		sequenceArr = sequenceArr[:sLen-1]
	}
	fastArray := TwoLineFasta(sequenceArr)
	if len(fastArray) < thread {
		fmt.Println("Please set a smaller number of threads")
		return
	}

	fmt.Println("-------Start splitting file------")
	in := SplitFile(fastArray, thread)
	fmt.Println("--------End of split file-------")
	fmt.Println("--------Start calculation-------")
	var wg sync.WaitGroup
	for i := 1; i <= thread; i++ {
		wg.Add(1)
		rk := reckon.New()
		rk.HashMatrix = hashMatrix
		rk.FileInput, _ = in.Load(i)
		rk.Thread = thread
		go rk.Init(&wg)
	}
	wg.Wait()
	fmt.Println("--------End of calculation-------")
	outfile := fmt.Sprintf("%s/pro", outDir)
	fmt.Println("---------Start merging--------")
	err = merge.AddSvmLabel(reckon.OS_PROPERTY, outfile)
	if err != nil {
		fmt.Printf("AddSvmLabel err : [%s]", err.Error())
		return
	}
	fmt.Println("---------End of merge-------")
	SvmPutFileName := fmt.Sprintf("%s/svm", outDir)
	SvmFile := fmt.Sprintf("%s/file", outDir)
	SvmTmp := fmt.Sprintf("%s/tmp", outDir)
	fmt.Println("-------Start vector calculation------")
	err = Libsvm(outfile, SvmPutFileName, SvmFile, SvmTmp, libsvmPath, cnciParameters, classModel)
	if err != nil {
		fmt.Printf("Libsvm err : [%s]", err.Error())
		return
	}
	fmt.Println("----------End of vector calculation--------")
	fmt.Println("Start output file")
	FirResult := PutResult(reckon.OS_DETIL, SvmFile)
	SvmFinalResult := fmt.Sprintf("%s/GO_CNCI.index", outDir)
	PrintResult(FirResult, SvmFinalResult)
	fmt.Println("---------End of output file----------")
	cost := time.Since(start)
	fmt.Printf("Time use [%s]", cost)
	_ = os.Remove(SvmPutFileName)
	_ = os.Remove(SvmFile)
	_ = os.Remove(SvmTmp)
	_ = os.Remove(outfile)
}
