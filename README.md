# Or-Channel

Пакет `or` предоставляет функцию для объединения нескольких done-каналов в один. 
Возвращаемый канал закрывается, как только закрывается любой из переданных каналов.

## Установка

```bash
go get github.com/Oleska1601/WBOrChannel
```


### Пример использования
```go
package main

import (
    "fmt"
    "time"
    or "github.com/Oleska1601/WBOrChannel"
)

func main() {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or.Or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)
	fmt.Printf("done after %v", time.Since(start))
}
```