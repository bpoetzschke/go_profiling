// Example taken from: https://github.com/davecheney/high-performance-go-workshop/blob/master/examples/block/block.go
// Modified to include an additional sleep to demonstrate the blockig profile

package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pkg/profile"
)

func generate(in <-chan int, out chan<- int) {
	i := <-in
	time.Sleep(100 * time.Millisecond)
	out <- i + i
}

func main() {
	defer profile.Start(profile.BlockProfile, profile.ProfilePath(".")).Stop()
	in := make(chan int, 1)
	in <- 1
	var out chan int
	for i := 0; i < 20; i++ {
		out = make(chan int)
		go generate(in, out)
		in = out
	}
	sleep := os.Getenv("SLEEP")
	fmt.Println(<-out)
	if strings.ToLower(sleep) == "true" {
		fmt.Println("Sleep")
		time.Sleep(2 * time.Second)
	}
}
