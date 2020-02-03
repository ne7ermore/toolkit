package code

import (
	"testing"
)

func Test_Code(t *testing.T) {
	println(NewCode(CodeResourceDuplicated, "").Message)
}

func Test_GetMessage(t *testing.T) {
	println(GetCodeMessage(CodeSystemError))
}
