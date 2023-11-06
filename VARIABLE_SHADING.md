```go
func fn(tracing bool) {
  var client *http.Client
	if tracing {
		client, err := createClient()
	} else {
		client, err := createDefaultClient()
	}
	// client не инициализирован из-за := для err
}
```
