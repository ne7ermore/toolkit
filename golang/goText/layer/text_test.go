package layer

import (
	"fmt"
	"runtime"
	"testing"
)

func Test_text(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	// a := InitArgs(Cbow, Softmax, "../testdata/input2.txt", "../testdata/CbowHs", "__label__", "", 0.025, 1e-4, 100, 100, 5, 1, 5, 1, 2, 16, 2, 0, 500000, 1, 2, 5)
	// fmt.Printf("model: %v\n", a.model)
	// fmt.Printf("input: %v\n", a.input)
	// fmt.Printf("output: %v\n", a.output)
	// fmt.Printf("loss: %v\n", a.loss)
	// fmt.Printf("lr: %v\n", a.lr)
	// fmt.Printf("lrUpdateRate: %v\n", a.lrUpdateRate)
	// fmt.Printf("dim: %v\n", a.dim)
	// fmt.Printf("ws : %v\n", a.ws)
	// fmt.Printf("epoch: %v\n", a.epoch)
	// fmt.Printf("minCount: %v\n", a.minCount)
	// fmt.Printf("minCountLabel: %v\n", a.minCountLabel)
	// fmt.Printf("neg : %v\n", a.neg)
	// fmt.Printf("bucket: %v\n", a.bucket)
	// fmt.Printf("minn: %v\n", a.minn)
	// fmt.Printf("maxn: %v\n", a.maxn)
	// fmt.Printf("thread: %v\n", a.thread)
	// fmt.Printf("t: %v\n", a.t)
	// fmt.Printf("label: %v\n", a.label)
	// fmt.Printf("verbose: %v\n", a.verbose)
	// fmt.Printf("pretrainedVectors: %v\n", a.pretrainedVectors)
	// fmt.Printf("saveOutput: %v\n", a.saveOutput)
	// fmt.Printf("wordNgrams: %v\n", a.wordNgrams)
	// Train(a)
	gotext := new(Text)
	gotext.LoadModel("../testdata/CbowHs.bin")
	res := gotext.Predict("波司登", 7, true)
	fmt.Println("input:")
	fmt.Println("波司登")
	fmt.Println("output:")
	fmt.Println(res)
}
