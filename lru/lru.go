package lru

import "container/list"

//缓存最底层，lru存储
//lru缓存结构，并发不安全
type LRUCache struct {

	//缓存最大容量，值为0代表没有限制
	MaxCapacity int
	//用户注册的回调函数，当一个node被从cache缓存中清除掉时，回调函数将被执行。
	OnEvicted func(key Key, value interface{})

	//双链表
	ll    *list.List
	cache map[interface{}]*list.Element
}

//哈希表的key类型
type Key interface{}

//双链表的节点类型，哈希表的value类型
type node struct {
	key   Key
	value interface{}
}

func New(MaxCapacity int, onEvicted func(Key, interface{})) *LRUCache {
	return &LRUCache{
		MaxCapacity: MaxCapacity,
		ll:          list.New(),
		cache:       make(map[interface{}]*list.Element),
		OnEvicted:   onEvicted,
	}
}

func (c *LRUCache) Add(key Key, value interface{}) {
	if c.cache == nil { //懒加载，使用时才创建
		c.cache = make(map[interface{}]*list.Element)
		c.ll = list.New()
	}
	if ee, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ee)
		ee.Value.(*node).value = value
		return
	}
	ele := c.ll.PushFront(&node{key, value})
	c.cache[key] = ele
	if c.MaxCapacity != 0 && c.ll.Len() > c.MaxCapacity {
		c.RemoveOldest()
	}
}

func (c *LRUCache) Get(key Key) (value interface{}, ok bool) {
	if c.cache == nil {
		return
	}
	if ele, hit := c.cache[key]; hit {
		c.ll.MoveToFront(ele)
		return ele.Value.(*node).value, true
	}
	return
}

func (c *LRUCache) Remove(key Key) {
	if c.cache == nil {
		return
	}
	if ele, hit := c.cache[key]; hit {
		c.removeElement(ele)
	}
}

func (c *LRUCache) RemoveOldest() {
	if c.cache == nil {
		return
	}
	ele := c.ll.Back()
	if ele != nil {
		c.removeElement(ele)
	}
}

func (c *LRUCache) removeElement(e *list.Element) {
	c.ll.Remove(e)
	kv := e.Value.(*node)
	delete(c.cache, kv.key)
	if c.OnEvicted != nil {
		c.OnEvicted(kv.key, kv.value) //调用用户的回调函数，就是要删除下游数据库中的数据
	}
}

func (c *LRUCache) Len() int {
	if c.cache == nil {
		return 0
	}
	return c.ll.Len()
}

func (c *LRUCache) Clear() {
	if c.OnEvicted != nil {
		for _, e := range c.cache {
			kv := e.Value.(*node)
			c.OnEvicted(kv.key, kv.value)
		}
	}
	c.ll = nil
	c.cache = nil
}
