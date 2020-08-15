package main

import (
	"bufio"
	"reflect"
	"fmt"
	"testing"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
)

type AnagramKey int32

// An Anagrams instance stores the anagrams for a word list.
type Anagrams struct {
	State map[AnagramKey][]string
	InsertionOrder []AnagramKey
	InvalidKeys map[AnagramKey]struct{}
}

func main() {
	if len(os.Args) == 1 {
		log.Fatal("Invocation: anagrams.go [words_file]")
	}

	_ = ioutil.ReadFile
	// b, err := ioutil.ReadFile(os.Args[1])
	// if err != nil {
	// 	log.Fatal("error reading anagrams file:", err)
	// }
	r, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal("error reading anagrams file:", err)
	}
	s := bufio.NewScanner(r)
	lines := []string{}
	for s.Scan() {
		lines = append(lines, s.Text())
	}
	// lines := strings.Split(string(b), "\n")
	// lines = lines[:len(lines)-1] // remove empty string from trailing newline
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

func max(x, y byte) byte {
	if x > y {
		return x
	}
	return y
}

func min(x, y byte) byte {
	if x < y {
		return x
	}
	return y
}

// NewAnagrams returns the Anagrams for the words.
func NewAnagrams(words []string) Anagrams {
	ana := make(map[AnagramKey][]string)
	io := []AnagramKey{}
	for _, word := range words {
		chars := []byte(word)
		chars = append(chars, '\n')
		a := chars[0]
		b := chars[1]
		c := chars[2]
		d := chars[3]

		e := max(a, c)
		f := min(a, c)

		g := max(b, d)
		h := min(b, d)

		i := max(e, g)
		j := min(e, g)

		k := max(f, h)
		l := min(f, h)

		m := max(j, k)
		n := min(j, k)

		k_str := ((AnagramKey(i) << 24) | (AnagramKey(m) << 16) | (AnagramKey(n) << 8) | (AnagramKey(l) << 0))

		_, ok := ana[k_str]
		if !ok {
			io = append(io, k_str)
		}
		ana[k_str] = append(ana[k_str], word)

	}
	return Anagrams{State: ana, InsertionOrder: io}
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

func (a *Anagrams) Filter() *Anagrams {
	a.InvalidKeys = make(map[AnagramKey]struct{})
	for k, v := range a.State {
		if len(v) == 1 {
			delete(a.State, k)
			a.InvalidKeys[k] = struct{}{}
		}
	}
	return a
}

func (a *Anagrams) SortedValues() [][]string {
	ret := [][]string{}
	for _, k := range a.InsertionOrder {
		_, ok := a.InvalidKeys[k]
		if !ok {
			ret = append(ret, a.State[k])
		}
	}
	return ret
}

func TestSample(t *testing.T) {
	sample := []string{"act", "cat", "tree", "race", "care", "acre", "bee"}
	expected := [][]string{[]string{"act", "cat"}, []string{"race", "care", "acre"}}
	a := NewAnagrams(sample)
	a.Filter()
	actual := a.SortedValues()
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("%v != %v", actual, expected)
	}
}

// func TestLongerWords(t *testing.T) {
// 	w := []string{"silent", "dreads", "gainly", "laying", "fried", "fired", "equal", "listen", "sadder", "reply", "title", "final"}
// 	expected := [][]string{[]string{"silent", "listen"}, []string{"dreads", "sadder"}, []string{"gainly", "laying"}, []string{"fried", "fired"}}
// 	a := NewAnagrams(w)
// 	a.Filter()
// 	actual := a.SortedValues()
// 	if !reflect.DeepEqual(expected, actual) {
// 		t.Errorf("%v != %v", actual, expected)
// 	}
// }

func TestCaseSensitive(t *testing.T) {
	w := []string{"God", "dog"}
	expected := [][]string{}
	a := NewAnagrams(w)
	a.Filter()
	actual := a.SortedValues()
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("%v != %v", actual, expected)
	}
}

func TestSpaceSensitive(t *testing.T) {
	w := []string{"tom marvolo riddle", "i am lord voldemort"}
	expected := [][]string{}
	a := NewAnagrams(w)
	a.Filter()
	actual := a.SortedValues()
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("%v != %v", actual, expected)
	}
}
