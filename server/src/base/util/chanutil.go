package util

func IsChanClosed(ch chan int) bool {
	if len(ch) == 0 {
		select {
		case v, ok := <-ch:
			if !ok {
				return true
			}
			ch <- v
		}
	}
	return false
}
