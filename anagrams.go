package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
)

// An Anagrams instance stores the anagrams for a word list.
type Anagrams struct {
	State map[string][]string
}

func main() {
	if len(os.Args) == 1 {
		log.Fatal("Invocation: anagrams.go [words_file]")
	}

	contents, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal("error reading anagrams file:", err)
	}
	lines := strings.Split(string(contents), "\n")
	lines = lines[:len(lines)-1] // remove empty string from trailing newline
	anagrams := NewAnagrams(lines)
	fmt.Println(anagrams.String())
}

// Merges the state of another Anagrams instance into this instance. This
// is to allow concurrent calculations in a map-reduce style.
func (a0 Anagrams) Merge(a1 Anagrams) Anagrams {
	for x := range a1.State {
		a0.State[x] = append(a0.State[x], a1.State[x]...)
	}
	return a0
}

// NewAnagrams returns the Anagrams for the words.
func NewAnagrams(words []string) Anagrams {
	a := make(map[string][]string)
	for _, word := range words {
		k := []byte(word)
		sort.Slice(k, func(i, j int) bool { return k[i] < k[j] })
		k_str := string(k)
		a[k_str] = append(a[k_str], word)
	}
	return Anagrams{State: a}
}

// The String method gives sorted lists of anagrams on each line.
func (a *Anagrams) String() string {
	ret := ""
	for _, words := range a.State {
		if len(words) > 1 {
			sort.Slice(words, func(i, j int) bool { return words[i] < words[j] })
			ret += fmt.Sprintln(strings.Join(words, " "))
		}
	}
	return ret
}
