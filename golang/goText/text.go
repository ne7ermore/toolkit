package main

import (
	"flag"
	"fmt"
	"runtime"

	"goText/layer"
)

var (
	// mandatory arguments
	command  = flag.String("command", "cbow", "goText command")
	input    = flag.String("input", "", "training file path")
	output   = flag.String("output", "", "output file path")
	preInput = flag.String("preInput", "", "preInput file path")

	// optional arguments
	lr                = flag.Float64("lr", 0.025, "learning rate")
	lrUpdateRate      = flag.Uint64("lrUpdateRate", 100, "change the rate of updates for the learning rate")
	dim               = flag.Uint64("dim", 100, "size of word vectors")
	minCount          = flag.Uint64("minCount", 5, "minimal number of word occurences")
	minCountLabel     = flag.Uint64("minCountLabel", 0, "minimal number of label occurences")
	neg               = flag.Uint64("neg", 5, "number of negatives sampled")
	loss              = flag.String("loss", "softmax", "loss function {ns, hs, softmax} [ns]")
	bucket            = flag.Uint64("bucket", 500000, "number of buckets")
	minn              = flag.Uint64("minn", 1, "min length of ngram")
	maxn              = flag.Uint64("maxn", 2, "max length of ngram")
	thread            = flag.Uint64("thread", 8, "number of threads")
	t                 = flag.Float64("t", 1e-4, "sampling threshold")
	label             = flag.String("label", "__label__", "labels prefix")
	verbose           = flag.Uint64("verbose", 2, "verbosity level")
	pretrainedVectors = flag.String("pretrainedVectors", "", "pretrained word vectors for supervised learning")
	saveOutput        = flag.Uint64("saveOutput", 0, "whether output params should be saved")
	wordNgrams        = flag.Int("wordNgrams", 2, "max length of word ngram")
	epoch             = flag.Int("epoch", 1, "number of epochs")
	ws                = flag.Int("ws", 5, "size of the context window")

	testdata = flag.String("testdata", "", "test data filename (if -, read from stdin)")
	k        = flag.Int("k", 1, "(optional; 1 by default) predict top k labels")
)

func printUsage() {
	fmt.Println("usage: goText <command> <args>\n")
	fmt.Println("The commands supported by goText are:\n")
	fmt.Println("supervised          train a supervised classifier")
	fmt.Println("test                evaluate a supervised classifier")
	fmt.Println("predict             predict most likely labels")
	fmt.Println("predict-prob        predict most likely labels with probabilities")
	fmt.Println("skipgram            train a skipgram model")
	fmt.Println("cbow                train a cbow model")
	fmt.Println("print-vectors       print vectors given a trained model")
}

func printTestUsage() {
	fmt.Println("usage: goText test <model> <test-data> [<k>]\n")
	fmt.Println("  <input model>      model filename")
	fmt.Println("  <test-data>  	  test data filename (if -, read from stdin)")
	fmt.Println("  <k>          	  (optional; 1 by default) predict top k labels")
}

func printPredictUsage() {
	fmt.Println("usage: goText predict[-prob] <model> <test-data> [<k>]\n")
	fmt.Println("  <input model>      model filename")
	fmt.Println("  <preInput>  		  pre data")
	fmt.Println("  <k>          	  (optional; 1 by default) predict top k labels")
}

func printPrintVectorsUsage() {
	fmt.Println("usage: goText print-vectors <model>\n")
	fmt.Println("  <model>      	  model filename")
}

func printPrintNgramsUsage() {
	fmt.Println("usage: goText print-ngrams <model> <word>\n")
	fmt.Println("  <model>      	  model filename")
	fmt.Println("  <word>       	  word to print")
}

func printHelp() {
	fmt.Printf("The following arguments are mandatory:\n")
	fmt.Printf("  -input              training file path\n")
	fmt.Printf("  -output             output file path\n\n")
	fmt.Printf("The following arguments are optional:\n")
	fmt.Printf("  -lr                 learning rate ['%v']\n", *lr)
	fmt.Printf("  -lrUpdateRate       change the rate of updates for the learning rate ['%v']\n", *lrUpdateRate)
	fmt.Printf("  -dim                size of word vectors ['%v']\n", *dim)
	fmt.Printf("  -ws                 size of the context window ['%v']\n", *ws)
	fmt.Printf("  -epoch              number of epochs ['%v']\n", *epoch)
	fmt.Printf("  -minCount           minimal number of word occurences ['%v']\n", *minCount)
	fmt.Printf("  -minCountLabel      minimal number of label occurences ['%v']\n", *minCountLabel)
	fmt.Printf("  -neg                number of negatives sampled ['%v']\n", *neg)
	fmt.Printf("  -wordNgrams         max length of word ngram ['%v']\n", *wordNgrams)
	fmt.Printf("  -loss               loss function {ns, hs, softmax} [ns]\n")
	fmt.Printf("  -bucket             number of buckets ['%v']\n", *bucket)
	fmt.Printf("  -minn               min length of char ngram ['%v']\n", *minn)
	fmt.Printf("  -maxn               max length of char ngram ['%v']\n", *maxn)
	fmt.Printf("  -thread             number of threads ['%v']\n", *thread)
	fmt.Printf("  -t                  sampling threshold ['%v']\n", *t)
	fmt.Printf("  -label              labels prefix ['%v']\n", *label)
	fmt.Printf("  -verbose            verbosity level ['%v']\n", *verbose)
	fmt.Printf("  -pretrainedVectors  pretrained word vectors for supervised learning []\n")
	fmt.Printf("  -saveOutput         whether output params should be saved ['%v']\n", *saveOutput)
}

func train() {
	var model string
	if *command == "skipgram" {
		model = layer.Sg
	} else if *command == "cbow" {
		model = layer.Cbow
	} else if *command == "supervised" {
		model = layer.Sup
	} else {
		printUsage()
		return
	}
	if *input == "" || *output == "" {
		printHelp()
		return
	}

	args := layer.InitArgs(model, *loss, *input, *output, *label, *pretrainedVectors, *lr, *t, *lrUpdateRate, *dim, *minCount, *minCountLabel, *neg, *minn, *maxn, *thread, *verbose, *saveOutput, *bucket, uint(*epoch), uint(*wordNgrams), *ws)
	layer.Train(args)
}

// 只按文件测试
func test() {
	if *input == "" || *testdata == "" {
		printTestUsage()
	}
	gotext := new(layer.Text)
	gotext.LoadModel(*input)
	gotext.Test(*testdata, uint(*k))
}

func predict() {
	if *input == "" || *preInput == "" {
		printTestUsage()
	}
	print_prob := *command == "predict-prob"
	gotext := new(layer.Text)
	gotext.LoadModel(*input)
	gotext.Predict(*preInput, uint(*k), print_prob)
}

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())
	if *command == "skipgram" || *command == "cbow" || *command == "supervised" {
		train()
	} else if *command == "test" {
		test()
	} else if *command == "predict" || *command == "predict-prob" {
		predict()
	} else {
		printUsage()
	}
}
