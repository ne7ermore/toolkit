package layer

// import (
// 	"fmt"
// 	"testing"
// )

// func TestVector(t *testing.T) {
// 	fmt.Println("======================TestVector=======================")
// 	v := NewVec()
// 	v.Init(5)
// 	fmt.Printf("init vec: %v\n", v)
// 	if v.Size() != 5 {
// 		t.Fail()
// 	}
// 	v.Mul(3.4)

// 	mat := NewMat()
// 	mat.Init(3, 5)
// 	mat.Uniform(2.3)

// 	v.AddRow(mat, 2)
// 	v.AddRowR(mat, 1, 0.9)

// 	v2 := NewVec()
// 	v2.Init(3)

// 	v2.MulMV(mat, v)
// 	fmt.Println(v2)
// 	fmt.Println(v2.Argmax())
// }
