package main

import (
	"fmt"
	"time"
)

func server1(ch chan string) {
	for {
		time.Sleep(4 * time.Second)
		ch <- fmt.Sprintf("Sending from server 1")
	}

}

func server2(ch chan string) {

	for {
		time.Sleep(3 * time.Second)
		ch <- fmt.Sprintf("Sending from server 2")

	}

}

func main() {

	ch1 := make(chan string)
	ch2 := make(chan string)

	go server1(ch1)
	go server2(ch2)

	for {
		select {
		case s1 := <-ch1:
			fmt.Println("Case one", s1)
		case s2 := <-ch1:
			fmt.Println("Case two", s2)
		case s3 := <-ch2:
			fmt.Println("Case three", s3)
		case s4 := <-ch2:
			fmt.Println("Case four", s4)
		}
	}

}
