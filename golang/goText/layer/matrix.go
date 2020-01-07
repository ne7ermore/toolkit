package layer

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

type Matrix struct {
	data_ []Real
	m_    uint64 // 行数 (m_个vector)
	n_    uint64 // 列数 (等于一个vector的长度)
}

func NewMat() *Matrix {
	return new(Matrix)
}

// m个vec,每个vec中n个成员
func (mat *Matrix) Init(m, n uint64) {
	mat.m_, mat.n_ = m, n
	mat.data_ = make([]Real, m*n)
}

func (mat *Matrix) Uniform(r Real) {
	for i := uint64(0); i < mat.m_*mat.n_; i++ {
		mat.data_[i] = Real(rand.Float64()-0.5) * r * 2.0
	}
}

func (mat *Matrix) AddRow(vec *Vector, i uint64, r Real) error {
	if i >= mat.m_ {
		return fmt.Errorf("AddRow Error! index: %v >= matrix m_: %v\n", i, mat.m_)
	}
	if vec.m_ != mat.n_ {
		return fmt.Errorf("AddRow Error! vec m_: %v != mat n_: %v\n", vec.m_, mat.n_)
	}
	for j := uint64(0); j < mat.n_; j++ {
		mat.data_[i*mat.n_+j] += r * vec.data_[j] // matrix和vector的区别就是增加一个倍数a
	}
	return nil
}

func (mat *Matrix) DotRow(vec *Vector, i uint64) (Real, error) {
	if i >= mat.m_ {
		return Real(0.0), fmt.Errorf("DotRow Error! index: %v >= matrix m_: %v\n", i, mat.m_)
	}
	if vec.m_ != mat.n_ {
		return Real(0.0), fmt.Errorf("DotRow Error! vec m_: %v != mat n_: %v\n", vec.m_, mat.n_)
	}
	var d Real = 0.0
	for j := uint64(0); j < mat.n_; j++ {
		d += mat.data_[i*mat.n_+j] * vec.data_[j] // 一行的dot = 每一列matrix * 每一列vector
	}
	return d, nil
}

func (mat *Matrix) Save(file *os.File, typeFile string) {
	var e error
	if _, e = file.WriteString(typeFile); e != nil {
		fmt.Println(e.Error())
		return
	}
	if _, e = file.WriteString(strconv.FormatUint(mat.m_, 10) + Space); e != nil {
		fmt.Println(e.Error())
		return
	}
	if _, e = file.WriteString(strconv.FormatUint(mat.n_, 10) + Space); e != nil {
		fmt.Println(e.Error())
		return
	}
	for _, d := range mat.data_ {
		if _, e = file.WriteString(strconv.FormatFloat(float64(d), 'E', -1, 64) + Dspace); e != nil {
			fmt.Println(e.Error())
			return
		}
	}
	file.WriteString("\n")
}

func (mat *Matrix) Load(content string, typeFile string) {
	mats := strings.Split(strings.TrimSpace(content), Space)
	if len(mats) != 3 {
		fmt.Printf("Load Error! Invalid type: %v, length of mats: %v\n", typeFile, typeFile, len(mats))
		return
	}
	// uint64
	m_, e := strconv.ParseUint(mats[0], 10, 64)
	if e != nil {
		fmt.Println(e.Error())
		return
	}
	mat.m_ = m_
	n_, e := strconv.ParseUint(mats[1], 10, 64)
	if e != nil {
		fmt.Println(e.Error())
		return
	}
	mat.n_ = n_
	mat.data_ = mat.data_[0:0]
	datas := strings.Split(strings.TrimSpace(mats[2]), Dspace)
	datas = datas[:len(datas)-1]
	for _, data := range datas {
		d, e := strconv.ParseFloat(data, 64)
		if e != nil {
			fmt.Printf("Error!, data: %v, err: %v\n", data, e.Error())
			return
		}
		mat.data_ = append(mat.data_, Real(d))
	}
	if uint64(len(mat.data_)) != mat.m_*mat.n_ {
		fmt.Printf("Load Error! length of data_: %v, m_*n_: %v\n", len(mat.data_), mat.m_*mat.n_)
		return
	}
}
