
```go
package main

import "fmt"

type Site struct {
	URL string
}

type Context struct {
	Site // ?? как это правильно называется
	Value string
}

func main() {

	ctx := Context{
		Site: Site{
			URL: "https://example.com",
		},
		Value: "Hello, World!",
	}

	fmt.Println(ctx.URL)
	fmt.Println(ctx.Site.URL)
	fmt.Println(ctx.Value)
}
```