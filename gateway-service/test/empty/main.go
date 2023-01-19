package main

import (
	"fmt"
	"time"
)

func main() {
	for {
		time.Sleep(5 * time.Second)
		fmt.Println("hello world")
	}

}
