package layer

import (
	"container/heap"
	"fmt"
	"math"
	"math/rand"
	"sort"

	"goText/common"
)

const (
	SIGMOID_TABLE_SIZE  = 512
	LOG_TABLE_SIZE      = 512
	MAX_SIGMOID         = 8
	NEGATIVE_TABLE_SIZE = 10000000
)

type Model struct {
	wi_        *Matrix
	wo_        *Matrix
	args_      *Args
	hidden_    *Vector
	output_    *Vector
	grad_      *Vector
	hsz_       uint64
	isz_       uint64
	osz_       uint64
	loss_      Real
	nexamples_ uint64
	t_sigmoid  []Real
	t_log      []Real
	negatives  []uint64
	negpos     int
	paths      [][]uint64
	codes      [][]bool
	tree       []Node
	rng        *rand.Rand

	pq itemPQ
}

type Node struct {
	parent uint64
	left   uint64
	right  uint64
	count  uint64
	binary bool
}

func NewModel(wi, wo *Matrix, args *Args, rng int64) *Model {
	m := new(Model)
	m.wi_ = wi
	m.wo_ = wo
	m.args_ = args
	m.isz_ = wi.m_
	m.osz_ = wo.m_
	m.hsz_ = uint64(args.dim)
	m.negpos = 0
	m.loss_ = 0.0
	m.nexamples_ = 1
	m.t_sigmoid = initSigmoid()
	m.t_log = initLog()
	m.rng = rand.New(rand.NewSource(rng))

	hidden, output, grad := NewVec(), NewVec(), NewVec()
	hidden.Init(uint64(args.dim))
	output.Init(wo.m_)
	grad.Init(uint64(args.dim))
	m.hidden_ = hidden
	m.output_ = output
	m.grad_ = grad
	return m
}

/**
无论是负采样还是层次 softmax，在神经网络的计算图中，所有 LR 都会依赖于 hidden_的值，所以 hidden_ 的梯度 grad_ 是各个 LR 的反向传播的梯度的累加
*/
func (m *Model) BinaryLogistic(target uint64, label bool, lr Real) (Real, error) {
	// 激活函数
	var labelR Real = 0.0
	if label {
		labelR = 1.0
	}
	r, e := m.wo_.DotRow(m.hidden_, target)
	if e != nil {
		return Real(0.0), e
	}

	if math.IsNaN(float64(r)) {
		r = 0.0
	}

	score := m.sigmoid(r)
	// 逻辑回归
	alpha := lr * (labelR - score)
	m.grad_.AddRowR(m.wo_, target, alpha)
	m.wo_.AddRow(m.hidden_, target, alpha)
	if label {
		return -m.log(score), nil
	} else {
		return -m.log(Real(1.0) - score), nil
	}
}

/**
训练时每次选择一个正样本，随机采样几个负样本，每种输出都对应一个参数向量，保存于 wo_ 的各行。对所有样本的参数更新，都是一次独立的 LR 参数更新。
*/
func (m *Model) NegativeSampling(target uint64, lr Real) (Real, error) {
	var loss Real = 0.0
	for n := uint64(0); n <= m.args_.neg; n++ {
		if n == 0 {
			if r, e := m.BinaryLogistic(target, true, lr); e != nil {
				return loss, e
			} else {
				loss += r
			}
		} else {
			if r, e := m.BinaryLogistic(m.GetNegative(target), false, lr); e != nil {
				return loss, e
			} else {
				loss += r
			}
		}
	}
	return loss, nil
}

/**
对于每个目标词，都可以在构建好的霍夫曼树上确定一条从根节点到叶节点的路径，路径上的每个非叶节点都是一个LR，参数保存在 wo_ 的各行上，训练时，这条路径上的 LR 各自独立进行参数更新。
*/
func (m *Model) HierarchicalSoftmax(target uint64, lr Real) (Real, error) {
	var loss Real = 0.0
	binaryCode, pathToRoot := m.codes[target], m.paths[target]
	for i := 0; i < len(pathToRoot); i++ {
		if r, e := m.BinaryLogistic(uint64(pathToRoot[i]), binaryCode[i], lr); e != nil {
			return 0.0, e
		} else {
			loss += r
		}
	}
	return loss, nil
}

// 计算output
func (m *Model) computeOutputSoftmax(hidden, output *Vector) {
	/**
	 * output matrix和隐藏层传输到输出层
	 */
	output.MulMV(m.wo_, hidden)
	// if output.Size() != m.osz_ {
	// 	fmt.Printf("Error: length of Output: %v does not match osz: %v\n", output.Size(), m.osz_)
	// 	return
	// }
	var max, z Real = output.data_[0], 0.0
	for _, v := range output.data_ {
		if v > max {
			max = v
		}
	}
	for i, v := range output.data_ {
		output.data_[i] = Real(math.Exp(float64(v - max)))
		z += output.data_[i]
	}
	for i, v := range output.data_ {
		output.data_[i] = v/z + math.SmallestNonzeroFloat64
	}
}

func (m *Model) Softmax(target uint64, lr Real) (Real, error) {
	m.computeOutputSoftmax(m.hidden_, m.output_)
	var label, alpha Real
	for i := uint64(0); i < m.osz_; i++ {
		if i == target {
			label = 1.0
		} else {
			label = 0.0
		}
		alpha = lr * (label - m.output_.data_[i])
		if e := m.grad_.AddRowR(m.wo_, i, alpha); e != nil {
			return label, e
		}
		if e := m.wo_.AddRow(m.hidden_, i, alpha); e != nil {
			return label, e
		}
	}

	// 无限趋近于0处理
	if math.IsNaN(float64(m.output_.data_[target])) {
		return -m.log(0.0), nil
	}
	return -m.log(m.output_.data_[target]), nil
}

func (m *Model) ComputeHidden(input []uint64, hidden *Vector) {
	if hidden.Size() != m.hsz_ {
		fmt.Printf("Error: length of Hidden: %v does not match hsz: %v\n", hidden.Size(), m.hsz_)
		return
	}
	// 隐藏层(hidden)的Vector.m_(长度)和输入层(wi)的每一行的列长度一致(Matrix.n_)，循环做的就是把输入层(wi)每一行index相同的的值相加并放入hidden.data_中
	for _, d := range input {
		e := hidden.AddRow(m.wi_, d)
		if e != nil {
			fmt.Printf("ComputeHidden Error: %v\n", e.Error())
			return
		}
	}
	// // 每个成员 * （1.0 / line的长度）
	// hidden.Mul(Real(1.0 / float64(len(input))))
}

/**
该函数有三个参数，分别是，“类标签”，“学习率”。

“输入”是一个 uint64 dictionary 里的ID。
分类问题(supervised)：这个数组代表输入的短文本；
word2vec：这个数组代表一个词的上下文。

“类标签”是一个 uint64 变量。
word2vec：它就是带预测的词的ID；
分类问题(supervised)：它就是类的 label 在 dictionary 里的 ID。因为 label和词在词表里一起存放，所以有统一的 ID 体系。
**/
func (m *Model) Update(input []uint64, target uint64, lr Real) {
	if target >= m.osz_ {
		fmt.Printf("Update Error! target: %v more than osz: %v\n", target, m.osz_)
		return
	}
	if len(input) == 0 {
		fmt.Printf("Update Error: length of input is 0\n")
		return
	}
	// 计算前向传播：输入层 -> 隐藏层
	m.ComputeHidden(input, m.hidden_)
	// 根据输出层的不同结构，调用不同的函数，在各个函数中，
	// 不仅通过前向传播算出了 loss_，还进行了反向传播，计算出了 grad_，后面逐一分析。
	// 1. 负采样
	if m.args_.loss == Ns {
		loss, e := m.NegativeSampling(target, lr)
		if e != nil {
			fmt.Println(e)
			return
		}
		m.loss_ += loss
	} else if m.args_.loss == Hs {
		loss, e := m.HierarchicalSoftmax(target, lr)
		if e != nil {
			fmt.Println(e)
			return
		}
		m.loss_ += loss
	} else if m.args_.loss == Softmax {
		loss, e := m.Softmax(target, lr)
		if e != nil {
			fmt.Println(e)
			return
		}
		m.loss_ += loss
	} else {
		fmt.Printf("Error: loss type: %v not support\n", m.args_.loss)
		return
	}
	m.nexamples_ += 1
	// 如果是在训练分类器，就将 grad_ 除以 input_ 的大小
	if m.args_.model == Sup {
		m.grad_.Mul(Real(1 / float64(len(input))))
	}
	// 反向传播，将 hidden_ 上的梯度传播到 wi_ 上的对应行
	for _, i := range input {
		m.wi_.AddRow(m.grad_, i, 1.0)
	}
}

func (m *Model) FindKBest(k uint, hidden, output *Vector) {
	m.computeOutputSoftmax(hidden, output)
	for i := uint64(0); i < m.osz_; i++ {
		if uint(m.pq.Len()) == k && byRealLess(m.pq[0].r, output.data_[i]) {
			continue
		}
		heap.Push(&m.pq, newItem(output.data_[i], i))
		if uint(m.pq.Len()) > k {
			heap.Pop(&m.pq)
		}
	}
}

func (m *Model) Dfs(k uint, node uint64, score Real, hidden *Vector) {
	if uint(m.pq.Len()) == k && byRealLess(m.pq[0].r, score) {
		return
	}
	// 只输出叶子节点的结果
	if m.tree[node].left == math.MaxUint64 && m.tree[node].right == math.MaxUint64 {
		heap.Push(&m.pq, newItem(score, node))
		if uint(m.pq.Len()) > k {
			heap.Pop(&m.pq)
		}
		return
	}
	rl, err := m.wo_.DotRow(hidden, node-m.osz_)
	if err != nil {
		fmt.Println(err)
		return
	}
	f := m.sigmoid(rl)
	// 将 score 累加后递归向下收集结果， 分数累加，右边为true,左边为false
	m.Dfs(k, m.tree[node].left, score+m.log(1.0-f), hidden)
	m.Dfs(k, m.tree[node].right, score+m.log(f), hidden)
}

/**
predict 函数可以用于给输入数据打上 1 ～ K 个类标签，并输出各个类标签对应的概率值，对于层次 softmax，我们需要遍历霍夫曼树，找到 top－K 的结果，对于普通 softmax（包括负采样和 softmax 的输出），我们需要遍历结果数组，找到 top－K。
*/
func (m *Model) Predict(k uint, input []uint64, hidden, output *Vector) {
	m.ComputeHidden(input, hidden)
	// 如果是层次 softmax，使用 dfs 遍历霍夫曼树的所有叶子节点，找到 top－k 的概率
	m.pq = make(itemPQ, 0, k)
	if m.args_.loss == Hs {
		m.Dfs(k, 2*m.osz_-2, 0.0, hidden)
	} else {
		// 如果是普通 softmax，在结果数组里找到 top-k
		m.FindKBest(k, hidden, output)
	}
	sort.Sort(byRank(m.pq)) // 排序
}

func (m *Model) SetTargetCounts(counts []uint64) {
	if uint64(len(counts)) != m.osz_ {
		fmt.Printf("Error! length of counts: %v does not match osz: %v\n", len(counts), m.osz_)
		return
	}
	if m.args_.loss == Ns {
		m.initTableNegatives(counts)
	} else if m.args_.loss == Hs {
		m.buildTree(counts)
	}
}

func (m *Model) initTableNegatives(counts []uint64) {
	var z Real = 0.0
	for _, count := range counts {
		z += Real(math.Pow(float64(count), 0.5)) // counts[i]^0.5
	}
	for _, count := range counts {
		c := Real(math.Pow(float64(count), 0.5))
		for i := uint64(0); i < uint64(c/z*NEGATIVE_TABLE_SIZE); i++ {
			m.negatives = append(m.negatives, i)
		}
	}
	m.negatives = common.Shuffle(m.negatives, m.rng)
}

func (m *Model) GetNegative(target uint64) (negative uint64) {
	for {
		negative = m.negatives[m.negpos]
		m.negpos = (m.negpos + 1) % len(m.negatives)
		if target == negative {
			break
		}
	}
	return negative
}

/**
在学信息论的时候接触过构建 Huffman 树的算法，课本中的方法描述往往是：
找到当前权重最小的两个子树，将它们合并
算法的性能取决于如何实现这个逻辑。网上的很多实现都是在新增节点都时遍历一次当前所有的树，这种算法的复杂度是 O(n2)O(n2)，性能很差。
聪明一点的方法是用一个优先级队列来保存当前所有的树，每次取 top 2，合并，加回队列。这个算法的复杂度是 O(nlogn)O(nlogn)，缺点是必需使用额外的数据结构，而且进堆出堆的操作导致常数项较大。
word2vec 以及 fastText 都采用了一种更好的方法，时间复杂度是 O(nlogn)O(nlogn)，只用了一次排序，一次遍历，简洁优美，但是要理解它需要进行一些推理。
算法首先对输入的叶子节点进行一次排序（O(nlogn)O(nlogn) ），然后确定两个下标 leaf 和 node，leaf 总是指向当前最小的叶子节点，node 总是指向当前最小的非叶子节点，所以，最小的两个节点可以从 leaf, leaf - 1, node, node + 1 四个位置中取得，时间复杂度 O(1)O(1)，每个非叶子节点都进行一次，所以总复杂度为 O(n)O(n)，算法整体复杂度为 O(nlogn)O(nlogn)。
*/
func (m *Model) buildTree(counts []uint64) {
	// counts 数组保存每个叶子节点的词频，降序排列
	// 分配所有节点的空间
	tree := make([]Node, 2*m.osz_-1)
	for i := uint64(0); i < 2*m.osz_-1; i++ {
		tree[i].parent = math.MaxUint64
		tree[i].left = math.MaxUint64
		tree[i].right = math.MaxUint64
		tree[i].count = math.MaxUint64
		tree[i].binary = false
	}
	for i := uint64(0); i < m.osz_; i++ {
		tree[i].count = counts[i]
	}
	// leaf 指向当前未处理的叶子节点的最后一个，也就是权值最小的叶子节点
	// node 指向当前未处理的非叶子节点的第一个，也是权值最小的非叶子节点
	leaf, node := m.osz_-1, m.osz_
	for i := m.osz_; i < 2*m.osz_-1; i++ {
		mini := make([]uint64, 2)
		for j := 0; j < 2; j++ {
			if leaf != math.MaxUint64 && tree[leaf].count < tree[node].count {
				mini[j] = leaf
				leaf -= 1
			} else {
				mini[j] = node
				node += 1
			}
		}
		tree[i].left, tree[i].right = mini[0], mini[1]
		tree[i].count = tree[mini[0]].count + tree[mini[1]].count
		tree[mini[0]].parent, tree[mini[1]].parent, tree[mini[1]].binary = i, i, true
	}

	for i := uint64(0); i < m.osz_; i++ {
		path, code := make([]uint64, 0), make([]bool, 0)
		j := i
		for tree[j].parent != math.MaxUint64 {
			path = append(path, tree[j].parent-m.osz_)
			code = append(code, tree[j].binary)
			j = tree[j].parent
		}
		m.paths = append(m.paths, path)
		m.codes = append(m.codes, code)
	}
	m.tree = tree
}

func (m *Model) getLoss() Real {
	return m.loss_ / Real(m.nexamples_)
}

func (m *Model) sigmoid(x Real) Real {
	if x < Real(-MAX_SIGMOID) {
		return Real(0.0)
	} else if x > Real(MAX_SIGMOID) {
		return Real(1.0)
	} else {
		i := int((x + Real(MAX_SIGMOID)) * Real(SIGMOID_TABLE_SIZE) / Real(MAX_SIGMOID) / 2)
		return m.t_sigmoid[i]
	}
}

func (m *Model) log(x Real) Real {
	if x > Real(1.0) {
		return Real(0.0)
	}
	i := int(x * Real(LOG_TABLE_SIZE))
	return m.t_log[i]
}

func initSigmoid() []Real {
	t_sigmoid := make([]Real, SIGMOID_TABLE_SIZE+1)
	for i := 0; i < SIGMOID_TABLE_SIZE+1; i++ {
		x := float64((i*2*MAX_SIGMOID)/SIGMOID_TABLE_SIZE - MAX_SIGMOID) // -MAX_SIGMOID ~ MAX_SIGMOID
		t_sigmoid[i] = Real(1.0 / (1.0 + math.Exp(x)))
	}
	return t_sigmoid
}

func initLog() []Real {
	t_log := make([]Real, LOG_TABLE_SIZE+1)
	for i := 0; i < LOG_TABLE_SIZE+1; i++ {
		x := float64((float64(i) + 1e-5) / LOG_TABLE_SIZE)
		t_log[i] = Real(math.Log(x))
	}
	return t_log
}

/**
 *  heap相关
 */
type Item struct {
	r     Real
	index uint64
}

type itemPQ []*Item

func (pq itemPQ) Len() int           { return len(pq) }
func (pq itemPQ) Less(i, j int) bool { return pq[i].r < pq[j].r } // 首位最小
func (pq itemPQ) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }
func (pq *itemPQ) Push(x interface{}) {
	item := x.(*Item)
	*pq = append(*pq, item)
}
func (pq *itemPQ) Pop() interface{} {
	n := len(*pq)
	item := (*pq)[n-1]
	*pq = (*pq)[0 : n-1]
	return item
}

func newItem(r Real, index uint64) *Item {
	i := new(Item)
	i.r, i.index = r, index
	return i
}
func byRealLess(itemR, opR Real) bool { return itemR > opR }

// 排序
type byRank []*Item

func (by byRank) Len() int           { return len(by) }
func (by byRank) Swap(i, j int)      { by[i], by[j] = by[j], by[i] }
func (by byRank) Less(i, j int) bool { return by[i].r > by[j].r }
