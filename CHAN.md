Вопрос:

```go
func Squares(c, quit chan int) {
  // ???
}

func main() {
  mychannel := make(chan int)
  quitchannel := make(chan int)
  sum := 0

  go func() {
    for i := 1; i <= 5; i++ {
      // ???
    }
    fmt.Println(sum)
    // ???
  }()

  Squares(mychannel, quitchannel)
}
```

Ответ:

```go
func Squares(c, quit chan int) {
	for {
		select {
		case <-quit:
			return
		case i := <-c:
			c <- i * i
		}
	}
}

func main() {
	mychannel := make(chan int)
	quitchannel := make(chan int)
	sum := 0

	go func() {
		for i := 1; i <= 5; i++ {
			mychannel <- i
			sum += <-mychannel
		}
		fmt.Println(sum)
		close(quitchannel)
	}()

	Squares(mychannel, quitchannel)
}
```