package lru

import "container/list"

type Cache struct {
	maxBytes int64 // 缓存的最大容量
	mbytes   int64 // 缓存的当前占用
	ll       *list.List
	cache    map[string]*list.Element

	OnEvicted func(key string, value Value)
}

// 定义双向链表的节点的内如
type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

// 定义缓存的构造函数
func New(maxBytes int64, OnEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: OnEvicted,
	}
}

// 定义缓存的查找功能，即get方法
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// 定义缓存的删除方法，独立于put方法之外写个函数
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key) //删除字典的映射
		c.mbytes -= int64(len(kv.key)) + int64(kv.value.Len())

		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// 定义缓存的Add方法
func (c *Cache) Add(key string, value Value) {
	// 先判断关键字是否存在
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.mbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else { // 如果关键字不存在，丢进去，然后看缓存空间是否已经满了
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.mbytes += int64(len(key)) + int64(value.Len())
	}

	// 清空超出的缓存空间
	for c.maxBytes != 0 && c.maxBytes < c.mbytes {
		c.RemoveOldest()
	}
}

// 查看链表中有多少条数据
func (c *Cache) Len() int {
	return c.ll.Len()
}
