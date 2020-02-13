package main

// Example taken from https://github.com/davecheney/high-performance-go-workshop/blob/master/examples/words/
// See original license: https://github.com/davecheney/high-performance-go-workshop#license-and-materials
// Changes were made to support a environment var to select between buffered reader and regular reader.

import (
	"fmt"
	"io"
	"log"
	"os"
	"unicode"

	"github.com/pkg/profile"
)

func readbyte(r io.Reader) (rune, error) {
	var buf [1]byte
	_, err := r.Read(buf[:])
	return rune(buf[0]), err
}

func main() {
	defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()

	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalf("could not open file %q: %v", os.Args[1], err)
	}

	words := countWords(f)

	fmt.Printf("%q: %d words\n", os.Args[1], words)
}

func countWords(f *os.File) int {
	words := 0
	inword := false

	for {
		r, err := readbyte(f)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("could not read file %q: %v", os.Args[1], err)
		}
		if unicode.IsSpace(r) && inword {
			words++
			inword = false
		}
		inword = unicode.IsLetter(r)
	}

	return words
}

// func countWords(f *os.File) int {
// 	words := 0
// 	inword := false

// 	b := bufio.NewReader(f)

// 	for {
// 		r, err := readbyte(b)
// 		if err == io.EOF {
// 			break
// 		}
// 		if err != nil {
// 			log.Fatalf("could not read file %q: %v", os.Args[1], err)
// 		}
// 		if unicode.IsSpace(r) && inword {
// 			words++
// 			inword = false
// 		}
// 		inword = unicode.IsLetter(r)
// 	}

// 	return words
// }
