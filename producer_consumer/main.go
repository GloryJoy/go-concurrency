package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

const NumberOfPizzas = 10

var pizzasMade, pizzasFailed, total int

type Producer struct {
	data chan PizzaOrder
	quit chan chan error
}

func (p *Producer) Close() error {
	ch := make(chan error)
	p.quit <- ch
	return <-ch

}

type PizzaOrder struct {
	pizzaNumber int
	message     string
	success     bool
}

func pizzaria(pizzaMaker *Producer) {
	// keep track of which pizza we are making
	var currentNumberOfMadePizza = 0

	// run forever or until we receive a quit notification

	// try to make pizzas

	for {
		currentPizza := makePizza(currentNumberOfMadePizza)
		if currentPizza != nil {
			currentNumberOfMadePizza = currentPizza.pizzaNumber
			select {
			case pizzaMaker.data <- *currentPizza:
			case quitChan := <-pizzaMaker.quit:
				close(pizzaMaker.data)
				close(quitChan)
				return
			}
		}
		// try to make a pizza
		// decision
	}
}

func makePizza(pizzaNumber int) *PizzaOrder {
	pizzaNumber++
	msg := ""
	success := false

	if pizzaNumber <= NumberOfPizzas {
		delay := rand.Intn(5) + 1
		fmt.Printf("Received order #%d\n", pizzaNumber)
		rnd := rand.Intn(12) + 1

		if rnd < 5 {
			pizzasFailed++
		} else {
			pizzasMade++
		}
		total++
		fmt.Printf("Making pizza #%d. It will take %d seconds....\n", pizzaNumber, delay)
		// delay for a bit

		time.Sleep(time.Duration(delay) * time.Second)

		if rnd <= 2 {
			msg = fmt.Sprintf("*** We ran out of ingredients for pizza #%d!\n", pizzaNumber)
		} else if rnd <= 4 {
			msg = fmt.Sprintf("*** The cook quit while making pizza", pizzaNumber)
		} else {
			success = true
			msg = fmt.Sprintf("Pizza order #%d is ready!", pizzaNumber)
		}

		p := PizzaOrder{
			pizzaNumber: pizzaNumber,
			message:     msg,
			success:     success,
		}

		return &p

	}
	return &PizzaOrder{
		pizzaNumber: pizzaNumber,
	}
}

func main() {
	//seed the random number generator
	rand.Seed(time.Now().UnixNano())
	// print out a message
	color.Cyan("The Pizzaria is open for business!")
	color.Cyan("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")

	// create a producer
	pizzaJob := &Producer{
		data: make(chan PizzaOrder),
		quit: make(chan chan error),
	}

	// run the producer in the background
	go pizzaria(pizzaJob)
	// create and run consumer
	for i := range pizzaJob.data {
		if i.pizzaNumber <= NumberOfPizzas {
			if i.success {
				color.Green(i.message)
				color.Green("Order #%d is out for delivery!", i.pizzaNumber)
			} else {
				color.Red(i.message)
				color.Red("Customer is really mad!")
			}
		} else {
			color.Cyan("Done making pizza....")
			err := pizzaJob.Close()
			if err != nil {
				color.Red("****Error", err)
			}
		}
	}

	// print out the ending message
	color.Cyan("Done for the day.")
	color.Cyan("We made %d pizzas, but failed to make %d, with %d attemps in total.", pizzasMade, pizzasFailed, total)
	switch {
	case pizzasFailed > 9:
		color.Red("It was an awful day...")
	case pizzasFailed >= 6:
		color.Red("It was not a very good day...")
	case pizzasFailed >= 4:
		color.Yellow("It was an okkay day...")

	}

}
