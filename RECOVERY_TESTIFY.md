```go
// Пример тестирования recovery
func TestSomethingThatMightPanic(t *testing.T) {
    assert.NotPanics(t, func() {
        // код, который может вызвать панику
    })

    // или для проверки конкретной паники
    assert.Panics(t, func() {
        panic("ожидаемая паника")
    })
}

// Более сложный пример с проверкой сообщения паники
func TestPanicMessage(t *testing.T) {
    assert.PanicsWithValue(t, "ожидаемое сообщение", func() {
        panic("ожидаемое сообщение")
    })
}
```