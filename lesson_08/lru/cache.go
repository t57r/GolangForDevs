package lru

type LruCache interface {
	Put(key, value string)
	Get(key string) (string, bool)
}

type Node struct {
	key   string
	value string
	prev  *Node
	next  *Node
}

type LruCacheImpl struct {
	capacity int
	items    map[string]*Node
	head     *Node
	tail     *Node
}

func NewLruCache(capacity int) LruCache {
	head := &Node{}
	tail := &Node{}
	head.next = tail
	tail.prev = head

	return &LruCacheImpl{
		capacity: capacity,
		items:    make(map[string]*Node, capacity),
		head:     head,
		tail:     tail,
	}
}

func (c *LruCacheImpl) Get(key string) (string, bool) {
	node, ok := c.items[key]
	if !ok {
		return "", false
	}
	c.moveToFront(node)
	return node.value, true
}

func (c *LruCacheImpl) Put(key, value string) {
	if c.capacity <= 0 {
		// Capacity 0 means we never store anything
		return
	}

	// Update existing
	node, ok := c.items[key]
	if ok {
		node.value = value
		c.moveToFront(node)
		return
	}

	// key does not exist, insert new as MRU (front)
	newNode := &Node{
		key:   key,
		value: value,
	}
	c.items[key] = newNode
	c.insertAfterHead(newNode)

	if len(c.items) > c.capacity {
		lru := c.removeLRU()
		delete(c.items, lru.key)
	}
}

// --- Doubly-linked list helpers

func (c *LruCacheImpl) moveToFront(node *Node) {
	c.detach(node)
	c.insertAfterHead(node)
}

func (c *LruCacheImpl) insertAfterHead(node *Node) {
	first := c.head.next
	node.prev = c.head
	node.next = first
	c.head.next = node
	first.prev = node
}

func (c *LruCacheImpl) detach(node *Node) {
	node.prev.next = node.next
	node.next.prev = node.prev
	node.prev = nil
	node.next = nil
}

func (c *LruCacheImpl) removeLRU() *Node {
	lru := c.tail.prev
	if lru == c.head {
		return nil
	}
	c.detach(lru)
	return lru
}
