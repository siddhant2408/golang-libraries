package inmemorycache

type simple struct {
	values map[interface{}]interface{}
}

func newSimple() *simple {
	return &simple{
		values: make(map[interface{}]interface{}),
	}
}

func (c *simple) Set(key interface{}, value interface{}) {
	c.values[key] = value
}

func (c *simple) Get(key interface{}) (value interface{}, found bool) {
	value, found = c.values[key]
	return value, found
}

func (c *simple) Remove(key interface{}) {
	delete(c.values, key)
}
