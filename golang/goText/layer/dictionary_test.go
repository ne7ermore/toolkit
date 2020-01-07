package layer

import (
	"fmt"
	"os"
	"testing"
)

func xTestDic(t *testing.T) {
	fmt.Println("======================TestDic=======================")
	a := InitArgs(Cbow, Hs, "../testdata/input-non-sup.txt", "../testdata/output.txt", "__label__", "", 0.05, 1e-4, 100, 100, 5, 1, 0, 1, 2, 12, 2, 0, 20000, 5, 1, 5)
	d := NewDic(a)
	file, err := os.Open("../testdata/input-non-sup.txt")
	if err != nil {
		t.Fail()
	}
	defer file.Close()
	d.ReadFromFile(file)

	fmt.Println(d.words_)

	fmt.Printf("words num: %d\n", d.Words())
	fmt.Printf("labels num: %d\n", d.Labels())
	fmt.Printf("tokens num: %d\n", d.Tokens())
	fmt.Printf("label couts: %d\n", d.GetCounts(labelType))
	fmt.Printf("word couts: %d\n", d.GetCounts(wordType))
}
