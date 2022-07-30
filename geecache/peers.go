package geecache

// peerPicker的接口  根据传入的key选择对应的peer
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// PeerGetter 定义查找缓存区的方法
type PeerGetter interface {
	Get(group string, key string) ([]byte, error)
}
