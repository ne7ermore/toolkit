package layer

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	"goText/common"
)

type Entry_type int

const (
	wordType Entry_type = iota + 0
	labelType

	MAX_VOCAB_SIZE uint64 = 30000000 // 语料库长度上限
	// MAX_VOCAB_SIZE uint64 = 220  // 语料库长度上限
	MAX_LINE_SIZE = 1024 // 一行文本最大长度

	BOW string = "<"
	EOW        = ">"
)

type Entry struct {
	word     string
	count    uint64
	wType    Entry_type
	subwords []uint64
}

type ByTorC []Entry

func (a ByTorC) Len() int      { return len(a) }
func (a ByTorC) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByTorC) Less(i, j int) bool {
	if a[i].wType != a[j].wType {
		return a[i].wType > a[j].wType
	}
	return a[i].count > a[j].count
}

// index of words_ -> int <- word2int_
type Dictionary struct {
	args_     *Args
	word2int_ []uint64 // word -> int -> words_
	words_    []Entry
	pdiscard_ []Real

	size_      uint64 // word + label 的数量
	nwords_    uint64 // word数量
	nlabels_   uint64 // label数量
	ntokens_   uint64 // 文本实际word数
	lineCount_ uint64 // 文本行数
}

func NewDic(args *Args) *Dictionary {
	word2int := make([]uint64, MAX_VOCAB_SIZE)
	for i := uint64(0); i < MAX_VOCAB_SIZE; i++ {
		word2int[i] = math.MaxUint64 // 1<<32 - 1
	}
	d := new(Dictionary)
	d.word2int_, d.args_ = word2int, args
	return d
}

func (dic *Dictionary) Find(w string) uint64 {
	h := common.Hash(w) % MAX_VOCAB_SIZE
	// 将MAX_VOCAB_SIZE填充满(参照Threshold方法，只会到0.75 * MAX_VOCAB_SIZE)
	for dic.word2int_[h] != math.MaxUint64 && dic.words_[dic.word2int_[h]].word != w {
		h = (h + 1) % MAX_VOCAB_SIZE
	}
	return h
}

func (dic *Dictionary) Add(w string) {
	h := dic.Find(w)
	dic.ntokens_++
	if dic.word2int_[h] == math.MaxUint64 {
		var e Entry
		e.word = w
		e.count = 1
		if strings.Index(w, dic.args_.label) == 0 {
			e.wType = labelType
		} else {
			e.wType = wordType
		}
		dic.words_ = append(dic.words_, e)
		dic.word2int_[h] = dic.size_
		dic.size_ += 1
	} else {
		dic.words_[dic.word2int_[h]].count += 1
	}
}

// word + label 的数量
func (dic *Dictionary) Size() uint64 {
	return dic.size_
}

// word 的数量
func (dic *Dictionary) Words() uint64 {
	return dic.nwords_
}

// label 的数量
func (dic *Dictionary) Labels() uint64 {
	return dic.nlabels_
}
func (dic *Dictionary) Tokens() uint64 {
	return dic.ntokens_
}

func (dic *Dictionary) getId(word string) uint64 {
	h := dic.Find(word)
	return dic.word2int_[h]
}

func (dic *Dictionary) GetWord(id uint64) string {
	if id >= uint64(dic.size_) {
		fmt.Errorf("Error! id: %v >= dic words size: %v", id, dic.size_)
		return ""
	}
	return dic.words_[id].word
}

func (dic *Dictionary) GetNgramsByI(i uint64) ([]uint64, error) {
	if i >= uint64(dic.nwords_) {
		return nil, fmt.Errorf("GetNgramsByI Error! index: %v >= dic words: %v", i, dic.nwords_)
	}
	return dic.words_[i].subwords, nil
}

func (dic *Dictionary) GetNgramsByW(word string) ([]uint64, error) {
	i := dic.getId(word)
	if i > 0 {
		return dic.GetNgramsByI(uint64(i))
	}
	ngrams, _ := dic.ComputeNgrams(word)
	return ngrams, nil
}

func (dic *Dictionary) GetNgrams(word string, ngrams *[]uint64, substrings *[]string) {
	i := dic.getId(word)
	if i >= 0 {
		*ngrams = append(*ngrams, uint64(i))
		*substrings = append(*substrings, dic.words_[i].word)
	} else {
		*ngrams = append(*ngrams, math.MaxUint64)
		*substrings = append(*substrings, word)
	}
	dNgrams, dSubstrings := dic.ComputeNgrams(BOW + word + EOW)
	*ngrams = append(*ngrams, dNgrams...)
	*substrings = append(*substrings, dSubstrings...)
}

/**
初始化ngram表，即每个词都对应一个ngram的表的id列表。比如词 "我想你" ，通过computeNgrams函数可以计算出相应ngram的词索引，假设ngram的词最短为1，最长为3，则就是"<我"，"我想"，"想你"，"你>"，<我想"，"我想你"，"想你>"的子词组成，这里有"<>"因为这里会自动添加这样的词的开始和结束位
*/
func (dic *Dictionary) ComputeNgrams(word string) (ngrams []uint64, substrings []string) {
	// 如果是非字母，数字就继续
	var runes []string
	for len(word) > 0 {
		r, size := utf8.DecodeRuneInString(word)
		if size != 3 {
			return
		}
		word = word[size:]
		runes = append(runes, string(r))
	}

	// 前后追加”<“ ”>“
	runes = append(runes, BOW)
	runes = append([]string{EOW}, runes...)
	for i := 0; i < len(runes); i++ {
		substr := runes[i:]
		for n := 1; n <= len(substr) && uint64(n) <= dic.args_.maxn; n++ {
			if uint64(n) < dic.args_.minn {
				continue
			}
			s := ""
			for _, r := range substr[:n] {
				s += r
			}
			substrings = append(substrings, s)
			ngrams = append(ngrams, uint64(dic.nwords_+common.Hash(s)%uint64(dic.args_.bucket)))
		}
	}
	return ngrams, substrings
}

func (dic *Dictionary) initNgrams() {
	for i := uint64(0); i < dic.size_; i++ {
		dic.words_[i].subwords = append(dic.words_[i].subwords, uint64(i))
		ngrams, _ := dic.ComputeNgrams(dic.words_[i].word)
		dic.words_[i].subwords = append(dic.words_[i].subwords, ngrams...)
	}
}

/**
初始化initTableDiscard表，对每个词根据词的频率获取相应的丢弃概率值，若是给定的阈值小于这个表的值那么就丢弃该词，这里是因为对于频率过高的词可能就是无用词，所以丢弃。比如"的"，"是"等；这里的实现与论文中有点差异，这里是当表中的词小于某个值表示该丢弃，这里因为这里没有对其求1-p形式，而是p+p^2。若是同理转为同方向，则论文是p，现实是p+p^2，这样的做法是使得打压更加宽松点，也就是更多词会被当作无用词丢弃。（不知道原因）
*/
func (dic *Dictionary) initTableDiscard() {
	dic.pdiscard_ = make([]Real, dic.size_)
	for i := uint64(0); i < dic.size_; i++ {
		f := float64(dic.words_[i].count / dic.ntokens_)
		dic.pdiscard_[i] = Real(math.Sqrt(dic.args_.t/f) + dic.args_.t/f)
	}
}

func (dic *Dictionary) discard(id uint64, rand Real) bool {
	if id >= dic.size_ {
		fmt.Printf("discard Error! id: %v, size_: %v\n", id, dic.size_)
		// return false
	}
	// 监督训练直接返回FALSE, 不过滤掉低频词
	if dic.args_.model == Sup {
		return false
	}
	return rand > dic.pdiscard_[id] // 大于随机数的返回false
}

func (dic *Dictionary) GetCounts(wType Entry_type) (counts []uint64) {
	for _, w := range dic.words_ {
		if w.wType == wType {
			counts = append(counts, w.count)
		}
	}
	return counts
}

func (dic *Dictionary) AddNgrams(line *[]uint64, n uint) {
	line_size := uint(len(*line))
	for i := uint(0); i < line_size; i++ {
		h := (*line)[i]
		for j := i + 1; j < line_size && j < i+n; j++ {
			h = h*116049371 + (*line)[j]
			*line = append(*line, (h % uint64(dic.args_.bucket)))
		}
	}
}

func (dic *Dictionary) GetLine(content string, words, labels *[]uint64, rng *rand.Rand) uint64 {
	var ntokens uint64 = 0
	if dic.args_.label != "" && strings.Index(content, dic.args_.label) == 0 {
		strs := strings.SplitN(strings.TrimSpace(content), " ", 2)
		wid := dic.getId(strs[0])
		ntokens++
		*labels = append(*labels, wid)
		for _, w := range common.GetSeg().Cut(strings.TrimSpace(strs[1])) {
			wid := dic.getId(w)
			// 不在dict.words里面，直接过滤
			if wid == math.MaxUint64 {
				continue
			}
			if dic.discard(wid, Real(rng.Float64())) {
				continue
			}
			ntokens++
			*words = append(*words, wid)
			if len(*words) > MAX_LINE_SIZE && dic.args_.model != Sup {
				break
			}
		}
		// } else if strings.Index(content, dic.args_.label) == -1 {
	} else {
		for _, w := range common.GetSeg().Cut(strings.TrimSpace(content)) {
			wid := dic.getId(w)
			// 不在dict.words里面，直接过滤
			if wid == math.MaxUint64 {
				continue
			}
			if dic.discard(wid, Real(rng.Float64())) {
				continue
			}
			ntokens++
			*words = append(*words, wid)
			if len(*words) > MAX_LINE_SIZE && dic.args_.model != Sup {
				break
			}
		}
	}
	return ntokens
}

func (dic *Dictionary) ReadFromFile(file *os.File) error {
	// 加载分词器
	var minThreshold uint64 = 1
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if dic.ntokens_%100000 == 0 {
			runtime.GC()
		}

		dic.lineCount_ += 1
		content := strings.TrimSpace(scanner.Text())

		// 判断有没有label
		if strings.Index(content, dic.args_.label) == 0 {
			strs := strings.SplitN(strings.TrimSpace(content), " ", 2)
			if len(strs) != 2 {
				fmt.Println("invalid line: " + content)
				continue
			}
			dic.Add(strs[0])
			for _, w := range common.GetSeg().Cut(strings.TrimSpace(strs[1])) {
				dic.Add(w)
				if dic.size_ > MAX_VOCAB_SIZE*3/4 {
					minThreshold++
					dic.Threshold(minThreshold, minThreshold)
				}
			}
		} else if strings.Index(content, dic.args_.label) == -1 {
			for _, w := range common.GetSeg().Cut(strings.TrimSpace(content)) {
				dic.Add(w)
				if dic.size_ > MAX_VOCAB_SIZE*3/4 {
					minThreshold++
					dic.Threshold(minThreshold, minThreshold)
				}
			}
		} else {
			fmt.Println("invalid line: " + content)
			continue
		}
	}

	if uint64(dic.args_.minCount) >= minThreshold || uint64(dic.args_.minCountLabel) >= minThreshold {
		dic.Threshold(uint64(dic.args_.minCount), uint64(dic.args_.minCountLabel))
	}
	dic.initTableDiscard()
	dic.initNgrams()
	if dic.size_ == 0 {
		return fmt.Errorf("Empty vocabulary. Try a smaller -minCount value.")
	}
	fmt.Printf("Read %d M words\n", dic.ntokens_/1000000)
	fmt.Printf("Number of words: %d\n", dic.nwords_)
	fmt.Printf("Number of labels: %d\n", dic.nlabels_)
	return nil
}

func (dic *Dictionary) rebuildWords(t, tl uint64) {
	// 排序，规则：label放前面，word放后面，word中：count大的在前面
	sort.Sort(ByTorC(dic.words_))
	// 找到word第一个的index
	pos := sort.Search(len(dic.words_), func(i int) bool {
		return dic.words_[i].wType == wordType
	})
	// 找到label中count小于tl的index
	lPos := sort.Search(len(dic.words_[:pos]), func(i int) bool {
		return dic.words_[i].count < tl
	})
	newLabels := dic.words_[:lPos]

	// 找到word中count小于t的index
	wPos := sort.Search(len(dic.words_[pos:]), func(i int) bool {
		return dic.words_[i].count < t
	})
	newWords := dic.words_[pos:wPos]

	dic.words_ = append(newLabels, newWords...)
}

func (dic *Dictionary) Threshold(t, tl uint64) {
	dic.rebuildWords(t, tl)
	// 初始化words, 语料，words和labels数量
	dic.size_, dic.nwords_, dic.nlabels_ = 0, 0, 0
	for i := uint64(0); i < MAX_VOCAB_SIZE; i++ {
		dic.word2int_[i] = math.MaxUint64
	}
	for _, e := range dic.words_ {
		h := dic.Find(e.word)
		dic.word2int_[h] = dic.size_
		dic.size_ += 1
		if e.wType == wordType {
			dic.nwords_ += 1
		} else {
			dic.nlabels_ += 1
		}
	}
}

func (dic *Dictionary) Save(file *os.File) {
	var e error
	if _, e = file.WriteString(PreDict); e != nil {
		fmt.Println(e.Error())
		return
	}
	if _, e := file.WriteString(strconv.FormatUint(uint64(dic.size_), 10) + Space); e != nil {
		fmt.Println(e.Error())
		return
	}
	if _, e := file.WriteString(strconv.FormatUint(uint64(dic.nwords_), 10) + Space); e != nil {
		fmt.Println(e.Error())
		return
	}
	if _, e := file.WriteString(strconv.FormatUint(uint64(dic.nlabels_), 10) + Space); e != nil {
		fmt.Println(e.Error())
		return
	}
	if _, e := file.WriteString(strconv.FormatUint(dic.ntokens_, 10) + Space); e != nil {
		fmt.Println(e.Error())
		return
	}
	for _, w := range dic.words_ {
		if _, e = file.WriteString(w.word + Dspace); e != nil {
			fmt.Println(e.Error())
			return
		}
		if _, e = file.WriteString(strconv.FormatUint(uint64(w.count), 10) + Dspace); e != nil {
			fmt.Println(e.Error())
			return
		}
		// word type
		var wType string = "0"
		if w.wType == labelType {
			wType = "1"
		}
		if _, e = file.WriteString(wType + Space); e != nil {
			fmt.Println(e.Error())
			return
		}
	}
	file.WriteString("\n")
}

func (dic *Dictionary) Load(content string) {
	dicts := strings.Split(strings.TrimSpace(content), Space)
	dic.words_ = dic.words_[0:0]
	for i := uint64(0); i < MAX_VOCAB_SIZE; i++ {
		dic.word2int_[i] = math.MaxUint64
	}
	var e error
	size_, e := strconv.ParseUint(dicts[0], 10, 32)
	if e != nil {
		fmt.Println(e.Error())
		return
	}
	dic.size_ = uint64(size_)
	nwords_, e := strconv.ParseUint(dicts[1], 10, 32)
	if e != nil {
		fmt.Println(e.Error())
		return
	}
	dic.nwords_ = uint64(nwords_)
	nlabels_, e := strconv.ParseUint(dicts[2], 10, 32)
	if e != nil {
		fmt.Println(e.Error())
		return
	}
	dic.nlabels_ = uint64(nlabels_)
	ntokens_, e := strconv.ParseUint(dicts[3], 10, 64)
	if e != nil {
		fmt.Println(e.Error())
		return
	}
	dic.ntokens_ = ntokens_
	dicts = dicts[4:]
	for i, dict := range dicts {
		wInfo := strings.Split(strings.TrimSpace(dict), Dspace)
		if uint64(i) >= dic.size_ {
			if len(wInfo) == 3 {
				fmt.Printf("Error! end of dict length == 3: %v\n", wInfo)
			}
			continue
		}
		if len(wInfo) != 3 {
			fmt.Printf("Error! dict length != 3: %v\n", wInfo)
			continue
		}
		var e Entry
		e.word = wInfo[0]
		count, err := strconv.ParseUint(wInfo[1], 10, 64)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		e.wType = wordType
		if wInfo[2] == "1" {
			e.wType = labelType
		}
		e.count = count
		dic.words_ = append(dic.words_, e)
		dic.word2int_[dic.Find(e.word)] = uint64(i)
	}
	dic.initTableDiscard()
	dic.initNgrams()
}
