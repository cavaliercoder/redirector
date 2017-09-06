package main

type ViewBag map[string]interface{}

func NewViewBag() ViewBag {
	return make(ViewBag, 0)
}

func (vb ViewBag) Add(key string, v interface{}) {
	vb[key] = v
}

func (vb ViewBag) Get(key string) interface{} {
	return vb[key]
}
