# GO_CNCI
## About GO_CNCI
 CPU resource consumption 100%+

## Documentation
```
The input file is gffread output
number_of_file_partitions propose 1-30
thread propose 4-12
 ./GO_CNCI reference_folder inputFile number_of_file_partitions outDir libsvmpath thread
```
## Install GO_CNCI
```
 git clone https://github.com/yingfeikeji/GO_CNCI.git
 unzip libsvm-3.0.zip
 cd libsvm-3.0
 make
 ----------------------------------------------------
 cd ..
 cd src/main
 go build main.go
```

## Example
```
 ./GO_CNCI ./CNCI_Parameters ./94d6346_candidate.fa ./test ./libsvm 8
```
## Citation
```
 https://github.com/www-bioinfo-org/CNCI.git
```
