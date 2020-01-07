package layer

// import (
// 	"fmt"
// 	"os"
// 	"testing"
// )

// func Test_saveArgs(t *testing.T) {
// 	fmt.Println("======================Test_saveArgs=======================")
// 	a := InitArgs(Cbow, Hs, "../testdata/input.txt", "../testdata/output.txt", "__label__", "", 0.05, 1e-4, 100, 100, 5, 1, 0, 1, 2, 12, 2, 0, 20000, 5, 1, 5)
// 	file, e := os.Create(a.output)
// 	if e != nil {
// 		fmt.Printf("Model file: %v cannot be opened for saving! Error: %v.bin\n", a.output, e.Error())
// 		return
// 	}
// 	defer file.Close()
// 	a.Save(file)
// }

// func Test_loadArgs(t *testing.T) {
// 	fmt.Println("======================Test_loadArgs=======================")
// 	a := new(Args)
// 	a.Load("<args>###cbow###hs###5###5###1###100###5###100###20000###1###2###0###1E-04")
// 	fmt.Printf("model: %v\n", a.model)
// 	fmt.Printf("input: %v\n", a.input)
// 	fmt.Printf("output: %v\n", a.output)
// 	fmt.Printf("loss: %v\n", a.loss)
// 	fmt.Printf("lr: %v\n", a.lr)
// 	fmt.Printf("lrUpdateRate: %v\n", a.lrUpdateRate)
// 	fmt.Printf("dim: %v\n", a.dim)
// 	fmt.Printf("ws : %v\n", a.ws)
// 	fmt.Printf("epoch: %v\n", a.epoch)
// 	fmt.Printf("minCount: %v\n", a.minCount)
// 	fmt.Printf("minCountLabel: %v\n", a.minCountLabel)
// 	fmt.Printf("neg : %v\n", a.neg)
// 	fmt.Printf("bucket: %v\n", a.bucket)
// 	fmt.Printf("minn: %v\n", a.minn)
// 	fmt.Printf("maxn: %v\n", a.maxn)
// 	fmt.Printf("thread: %v\n", a.thread)
// 	fmt.Printf("t: %v\n", a.t)
// 	fmt.Printf("label: %v\n", a.label)
// 	fmt.Printf("verbose: %v\n", a.verbose)
// 	fmt.Printf("pretrainedVectors: %v\n", a.pretrainedVectors)
// 	fmt.Printf("saveOutput: %v\n", a.saveOutput)
// 	fmt.Printf("wordNgrams: %v\n", a.wordNgrams)
// }
