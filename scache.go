package scache

import (
	"fmt"
	"log"
	pb "scache/scachepb"
	"scache/singleflight"
	"sync"
)

type (
	//远程获取值接口
	Getter interface {
		Get(key string) ([]byte, error)
	}
	SCache struct {
		name      string
		picker    PeerPicker
		mainCache cache
		getter    Getter
		//解决并发缓存击穿
		loader singleflight.SingleFlight
	}
)

var (
	mu      sync.RWMutex
	sCaches = make(map[string]*SCache) //所有节点的map
)

//GetSCache 根据名字获取节点
func GetSCache(name string) *SCache {
	mu.RLock()
	g := sCaches[name]
	mu.RUnlock()
	return g
}

func NewSCache(name string, cacheBytes int, getter Getter) *SCache {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &SCache{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
		loader:    &singleflight.Group{},
	}
	sCaches[name] = g
	return g
}

func (s *SCache) Name() string {
	return s.name
}

//Get 对外提供方法
func (s *SCache) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}

	if v, ok := s.mainCache.get(key); ok {
		log.Println("[GeeCache] hit")
		return v, nil
	}

	//本地没有，调用load加载
	return s.load(key)
}

//RegisterPicker 注册从其他节点获取方法
func (s *SCache) RegisterPicker(peers PeerPicker) {
	if s.picker != nil {
		panic("RegisterPeerPicker called more than once")
	}
	s.picker = peers
}

func (s *SCache) load(key string) (value ByteView, err error) {

	//每个key只获取一次，本地或者远端，使用sync.Once保证唯一
	view, err, _ := s.loader.Do(key, func() (interface{}, error) {
		if s.picker != nil {
			if peer, ok := s.picker.PickPeer(key); ok {
				if value, err = s.getFromPeer(peer, key); err == nil {
					return value, nil
				}
				log.Println("[GeeCache] Failed to get from peer", err)
			}
		}

		return s.getLocally(key)
	})

	if err == nil {
		return view.(ByteView), nil
	}
	return
}

func (s *SCache) populateCache(key string, value ByteView) {
	s.mainCache.add(key, value)
}

func (s *SCache) getLocally(key string) (ByteView, error) {
	bytes, err := s.getter.Get(key)
	if err != nil {
		return ByteView{}, err

	}
	value := ByteView{b: cloneBytes(bytes)}
	s.populateCache(key, value)
	return value, nil
}

func (s *SCache) getFromPeer(peer ProtoGetter, key string) (ByteView, error) {
	req := &pb.Request{
		Group: s.name,
		Key:   key,
	}
	res := &pb.Response{}
	err := peer.Get(req, res)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: res.Value}, nil
}
