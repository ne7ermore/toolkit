// Copyright 2017 All rights reserved.
// Author: ne7ermore.
//
package goHash

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type MapWord struct {
	MW     GoHash
	Length int64
}

func NewMapWord() *MapWord {
	return &MapWord{newHash(), 0}
}

func (m *MapWord) LoadMapWords(file string) error {
	fmt.Println("Loading map words, file: " + file)
	wFile, err := os.Open(file)
	if err != nil {
		return err
	}
	defer wFile.Close()
	scanner := bufio.NewScanner(wFile)
	var length int64 = 0
	for scanner.Scan() {
		length++
		line := strings.TrimSpace(scanner.Text())
		fields := strings.Split(line, "\t")
		if len(fields) != 2 {
			return fmt.Errorf("fmt error line: %v", length)

		}
		if !m.MW.Has(fields[0]) {
			m.MW.Tuh[fields[0]] = fields[1]
		}
	}
	if e := scanner.Err(); e != nil {
		return fmt.Errorf("error occurred while reading: %v", file)
	}
	m.Length = length
	fmt.Sprintf("load %v words", m.Length)
	return nil
}
