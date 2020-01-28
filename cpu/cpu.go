package main

// Example taken from https://github.com/davecheney/high-performance-go-workshop/blob/master/examples/words/

import (
	"bufio"
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

	version := os.Getenv("VERSION")
	if version == "" {
		version = "0"
	}

	words := 0

	switch version {
	case "0":
		words = v0(f)
	case "1":
		words = v1(f)
	default:
		words = v0(f)
	}

	fmt.Printf("%q: %d words\n", os.Args[1], words)
}

func v0(f *os.File) int {
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

func v1(f *os.File) int {
	words := 0
	inword := false

	b := bufio.NewReader(f)

	for {
		r, err := readbyte(b)
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
