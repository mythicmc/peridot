package utils

import "sync"

func PickNonEmptyString(a, b string) string {
	if a != "" {
		return a
	}
	return b
}

func WaitGroupDoneAsChannel(wg *sync.WaitGroup) <-chan struct{} {
	ch := make(chan struct{})
	go func() {
		defer close(ch)
		wg.Wait()
	}()
	return ch
}
