package layer

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	Cbow string = "cbow"
	Sg          = "sg"
	Sup         = "sup"

	Hs      = "hs"
	Ns      = "ns"
	Softmax = "softmax"
)

type Args struct {
	model  string
	input  string
	output string

	loss              string
	lr                float64
	lrUpdateRate      uint64
	dim               uint64
	ws                int
	epoch             uint
	minCount          uint64
	minCountLabel     uint64
	neg               uint64
	bucket            uint64
	minn              uint64
	maxn              uint64
	thread            uint64
	t                 float64
	label             string
	verbose           uint64
	pretrainedVectors string
	saveOutput        uint64
	wordNgrams        uint
}

func InitArgs(model, loss, input, output, label, pretrainedVectors string, lr, t float64, lrUpdateRate, dim, minCount, minCountLabel, neg, minn, maxn, thread, verbose, saveOutput, bucket uint64, epoch, wordNgrams uint, ws int) *Args {
	return &Args{
		model:             model,
		input:             input,
		output:            output,
		loss:              loss,
		lr:                lr,
		lrUpdateRate:      lrUpdateRate,
		dim:               dim,
		minCount:          minCount,
		minCountLabel:     minCountLabel,
		neg:               neg,
		bucket:            bucket,
		minn:              minn,
		maxn:              maxn,
		thread:            thread,
		t:                 t,
		label:             label,
		verbose:           verbose,
		pretrainedVectors: pretrainedVectors,
		saveOutput:        saveOutput,
		epoch:             epoch,
		wordNgrams:        wordNgrams,
		ws:                ws,
	}
}

func (a *Args) Save(file *os.File) {
	var e error
	var str string = PreArgs + a.model + Space + a.loss + Space + strconv.FormatInt(int64(a.ws), 10) + Space + strconv.FormatUint(uint64(a.epoch), 10) + Space + strconv.FormatUint(uint64(a.wordNgrams), 10) + Space + strconv.FormatUint(a.dim, 10) + Space + strconv.FormatUint(a.minCount, 10) + Space + strconv.FormatUint(a.lrUpdateRate, 10) + Space + strconv.FormatUint(a.bucket, 10) + Space + strconv.FormatUint(a.minn, 10) + Space + strconv.FormatUint(a.maxn, 10) + Space + strconv.FormatUint(a.neg, 10) + Space + strconv.FormatFloat(a.t, 'E', -1, 64) + "\n"

	if _, e = file.WriteString(str); e != nil {
		fmt.Println(e.Error())
		return
	}
}

func (a *Args) Load(content string) {
	args := strings.Split(strings.TrimSpace(content), Space)
	if len(args) != 13 {
		fmt.Printf("Error! invalid args: %v, length: %v\n", args, len(args))
		return
	}
	var e error

	//string
	a.model = args[0]
	a.loss = args[1]

	// int
	ws, e := strconv.ParseInt(args[2], 10, 0)
	if e != nil {
		fmt.Println(e.Error())
		return
	}
	a.ws = int(ws)

	// uint
	epoch, e := strconv.ParseUint(args[3], 10, 0)
	if e != nil {
		fmt.Println(e.Error())
		return
	}
	a.epoch = uint(epoch)
	wordNgrams, e := strconv.ParseUint(args[4], 10, 0)
	if e != nil {
		fmt.Println(e.Error())
		return
	}
	a.wordNgrams = uint(wordNgrams)

	// uint64
	dim, e := strconv.ParseUint(args[5], 10, 64)
	if e != nil {
		fmt.Println(e.Error())
		return
	}
	a.dim = dim
	minCount, e := strconv.ParseUint(args[6], 10, 64)
	if e != nil {
		fmt.Println(e.Error())
		return
	}
	a.minCount = minCount
	lrUpdateRate, e := strconv.ParseUint(args[7], 10, 64)
	if e != nil {
		fmt.Println(e.Error())
		return
	}
	a.lrUpdateRate = lrUpdateRate
	bucket, e := strconv.ParseUint(args[8], 10, 64)
	if e != nil {
		fmt.Println(e.Error())
		return
	}
	a.bucket = bucket
	minn, e := strconv.ParseUint(args[9], 10, 64)
	if e != nil {
		fmt.Println(e.Error())
		return
	}
	a.minn = minn
	maxn, e := strconv.ParseUint(args[10], 10, 64)
	if e != nil {
		fmt.Println(e.Error())
		return
	}
	a.maxn = maxn
	neg, e := strconv.ParseUint(args[11], 10, 64)
	if e != nil {
		fmt.Println(e.Error())
		return
	}
	a.neg = neg

	// float64
	t, e := strconv.ParseFloat(args[12], 64)
	if e != nil {
		fmt.Println(e.Error())
		return
	}
	a.t = t
}
