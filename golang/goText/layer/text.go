package layer

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const dict string = "dicts/dict.txt"

const (
	Space   string = "#"
	Dspace         = "^"
	PreArgs        = "<args>" + Space
	PreDict        = "<dict>" + Space
	PreIP          = "<input>" + Space
	PreOP          = "<output>" + Space
	PreVec         = "<vector>" + Space

	maxCapacity int = 1024 * 1024 * 1024 * 1024
)

type Text struct {
	args_      *Args
	dict_      *Dictionary
	input_     *Matrix
	output_    *Matrix
	model_     *Model
	tokenCount uint64
	start      time.Time
}

func (text *Text) LoadModel(fp string) {
	file, err := os.Open(fp)
	if err != nil {
		fmt.Printf("Model file cannot be opened for loading! Error: %v", err.Error())
		return
	}
	defer file.Close()

	text.args_ = &Args{}
	text.dict_ = NewDic(text.args_)
	text.input_ = NewMat()
	text.output_ = NewMat()

	scanner := bufio.NewScanner(file)
	// set scanner buffer max
	buf := make([]byte, 0)
	scanner.Buffer(buf, maxCapacity)
	for scanner.Scan() {
		content := strings.TrimSpace(scanner.Text())
		// handle args
		if strings.Index(content, PreArgs) == 0 {
			fmt.Printf("Loading %v\n", PreArgs)
			strs := strings.Split(content, PreArgs)
			if len(strs) != 2 {
				fmt.Printf("invalid file: %v, args Error!", fp)
				return
			}
			text.args_.Load(strs[1])
			continue
		}
		// handle dict
		if strings.Index(content, PreDict) == 0 {
			strs := strings.Split(content, PreDict)
			fmt.Printf("Loading %v\n", PreDict)
			if len(strs) != 2 {
				fmt.Printf("invalid file: %v, dict Error!", fp)
				return
			}
			text.dict_.Load(strs[1])
			continue
		}
		// handle input
		if strings.Index(content, PreIP) == 0 {
			strs := strings.Split(content, PreIP)
			fmt.Printf("Loading %v\n", PreIP)
			if len(strs) != 2 {
				fmt.Printf("invalid file: %v, input Error!", fp)
				return
			}
			text.input_.Load(strs[1], PreIP)
			continue
		}
		// handle output
		if strings.Index(content, PreOP) == 0 {
			strs := strings.Split(content, PreOP)
			fmt.Printf("Loading %v\n", PreOP)
			if len(strs) != 2 {
				fmt.Printf("invalid file: %v, output Error!", fp)
				return
			}
			text.output_.Load(strs[1], PreOP)
			continue
		}
	}

	text.model_ = NewModel(text.input_, text.output_, text.args_, 0)
	if text.args_.model == Sup {
		text.model_.SetTargetCounts(text.dict_.GetCounts(labelType))
	} else {
		text.model_.SetTargetCounts(text.dict_.GetCounts(wordType))
	}
	fmt.Println("#############################Loading Done##############################")
}

func (text *Text) GetVector(vec *Vector, word string) {
	ngrams, err := text.dict_.GetNgramsByW(word)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}
	vec.Zero()
	for _, n := range ngrams {
		vec.AddRow(text.input_, n)
	}
	if len(ngrams) > 0 {
		vec.Mul(Real(1.0 / float64(len(ngrams))))
	}
}

func (text *Text) NgramVectors(word string) {
	var ngrams []uint64
	var substrings []string
	vec := NewVec()
	vec.Init(text.args_.dim)
	text.dict_.GetNgrams(word, &ngrams, &substrings)
	for i := 0; i < len(ngrams); i++ {
		vec.Zero()
		if ngrams[i] >= 0 {
			vec.AddRow(text.input_, ngrams[i])
		}
		fmt.Printf("%v %v", substrings[i], vec)
	}
}

func (text *Text) WordVectors() {
	input, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		fmt.Printf("input Error: %v", err.Error())
		return
	}
	vec := NewVec()
	vec.Init(text.args_.dim)
	text.GetVector(vec, input)
}

func (text *Text) TextVectors() {
	var line, labels []uint64
	vec := NewVec()
	vec.Init(text.args_.dim)

	input, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		fmt.Printf("input Error: %v", err.Error())
		return
	}
	text.dict_.GetLine(strings.TrimSpace(input), &line, &labels, text.model_.rng)
	text.dict_.AddNgrams(&line, text.args_.wordNgrams)
	for _, l := range line {
		vec.AddRow(text.input_, uint64(l))
	}
	if len(line) != 0 {
		vec.Mul(Real(1.0 / float64(len(line))))
	}
	fmt.Println(vec)
}

func (text *Text) PrintVectors() {
	if text.args_.model == Sup {
		text.TextVectors()
	} else {
		text.WordVectors()
	}
}

func (text *Text) trainThread(threadId uint64) {
	model := NewModel(text.input_, text.output_, text.args_, int64(threadId))
	// dict_.getCounts(): 返回的是一个词频数组 word/label
	if text.args_.model == Sup {
		model.SetTargetCounts(text.dict_.GetCounts(labelType))
	} else {
		model.SetTargetCounts(text.dict_.GetCounts(wordType))
	}

	// 训练数据实际行数
	ntokens, localTokenCount := text.dict_.Tokens(), uint64(0)
	var line, labels []uint64
	for text.tokenCount < uint64(text.args_.epoch)*ntokens/2 {
		progress := Real(float64(text.tokenCount) / float64(uint64(text.args_.epoch)*ntokens))
		lr := Real(text.args_.lr) * (Real(1) - progress)
		/**
		 * 每次循环初始化
		 * 1、重置文件流
		 * 2、重置words(line)、和labels
		 */
		line, labels = line[0:0], labels[0:0]
		file, _ := os.Open(text.args_.input)
		scanner, lines := bufio.NewScanner(file), 0
		for scanner.Scan() {
			// 在这里线程中的数据才训练
			lines++
			if uint64(lines) < text.dict_.lineCount_*threadId/uint64(text.args_.thread) || uint64(lines) >= text.dict_.lineCount_*(threadId+1)/uint64(text.args_.thread) {
				continue
			}
			localTokenCount += text.dict_.GetLine(strings.TrimSpace(scanner.Text()), &line, &labels, model.rng)
		}
		if text.args_.model == Sup {
			text.dict_.AddNgrams(&line, text.args_.wordNgrams)
			// 按行来，一行只有一个label，如果是多个就随机选其中一个
			supervised(model, lr, line, labels)
		} else if text.args_.model == Cbow {
			text.cbow(model, lr, line)
		} else if text.args_.model == Sg {
			text.skipgram(model, lr, line)
		}

		if localTokenCount > text.args_.lrUpdateRate {
			text.tokenCount += localTokenCount
			localTokenCount = uint64(0)
		}
		if threadId == 0 {
			fmt.Printf("tokenCount: %v\n", text.tokenCount)
		}
	}
}

func supervised(model *Model, lr Real, line, labels []uint64) {
	if len(labels) == 0 || len(line) == 0 {
		fmt.Printf("labels : %v or line: %v should not be null\n", labels, line)
		return
	}
	i := model.rng.Intn(len(labels) - 1)
	model.Update(line, labels[i], lr)
}

func (text *Text) cbow(model *Model, lr Real, line []uint64) {
	if model.args_.ws <= 0 {
		fmt.Printf("invalid ws: %v", model.args_.ws)
		return
	}
	var bow []uint64
	for w := 0; w < len(line); w++ {
		boundary := model.rng.Intn(text.args_.ws-1) + 1
		bow = bow[0:0]
		for c := -boundary; c <= boundary; c++ {
			// 不能是词本身，位置不能是句子负数，不能超过句子长度
			if c != 0 && w+c >= 0 && w+c < len(line) {
				// 实际被加入 input 的不止是词本身，还有词的 word n-gram
				ngrams, e := text.dict_.GetNgramsByI(line[w+c])
				if e != nil {
					fmt.Println(e.Error())
					return
				}
				bow = append(bow, ngrams...)
			}
		}
		model.Update(bow, line[w], lr)
	}
}

func (text *Text) skipgram(model *Model, lr Real, line []uint64) {
	if model.args_.ws <= 0 {
		fmt.Printf("invalid ws: %v", model.args_.ws)
		return
	}
	for w := 0; w < len(line); w++ {
		boundary := model.rng.Intn(text.args_.ws-1) + 1
		ngrams, e := text.dict_.GetNgramsByI(line[w])
		if e != nil {
			fmt.Println(e.Error())
			return
		}
		// 对上下文的每一个词分别更新一次模型
		for c := -boundary; c <= boundary; c++ {
			if c != 0 && w+c >= 0 && w+c < len(line) {
				model.Update(ngrams, line[w+c], lr)
			}
		}
	}
}

func (text *Text) Test(fp string, k uint) {
	file, err := os.Open(text.args_.input)
	if err != nil {
		fmt.Printf("Error: can not open text file from: %v! Error: %v", fp, err.Error())
		return
	}
	defer file.Close()

	var nexamples, nlabels uint
	var precision float32
	line, labels := make([]uint64, 0), make([]uint64, 0)

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		if nexamples%10000 == 0 {
			runtime.GC()
		}
		content := strings.TrimSpace(scanner.Text())
		text.dict_.GetLine(content, &line, &labels, text.model_.rng)
		text.dict_.AddNgrams(&line, text.args_.wordNgrams)
		if len(line) > 0 && len(labels) > 0 {
			text.model_.Predict(k, line, text.model_.hidden_, text.model_.output_)
			for _, item := range text.model_.pq {
				for _, l := range labels {
					if item.index == l {
						precision += 1.0
						continue
					}
				}
			}
			nexamples += 1
			nlabels += uint(len(labels))
		}
	}
	fmt.Printf("P@%v: %v\n", k, precision/float32(k*nexamples))
	fmt.Printf("R@%v: %v", k, precision/float32(nlabels))
	fmt.Printf("Number of examples: %v", nexamples)
}

type itemstr struct {
	r     Real
	label string
}

func (text *Text) Predict(content string, k uint, print_prob bool) string {
	words, labels := make([]uint64, 0), make([]uint64, 0)
	text.dict_.GetLine(strings.TrimSpace(content), &words, &labels, text.model_.rng)
	text.dict_.AddNgrams(&words, text.args_.wordNgrams)
	if len(words) == 0 {
		fmt.Printf("length of words == 0: %v\n", words)
		return ""
	}
	hidden, output := NewVec(), NewVec()
	hidden.Init(text.args_.dim)
	output.Init(text.dict_.Labels())

	text.model_.Predict(k, words, hidden, output)
	if len(text.model_.pq) == 0 {
		fmt.Println("n/a")
		return ""
	}
	var result string
	for _, item := range text.model_.pq {
		word := text.dict_.GetWord(item.index)
		result += fmt.Sprint(word)
		if print_prob {
			result += fmt.Sprint("/", item.r)
		}
		result += "\n"
	}
	return result
}

func (text *Text) saveModel() {
	file, e := os.Create(text.args_.output + ".bin")
	if e != nil {
		fmt.Printf("Model file: %v cannot be opened for saving! Error: %v.bin\n", text.args_.output, e.Error())
		return
	}
	defer file.Close()
	text.args_.Save(file)
	text.dict_.Save(file)
	text.input_.Save(file, PreIP)
	text.output_.Save(file, PreOP)
}

func (text *Text) getVector(vec *Vector, word string) {
	ngrams, e := text.dict_.GetNgramsByW(word)
	if e != nil {
		fmt.Printf("get ngrams error: %v\n", e.Error())
		return
	}
	vec.Zero()
	/**
	 * 首先，vec长度 等于 len(ngrams)
	 * 然后，遍历vec中每个成员
	 * 最后，将ngrams中word对应位置的input值累加到vec上
	 */
	for _, n := range ngrams {
		vec.AddRow(text.input_, n)
	}
	// vec上每个成员值除以长度
	if len(ngrams) > 0 {
		vec.Mul(Real(1 / float64(len(ngrams))))
	}
}

func (text *Text) saveOutput() {
	file, e := os.Create(text.args_.output + ".output")
	if e != nil {
		fmt.Printf("OutPut file: %v cannot be opened for saving! Error: %v.output\n", text.args_.output, e.Error())
		return
	}
	defer file.Close()

	var str string = PreOP + strconv.FormatUint(uint64(text.dict_.Words()), 10) + Space + strconv.FormatUint(text.args_.dim, 10) + "\n"
	if _, e = file.WriteString(str); e != nil {
		fmt.Println(e.Error())
		return
	}

	vec := NewVec()
	vec.Init(text.args_.dim)
	for i := uint64(0); i < text.dict_.Words(); i++ {
		word := text.dict_.GetWord(i)
		vec.Zero()
		if e := vec.AddRow(text.output_, i); e != nil {
			fmt.Printf("Save output file Error: %v", e.Error())
			return
		}
		str = word + Space
		for _, d := range vec.data_ {
			str += strconv.FormatFloat(float64(d), 'E', -1, 64) + Dspace
		}
		file.WriteString(str + "\n")
	}
}

func (text *Text) saveVectors() {
	file, e := os.Create(text.args_.output + ".vec")
	if e != nil {
		fmt.Printf("Vector file: %v cannot be opened for saving! Error: %v.vec\n", text.args_.output, e.Error())
		return
	}
	defer file.Close()

	var str string = PreVec + strconv.FormatUint(uint64(text.dict_.Words()), 10) + Space + strconv.FormatUint(text.args_.dim, 10) + "\n"
	if _, e = file.WriteString(str); e != nil {
		fmt.Println(e.Error())
		return
	}

	vec := NewVec()
	vec.Init(text.args_.dim)
	for i := uint64(0); i < text.dict_.Words(); i++ {
		word := text.dict_.GetWord(i)
		text.getVector(vec, word)
		str = word + Space
		for _, d := range vec.data_ {
			str += strconv.FormatFloat(float64(d), 'E', -1, 64) + Dspace
		}
		file.WriteString(str + "\n")
	}
}

func (text *Text) LoadVectors(fp string) {
	file, err := os.Open(fp)
	if err != nil {
		fmt.Printf("Pretrained vectors file cannot be opened!: %v", err.Error())
		return
	}
	defer file.Close()

	var words []string
	mat := NewMat()
	var n, dim uint64

	scanner := bufio.NewScanner(file)
	// handle first line
	for scanner.Scan() {
		content := strings.TrimSpace(scanner.Text())
		if strings.Index(content, PreVec) == 0 {
			strs := strings.Split(content, PreVec)
			if len(strs) != 2 {
				fmt.Printf("invalid vectors file: %v\n", content)
				return
			}
			strs = strings.Split(strs[1], PreVec)
			if len(strs) != 2 {
				fmt.Printf("invalid vectors file: %v\n", content)
				return
			}

			if s, err := strconv.ParseUint(strs[0], 10, 64); err == nil {
				n = s
			} else {
				fmt.Println(err)
				return
			}

			if s, err := strconv.ParseUint(strs[1], 10, 64); err == nil {
				dim = s
			} else {
				fmt.Println(err)
				return
			}

			if dim != text.args_.dim {
				fmt.Println("Dimension of pretrained vectors does not match -dim option")
				return
			}
			break
		}
	}
	mat.Init(n, dim)
	var line uint64 = 0

	for scanner.Scan() {
		content := strings.TrimSpace(scanner.Text())
		if strings.Index(content, PreVec) == 0 {
			continue
		}

		if line%100000 == 0 {
			runtime.GC()
		}

		strs := strings.Split(content, Space)
		if len(strs) != 2 {
			fmt.Printf("invalid vectors line: %v\n", content)
			continue
		}
		words = append(words, strs[0])
		text.dict_.Add(strs[0])

		strs = strings.Split(strs[1], Dspace)
		if uint64(len(strs)) != dim {
			fmt.Printf("invalid vectors line: %v\n", content)
			continue
		}
		for i, str := range strs {
			s, e := strconv.ParseFloat(str, 32)
			if e != nil {
				fmt.Printf("Error: %v", e)
				continue
			}
			text.input_.data_[line*dim+uint64(i)] = Real(float32(s))
		}
		line += 1
	}
	text.dict_.Threshold(1, 0)
	text.input_ = NewMat()
	text.input_.Init(text.dict_.Words()+text.args_.bucket, text.args_.dim)
	text.input_.Uniform(Real(1.0 / float64(text.args_.dim)))

	for i := uint64(0); i < n; i++ {
		idx := text.dict_.getId(words[i])
		if idx < 0 || idx >= text.dict_.Words() {
			continue
		}
		for j := uint64(0); j < dim; j++ {
			text.input_.data_[idx*dim+uint64(j)] = mat.data_[i*dim+uint64(j)]
		}
	}
}

/**
 * step 1: 读取文本，分词 -> 添加词 -> word(string) => word hash(uint64) => word2int[h](int) => words_[int](type) -> 传给text
 * step 2: 初始化agrs
 * step 3: 初始化input, (words长度 + bucket) * dim(vec长度), 激活函数： 随机数~[-1, 1] * r
 * step 4: 初始化output, 输出len(labels) * vec/输出len(words) * vec
 * step 5: 训练
 * step 6: 保存数据
 */
func Train(args *Args) {
	file, err := os.Open(args.input)
	if err != nil {
		fmt.Printf("Error: can not open input file from: %v", err.Error())
		return
	}
	defer file.Close()
	text := new(Text)
	text.args_, text.dict_ = args, NewDic(args)

	err = text.dict_.ReadFromFile(file)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if text.args_.pretrainedVectors != "" {
		// TODO loadVectors
	} else {
		text.input_ = NewMat()
		text.input_.Init(uint64(text.dict_.Words()+text.args_.bucket), uint64(text.args_.dim))
		text.input_.Uniform(Real(1.0 / float64(text.args_.dim)))
	}

	if text.args_.model == Sup {
		text.output_ = NewMat()
		text.output_.Init(uint64(text.dict_.Labels()), uint64(text.args_.dim)) // 输出len(labels) * vec
	} else {
		text.output_ = NewMat()
		text.output_.Init(uint64(text.dict_.Words()), uint64(text.args_.dim)) // 输出len(words) * vec
	}

	text.start = time.Now()
	text.tokenCount = 0
	thread := text.args_.thread
	ch := make(chan bool, thread)
	fmt.Println("###################Start training......###################")
	trainT := time.Now()
	doTrain := func(i uint64) {
		text.trainThread(i)
		ch <- true
	}
	for i := uint64(0); i < thread; i++ {
		go doTrain(i)
	}
	for i := uint64(0); i < thread; i++ {
		<-ch
	}
	fmt.Printf("###################training ended, time cost: %v......###################\n", time.Now().Sub(trainT).String())

	fmt.Println("###################Saving Model###################")
	text.model_ = NewModel(text.input_, text.output_, text.args_, 0)
	text.saveModel()
	if text.args_.model != Sup {
		text.saveVectors()
		if text.args_.saveOutput > 0 {
			text.saveOutput()
		}
	}
}
