package common

import (
	"fmt"
	"testing"
)

func TestSeg(t *testing.T) {
	var s string = "我，《》？“’、。，；：！@#%……&*（）——-、】【‘；、。，·来到 北京free清华大学    你好2222?"

	fmt.Println(GetSeg().Cut(s))
}
