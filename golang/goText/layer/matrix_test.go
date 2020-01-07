package layer

// import (
// 	"fmt"
// 	"testing"
// )

// func TestMatrix(t *testing.T) {
// 	fmt.Println("======================TestMatrix=======================")
// 	mat := NewMat()
// 	mat.Init(3, 5)
// 	mat.Uniform(2.3)

// 	v := NewVec()
// 	v.Init(5)
// 	for i := uint64(0); i < 5; i++ {
// 		v.data_[i] = Real(i)
// 	}
// 	v.Mul(3.4)

// 	e := mat.AddRow(v, 2, 0.1)
// 	if e != nil {
// 		t.Fail()
// 	}
// 	fmt.Println(mat)

// 	r, e := mat.DotRow(v, 1)
// 	if e != nil {
// 		fmt.Println(e)
// 		t.Fail()
// 	}
// 	fmt.Println(r)
// }
