package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32 //哈希环，2^32个节点

type Map struct {
	hash          Hash
	virtualFactor int            //虚拟节点倍数
	keys          []int          // 有序数组表示哈希环
	hashMap       map[int]string //虚拟节点映射
}

func New(virtualFactor int, fn Hash) *Map {
	m := &Map{
		virtualFactor: virtualFactor,
		hash:          fn,
		hashMap:       make(map[int]string),
	}
	if m.hash == nil { //业务可以自己传递哈希方法
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

func (m *Map) IsEmpty() bool {
	return len(m.keys) == 0
}

func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.virtualFactor; i++ { //添加虚拟节点
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
	sort.Ints(m.keys)
}

func (m *Map) Get(key string) string {
	if m.IsEmpty() {
		return ""
	}

	hash := int(m.hash([]byte(key)))

	idx := sort.Search(len(m.keys), func(i int) bool { return m.keys[i] >= hash })

	if idx == len(m.keys) {
		idx = 0
	}

	return m.hashMap[m.keys[idx]]
}
