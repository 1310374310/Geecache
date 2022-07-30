package geecache

import (
	"fmt"
	"log"
	"singleflight"
	"sync"
)

// 回调Getter，当缓存不存在时从数据源获取数据并添加到缓存中
type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

// 定义更广泛的接口型函数，提升扩展性
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

/* 整体逻辑框架
                            是
接收 key --> 检查是否被缓存 -----> 返回缓存值 ⑴
                |  否                         是
                |-----> 是否应当从远程节点获取 -----> 与远程节点交互 --> 返回缓存值 ⑵
                            |  否
                            |-----> 调用`回调函数`，获取值并添加到缓存 --> 返回缓存值 ⑶
*/

type Group struct {
	name      string
	getter    Getter // 数据源
	mainCache cache
	peers     PeerPicker //根据key获取节点

	// 如果一个key有多个并发请求，则只找一次db
	loader *singleflight.Group
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

// RegisterPeers 为组注册一个PeerPicker 用于选择远端节点
func (g *Group) RegisterPeers(peers PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = peers
}

// NewGroup  创建一个用户组示例
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()

	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
		loader:    &singleflight.Group{},
	}

	groups[name] = g
	return g
}

// GetGroup 根据组名返回组对应的对象
func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()
	g := groups[name]
	return g
}

// Get 实现Group的Get方法 实现Getter接口
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}
	// 判断是否被缓存，如果是，直接返回
	if v, ok := g.mainCache.get(key); ok {
		log.Println("[GeeCache] hit")
		return v, nil
	}

	// 如果没有被缓存，尝试加载
	return g.load(key)
}

func (g *Group) load(key string) (value ByteView, err error) {
	// 1 先考虑通过分布式缓存获取 调用 getFromPeer 从其他节点获取
	viewi, err := g.loader.Do(key, func() (interface{}, error) {
		if g.peers != nil {
			if peer, ok := g.peers.PickPeer(key); ok {
				if value, err = g.getFromPeer(peer, key); err == nil {
					return value, nil
				}
				log.Println("[GeeCache] Failed to get from peer", err)
			}
		}
		// 2 如果远端节点中不存在，则从本地数据获取，并加载到缓存
		return g.getLocally(key)
	})
	if err == nil {
		return viewi.(ByteView), nil
	}

	return
}

// getFromPeer  从peer中获取缓存数据
func (g *Group) getFromPeer(peer PeerGetter, key string) (ByteView, error) {
	bytes, err := peer.Get(g.name, key)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: bytes}, nil
}

// getLocally 从本地获取数据，并加入缓存
func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key) // 从数据源获取数据
	if err != nil {
		return ByteView{}, err

	}
	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}
