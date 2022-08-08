package bloomfilter

type BitSet struct {
	length uint
	set    []uint64
}

// 字长
const wordSize = uint(64)

// log2WordSize is log2(wordSize)
const log2WordSize = uint(6)

// wordsNeeded 计算需要开辟的uint64数组长度
func wordsNeeded(i uint) int {
	if i > (Cap() - wordSize + 1) {
		return int(Cap() >> log2WordSize)
	}
	return int((i + (wordSize - 1)) >> log2WordSize)
}

//Cap 获取uint在机器中的最大值
//通过位运算符^ 对unit的最小值 uint(0)取反，即可得到最大值
func Cap() uint {
	return ^uint(0)
}
func New(length uint) (bset *BitSet) {
	defer func() {
		if r := recover(); r != nil {
			bset = &BitSet{
				0,
				make([]uint64, 0),
			}
		}
	}()

	bset = &BitSet{
		length,
		make([]uint64, wordsNeeded(length)),
	}

	return bset
}
