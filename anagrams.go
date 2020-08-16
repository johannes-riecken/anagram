package main

import (
	"sync"
	"flag"
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
)

type AnagramKey string
type AnagramValue []string

type OrderedMap struct {
	m map[AnagramKey]AnagramValue
	Keys []AnagramKey
	i int
}

func NewOrderedMap() OrderedMap {
	o := OrderedMap{}
	o.m = make(map[AnagramKey]AnagramValue)
	return o
}

func (o *OrderedMap) Key() AnagramKey {
	return o.Keys[o.i]
}

func (o *OrderedMap) Value() AnagramValue {
	return o.m[o.Keys[o.i]]
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

func (o *OrderedMap) Insert(k AnagramKey, v AnagramValue) {
	_, ok := o.m[k]
	o.m[k] = v
	if !ok {
		o.Keys = append(o.Keys, k)
	}
}

func (o *OrderedMap) AppendValues(k AnagramKey, v AnagramValue) {
	_, ok := o.m[k]
	o.m[k] = append(o.m[k], v...)
	if !ok {
		o.Keys = append(o.Keys, k)
	}
}

func lines(f string) <-chan string {
	ch := make(chan string)
	go func() {
		r, err := os.Open(f)
		if err != nil {
			log.Fatal("error reading anagrams file:", err)
		}
		s := bufio.NewScanner(r)
		for s.Scan() {
			ch <- s.Text()
		}
		close(ch)
	}()
	return ch
}

func anagrams(a <-chan string) <-chan OrderedMap {
	ch := make(chan OrderedMap)
	go func() {
		o := NewOrderedMap()
		for w := range a {
			b := []byte(w)
			sort.Slice(b, func(i, j int) bool {return b[i] < b[j];})
			o.AppendValues(AnagramKey(b), AnagramValue{w})
		}
		ch <- o
		close(ch)
	}()
	return ch
}

func merge(chans []<-chan OrderedMap) OrderedMap {
	o := NewOrderedMap()
	out := make(chan OrderedMap)
	wg := sync.WaitGroup{}
	wg.Add(len(chans))
	for _, ch := range chans {
		go func(c <-chan OrderedMap) {
			out <- <-c
			wg.Done()
		}(ch)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	for oo := range out {
		for it := oo.Front(); it != nil; it = it.Next() {
			o.AppendValues(it.Key(), it.Value())
		}
	}
	return o
}


func main() {
	_ = flag.Int
	_ = ioutil.ReadFile
	_ = os.Open
	_ = log.Fatal
	_ = bufio.NewScanner
	if len(os.Args) == 1 {
		log.Fatal("Invocation: anagrams.go [words_file]")
	}

	chans := make([]<-chan OrderedMap, len(os.Args[1:]))
	for i, a := range os.Args[1:] {
		chans[i] = anagrams(lines(a))
	}
	o := merge(chans)
	for it := o.Front(); it != nil; it = it.Next() {
		v := it.Value()
		if len(v) > 1 {
			fmt.Println(strings.Join(v, " "))
		}
	}
}
