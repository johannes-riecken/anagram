package main

import (
	"testing"
	"fmt"
	"strings"
)

func TestIntegration(t *testing.T) {
	// "lines" workers for all files hand over lines to "anagrams" workers
	args := []string{"sample0.txt", "sample1.txt", "sample2.txt", "sample3.txt"}
	chans := make([]<-chan OrderedMap, len(args))
	for i, a := range args {
		chans[i] = Anagrams(Lines(a))
	}
	// A "merge" worker merges results from workers in FIFO order
	o := Merge(chans)
	// The results are printed out
	for it := o.Front(); it != nil; it = it.Next() {
		v := it.Value()
		if len(v) > 1 {
			fmt.Println(strings.Join(v, " "))
		}
	}
}
