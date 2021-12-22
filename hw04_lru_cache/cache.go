package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type cacheItem struct {
	key   Key
	value interface{}
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	if item, ok := c.items[key]; ok {
		item.Value = cacheItem{key, value}
		c.queue.MoveToFront(item)
		return true
	}

	c.items[key] = c.queue.PushFront(cacheItem{key, value})
	c.CheckCapacity()

	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	if item, ok := c.items[key]; ok {
		itemVal := item.Value.(cacheItem)
		c.queue.MoveToFront(item)
		return itemVal.value, true
	}

	return nil, false
}

func (c *lruCache) Clear() {
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}

func (c *lruCache) CheckCapacity() {
	if c.queue.Len() <= c.capacity {
		return
	}

	lastItem := c.queue.Back()
	lastItemVal := lastItem.Value.(cacheItem)

	c.queue.Remove(lastItem)
	delete(c.items, lastItemVal.key)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
