package layer

// import (
// 	"fmt"
// 	"testing"
// )

// func Test_NewModel(t *testing.T) {
// 	fmt.Println("======================Test_NewModel=======================")
// 	a := InitArgs(Cbow, Ns, "../testdata/input", "../testdata/output", "__label__", "", 0.05, 1e-4, 100, 6, 5, 1, 2, 1, 2, 12, 2, 0, 20000, 5, 1, 5)
// 	wi, wo := NewMat(), NewMat()
// 	wi.Init(6, 6)
// 	wo.Init(6, 6)
// 	m := NewModel(wi, wo, a, 1)
// 	fmt.Println(m)

// 	fmt.Println("======================Test_BuildTree=======================")
// 	m.buildTree([]uint64{15, 8, 6, 5, 3, 1})
// 	fmt.Println(m.tree)

// 	fmt.Println("======================initTableNegatives=======================")
// 	m.SetTargetCounts([]uint64{15, 8, 6, 5, 3, 1})
// 	fmt.Println(m.negatives)

// 	fmt.Println("======================Uniform input=======================")
// 	m.wi_.Uniform(Real(1.0 / 6))
// 	fmt.Println(m.wi_)

// 	fmt.Println("======================Uniform output=======================")
// 	m.wo_.Uniform(Real(1.0 / 6))
// 	fmt.Println(m.wo_)

// 	fmt.Println("======================ComputeHidden=======================")
// 	m.ComputeHidden([]uint64{5, 4, 3, 2, 1, 0}, m.hidden_)
// 	fmt.Println(m.hidden_)

// 	fmt.Println("======================BinaryLogistic=======================")
// 	r, e := m.BinaryLogistic(1, true, 0.1)
// 	if e != nil {
// 		fmt.Println(e)
// 		t.Fail()
// 	}
// 	fmt.Println(r)

// 	fmt.Println("======================NegativeSampling=======================")
// 	loss, e := m.NegativeSampling(1, 0.1)
// 	if e != nil {
// 		fmt.Println(e)
// 		t.Fail()
// 	}
// 	fmt.Println(r)

// 	fmt.Println("======================Softmax=======================")
// 	loss, e = m.Softmax(1, 0.1)
// 	if e != nil {
// 		fmt.Println(e)
// 		t.Fail()
// 	}
// 	fmt.Println(loss)

// 	fmt.Println("======================HierarchicalSoftmax=======================")
// 	loss, e = m.HierarchicalSoftmax(1, 0.1)
// 	if e != nil {
// 		fmt.Println(e)
// 		t.Fail()
// 	}
// 	fmt.Println(loss)

// 	fmt.Println("======================getLoss=======================")
// 	m.loss_ = 1.11
// 	m.nexamples_ = 32
// 	fmt.Println(m.getLoss())

// 	fmt.Println("======================Predict=======================")

// 	a.loss = Hs
// 	m.Predict(3, []uint64{5, 4, 3, 2, 1, 0}, m.hidden_, m.output_)
// 	for _, item := range m.pq {
// 		fmt.Println(item.r)
// 	}
// 	fmt.Println("======================let me cut=======================")
// 	a.loss = Softmax
// 	m.Predict(3, []uint64{5, 4, 3, 2, 1, 0}, m.hidden_, m.output_)
// 	for _, item := range m.pq {
// 		fmt.Println(item.r)
// 	}
// }

// func Test_init(t *testing.T) {
// 	fmt.Println("======================initLog=======================")
// 	fmt.Println(initLog())

// 	fmt.Println("======================initSigmoid=======================")
// 	fmt.Println(initSigmoid())
// }
