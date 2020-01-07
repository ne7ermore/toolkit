package layer

import (
	"fmt"
)

type Real float32

type Vector struct {
	m_    uint64 // 向量长度，对应matrix中的列数Matrix.n_
	data_ []Real
}

func NewVec() *Vector {
	return new(Vector)
}

func (v *Vector) Init(m uint64) {
	v.m_ = m
	v.data_ = make([]Real, m)
}

func (v *Vector) Size() uint64 {
	return v.m_
}

func (v *Vector) Zero() {
	for i := uint64(0); i < v.m_; i++ {
		v.data_[i] = 0.0
	}
}

func (v *Vector) Mul(a Real) {
	for i := uint64(0); i < v.m_; i++ {
		v.data_[i] *= a
	}
}

// 将matrix中第i组中a.n_(等于v.m_)个成员加到vec中v.m_中
func (v *Vector) AddRow(a *Matrix, i uint64) error {
	if v.m_ != a.n_ {
		return fmt.Errorf("AddRow ERROR! vec.m_: %v != Matrix.n_: %v\n", v.m_, a.n_)
	}
	if i >= a.m_ {
		return fmt.Errorf("AddRow ERROR! Index: %v >= Matrix.m_: %v\n", i, a.m_)
	}
	for j := uint64(0); j < a.n_; j++ {
		v.data_[j] += a.data_[i*a.n_+j]
	}
	return nil
}

func (v *Vector) AddRowR(a *Matrix, i uint64, r Real) error {
	if v.m_ != a.n_ {
		return fmt.Errorf("AddRowR ERROR! vec.m_: %v != Matrix.n_: %v\n", v.m_, a.n_)
	}
	if i >= a.m_ {
		return fmt.Errorf("AddRowR ERROR! Index: %v >= Matrix.m_: %v\n", i, a.m_)
	}
	for j := uint64(0); j < a.n_; j++ {
		v.data_[j] += r * a.data_[i*a.n_+j]
	}
	return nil
}

/**
 * 将matrix分成： n_(vec.m_) * m_(v.n_)
 * 然后将vec的每个成员 * matrix第i([0, m_))列每个成员再相加
 * 成为v的第i个成员值
 */
func (v *Vector) MulMV(a *Matrix, vec *Vector) error {
	if v.m_ != a.m_ {
		return fmt.Errorf("MulMV ERROR! v.m_: %v != Matrix.m_: %v\n", v.m_, a.m_)
	}
	if a.n_ != vec.m_ {
		return fmt.Errorf("MulMV ERROR! a.n_: %v != vec.m_: %v\n", a.n_, vec.m_)
	}
	for i := uint64(0); i < v.m_; i++ {
		v.data_[i] = 0.0
		for j := uint64(0); j < a.n_; j++ {
			v.data_[i] += a.data_[i*a.n_+j] * vec.data_[j]
		}
	}
	return nil
}

// 返回Verctor.data_中最大值得index
func (v *Vector) Argmax() uint64 {
	var max Real = v.data_[0]
	var argmax uint64 = 0
	for i := uint64(0); i < v.m_; i++ {
		if v.data_[i] > max {
			max = v.data_[i]
			argmax = uint64(i)
		}
	}
	return argmax
}
