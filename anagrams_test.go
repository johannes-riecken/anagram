package main

import (
	"flag"
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

type OrderedMap struct {
	Map map[string][]string
	Keys []string
	i int
}

func (o *OrderedMap) Value() []string {
	return o.Map[o.Keys[o.i]]
}

func (o *OrderedMap) Front() *OrderedMap {
	o.i = 0
	return o
}

func (o *OrderedMap) Next() *OrderedMap {
	o.i++
	if o.i < len(o.Keys) {
		return o
	}
	return nil
}

func (o *OrderedMap) Insert(k string, v []string) {
	_, ok := o.Map[k]
	o.Map[k] = v
	if !ok {
		o.Keys = append(o.Keys, k)
	}
}

func (o *OrderedMap) AppendValues(k string, v []string) {
	_, ok := o.Map[k]
	o.Map[k] = append(o.Map[k], v...)
	if !ok {
		o.Keys = append(o.Keys, k)
	}
}

func lines(s bufio.Scanner) <-chan string {
	ch := make(chan string)
	for s.Scan() {
		ch <- s.Text()
	}
	close(ch)
	return ch
}

// func anagrams(a []string, co chan OrderedMap) {
// 	m := OrderedMap{}
// 	for _, w := range a {
// 		b := []byte(w)
// 		sort.Slice(b, func(i, j int) bool {return b[i] < b[j];})
// 		m.AppendValues(string(b), []string{w})
// 	}
// 	return ch
// }

// func merge(chans []chan OrderedMap) chan OrderedMap {
// 	co := make(chan OrderedMap)
// 	for ch := range chans {
// 	}
// 	return co
// }


func main() {
	_ = flag.Int
	_ = ioutil.ReadFile
	jobs := flag.Int("j", 1, "Number of jobs to run simultaneously")
	if len(os.Args) == 1 {
		log.Fatal("Invocation: anagrams.go [words_file]")
	}

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

func TestPseudoParallel(t *testing.T) {
	w0 := []string{"act", "cat", "tree", "race"}
	w1 := []string{"silent"}
	w2 := []string{"care"}
	w3 := []string{"acre", "bee"}
	a0 := NewAnagrams(w0)
	a1 := NewAnagrams(w1)
	a2 := NewAnagrams(w2)
	a3 := NewAnagrams(w3)
	a0.Merge(a1)
	a0.Merge(a2)
	a0.Merge(a3)
	a0.Filter()
	actual := a0.SortedValues()
	expected := [][]string{[]string{"act", "cat"}, []string{"race", "care", "acre"}}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("%v != %v", actual, expected)
	}
}

func TestParallel(t *testing.T) {
	// w0 := []string{"act", "cat", "tree", "race"}
	// w1 := []string{"silent"}
	// w2 := []string{"care"}
	// w3 := []string{"acre", "bee"}
	// a0 := NewAnagrams(w0)
	// a1 := NewAnagrams(w1)
	// a2 := NewAnagrams(w2)
	// a3 := NewAnagrams(w3)
	// each file partition has a job to put lines in its own channel
	// each above job is connected to an anagram loop job
	// a merge job takes all anagram channels and does a while select until all
	// are closed
	// then the merge job returns all info to main, which calls a print function
}
