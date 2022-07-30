package geecache

import (
	"lru"
	"sync"
)

type cache struct {
	mu         sync.Mutex
	lru        *lru.Cache // 实例化LRU
	cacheBytes int64
}

// 互斥锁封装缓存的并发写
func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.lru == nil { // 懒初始化 减少程序内存占用
		c.lru = lru.New(c.cacheBytes, nil)
	}
	c.lru.Add(key, value)
}

// 互斥锁封装缓存的并发读写（因为涉及LRU内部的链表调整）
func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.lru == nil {
		return
	}

	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}
	return
}
