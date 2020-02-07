// Example taken from https://rakyll.org/mutexprofile/
// Modified to work as standalone application. Included print and additional sleep

package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/pkg/profile"
)

func main() {
	defer profile.Start(profile.MutexProfile, profile.ProfilePath(".")).Stop()

	var mu sync.Mutex
	var items = make(map[int]struct{})
	var wg = sync.WaitGroup{}

	fmt.Println("Start mutex contention")

	// 6.2 million go routines
	for i := 0; i < 2500*2500; i++ {
		wg.Add(1)
		go func(i int) {
			mu.Lock()
			defer mu.Unlock()
			defer wg.Done()

			items[i] = struct{}{}
		}(i)
	}

	fmt.Println("Wait")
	wg.Wait()
	fmt.Println("Sleep 5s")
	time.Sleep(5 * time.Second)
}
