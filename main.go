package main

import (
	"fmt"
	"sync"
)

var msg string
var wg sync.WaitGroup

func updateMessage(s string, m *sync.Mutex) {
	defer wg.Done()
	m.Lock()
	msg = s
	m.Unlock()
}
func main() {

	// var mutex sync.Mutex

	// msg = "Hello, world!"
	// wg.Add(2)
	// go updateMessage("Hello, Universe!", &mutex)
	// go updateMessage("Hello, cosmos!", &mutex)

	// wg.Wait()

	// fmt.Println(msg)

	var journal Journey
	journal.AddEntry("Hello World!")
	fmt.Println(journal.entries)

}

var entryCount int

type Journey struct {
	entries []string
}

func (j *Journey) AddEntry(text string) int {
	entryCount++
	entry := fmt.Sprintf("This entry is %d and the text is %s", entryCount, text)
	j.entries = append(j.entries, entry)
	return entryCount
}
