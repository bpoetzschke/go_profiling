// Example taken from: https://github.com/davecheney/high-performance-go-workshop/blob/master/examples/esc/sum.go

package main

import "fmt"

func Sum() int {
	var count = 100
	numbers := make([]int, count)
	for i := range numbers {
		numbers[i] = i + 1
	}

	var sum int
	for _, i := range numbers {
		sum += i
	}
	return sum
}

func main() {
	answer := Sum()
	fmt.Println(answer)
}
