package singleflight

import "sync"

// call 表示正在进行中，或已经结束的请求。
type call struct {
	wg  sync.WaitGroup // 避免重入
	val interface{}
	err error
}

// Group 是singleflight的主数据结构，管理不同key的请求
type Group struct {
	mu sync.Mutex //保护m变量不被并发读写而加上的锁
	m  map[string]*call
}

// Do  等待组的的执行逻辑，根据传入fn方法进行数据读取，同时防止
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}

	// 如果在字典中，则直接读取
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait() // 等待所有的任务完成
		return c.val, c.err
	}

	c := new(call)
	c.wg.Add(1)
	g.m[key] = c // 添加刀g.m 表明key已经有对应的请求在处理
	g.mu.Unlock()

	c.val, c.err = fn() //调用fn， 发起请求
	c.wg.Done()         // 请求结束

	g.mu.Lock()
	delete(g.m, key) // 更新g.m
	g.mu.Unlock()

	return c.val, c.err
}
