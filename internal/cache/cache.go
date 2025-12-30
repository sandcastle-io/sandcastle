package cache

type Cache struct {
	pods map[string]string
}

func NewCache() *Cache {
	return &Cache{
		pods: make(map[string]string),
	}
}
