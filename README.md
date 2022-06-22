# GO_CNCI
## About GO_CNCI
```
 CPU resource consumption 100%+
 Data fault tolerance 99.96%
```

## Documentation
```
The input file is gffread output
thread propose 4-12
This tool will execute in hyper thread (t*6)
assign the classification models ("ve" for vertebrate species, "pl" for plat species)
 ./GO_CNCI reference_folder inputFile outDir libsvmpath model thread
```
## Install GO_CNCI
```
 cd GO_CNCI
 git clone https://github.com/YspCoder/GO_CNCI.git
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
 ./GO_CNCI ./CNCI_Parameters ./94d6346_candidate.fa ./test ./libsvm ve 8
```
## Citation
```
 https://github.com/www-bioinfo-org/CNCI.git
```
