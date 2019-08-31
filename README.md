# snoc
A Go package to allow you to send message on Sigfox network with SNOC board.

## Usage
```go
package main

import (
	"github.com/k-yak/snoc"
)

func main() {
	var s snoc.Sigfox
	s.Init("/dev/ttyAMA0")
	err := s.SendMessage("4242")
	fmt.Println(err)
}
```
