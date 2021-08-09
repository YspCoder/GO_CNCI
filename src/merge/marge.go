/**
 * @Author: lipengfei
 * @Description:
 * @File:  marge
 * @Version: 1.0.0
 * @Date: 2021/08/06 16:39
 */
package merge

import (
	. "GO_CNCI/src/base"
	"fmt"
	"io/ioutil"
	"os"
)

func Merge(inputDir, outDir string) error {
	score, err := os.Create(outDir + "/CNCI_score")
	if err != nil {
		Error("Merge Create outDir Err : [%v]", err)
		return err
	}
	for i := 1; i <= 50; i++ {
		Score_string := fmt.Sprintf("%s/CNCI_file_score%v", inputDir, i)
		f, err := os.OpenFile(Score_string, os.O_RDONLY, os.ModePerm)
		if err != nil {
			Error("Merge OpenFile Err : [%v]", err.Error())
			continue
		}
		b, err := ioutil.ReadAll(f)
		if err != nil {
			Error("Merge ReadAll Err : [%v]", err.Error())
			continue
		}
		_, _ = score.Write(b)
		_ = f.Close()
	}

	detil, err := os.Create(outDir + "CNCI_detil")
	if err != nil {
		Error("Merge Create outDir Err : [%v]", err)
		return err
	}
	for i := 1; i <= 50; i++ {
		Detil_string := fmt.Sprintf("%s/CNCI_file_detil%v", inputDir, i)
		f, err := os.OpenFile(Detil_string, os.O_RDONLY, os.ModePerm)
		if err != nil {
			Error("Merge OpenFile Err : [%v]", err.Error())
			continue
		}
		b, err := ioutil.ReadAll(f)
		if err != nil {
			Error("Merge ReadAll Err : [%v]", err.Error())
			continue
		}
		_, _ = detil.Write(b)
		_ = f.Close()
	}
	return nil
}
