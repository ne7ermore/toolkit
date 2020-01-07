package common

import (
	"fmt"
	"sync"
	"unicode"

	"github.com/wangbin/jiebago"
)

var seger Segmenter

type Segmenter struct {
	seg   *jiebago.Segmenter
	mutex sync.Mutex
}

func init() {
	fmt.Println("init segmenter")
	seger.seg = new(jiebago.Segmenter)
	seger.seg.LoadDictionary("dict/dict.txt")
}

func GetSeg() *Segmenter {
	return &seger
}

func (seger *Segmenter) Cut(s string) []string {
	seger.mutex.Lock()
	defer seger.mutex.Unlock()

	var ts []string
	words := seger.seg.Cut(s, false)
	for w := range words {
		// 过滤标点、空格和tab
		if unicode.IsPunct([]rune(w)[0]) || unicode.IsSpace([]rune(w)[0]) || w == "＝" || w == "｀" {
			continue
		}
		ts = append(ts, w)
	}
	return ts
}
