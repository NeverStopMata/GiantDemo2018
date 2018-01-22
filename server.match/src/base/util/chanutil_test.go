package util

import (
	"fmt"
	"testing"
)

func TestChan(t *testing.T) {
	ch := make(chan int, 100)
	go func() {
		for i := 0; i < 5000; i++ {
			ch <- i
		}
	}()
	for i := 0; i < 5000; i++ {
		if !IsChanClosed(ch) {
			fmt.Println(i, "--", <-ch)
		}
	}

}
