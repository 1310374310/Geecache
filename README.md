# Go语言尝试实现分布式缓存  参照groupcache
1. LRU算法实现
2. 单机并发缓存  基于`sync.Mutex`封装LRU的方法，使其支持并发读写。
3. HTTP服务端
4. 一致性哈希
5. 分布式节点
6. 防止缓存击穿
7. 使用Protobuf通信