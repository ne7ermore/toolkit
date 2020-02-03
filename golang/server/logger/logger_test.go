package logger

import (
	"fmt"
	"testing"
)

func Test_lt(t *testing.T) {
	lt := Getlogger()
	lt.Info("ooooo")
	lt.Warn(fmt.Sprintf("bbbbb %d", 123123))
	lt.Err("ccccc")
}
