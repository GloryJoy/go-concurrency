package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func shout(ping chan string, pong chan string) {
	for {
		s := <-ping
		pong <- fmt.Sprintf("%s!!!", strings.ToUpper(s))
	}
}
func main() {
	ping := make(chan string)
	pong := make(chan string)

	go shout(ping, pong)

	fmt.Println("Type something and then press Enter(enter q to quite)")

	inputScan := bufio.NewScanner(os.Stdin)
	var textLine string

	for {
		fmt.Print("-->")
		// var userInput string
		if inputScan.Scan() {
			textLine = inputScan.Text()
		}

		// _, _ = fmt.Scanln(&userInput)

		// if userInput == strings.ToLower()
		if s := strings.ToLower(textLine); s == "q" {
			break
		}

		ping <- textLine

		response := <-pong

		fmt.Println("Response is ", response)
	}

	fmt.Println("All done, closing channels")
	close(ping)
	close(pong)
}
