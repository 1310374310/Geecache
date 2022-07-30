package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// Hash maps bytes to uint32
type Hash func(data []byte) uint32

// Map constains all hashed keys
type Map struct {
	hash     Hash
	replicas int   // 虚拟节点数目
	keys     []int // 排序过后的哈希值
	hashMap  map[int]string
}

// New creates a Map instance
func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}

	// 若没有自定义，则使用默认的crc32.ChecksumIEEE
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// Add 添加真实节点/机器
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			// 对于每个节点/机器 key 创建对应的m.replicas个虚拟节点
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)

			// 添加节点哈希值与真实节点的映射关系
			m.hashMap[hash] = key
		}
	}

	// 对环上的所有哈希值进行排序
	sort.Ints(m.keys)
}

// Get 获取与对应key最近的节点/机器
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}

	// 非节点的key不存储
	hash := int(m.hash([]byte(key)))

	//  找到第一个匹配的虚拟节点的下标
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})

	// 返回对应的节点/机器
	return m.hashMap[m.keys[idx%len(m.keys)]]
}
