```go
type IntConfig struct {
	v int
}

func (c *IntConfig) Get() int {
	return c.v
}

func (c *IntConfig) Set(v int) {
	c.v = v
}

type IntConfigGetter interface {
	Get() int
}

type Foo struct {
	config IntConfigGetter
}

func (f Foo) Bar() {
	println(f.config.Get()) // у config есть только .Get(), но не .Set()
}

func NewFoo(config IntConfigGetter) *Foo {
	return &Foo{config}
}
```
