package main

import (
	"math/rand"
	"sync"
	"time"

	"github.com/fatih/color"
)

const (
	seatingCapacity = 10
	arrivalRate     = 100
	cutDuration     = 2 * time.Second
	timeOpen        = 10 * time.Second
)

var wg sync.WaitGroup

func main() {
	rand.Seed(time.Now().UnixNano())
	color.Yellow("The Sleeping Barber")
	color.Yellow("-------------------")

	clientChan := make(chan string, seatingCapacity)
	doneChan := make(chan bool)
	shopClientChan := make(chan string, 3)
	clientShopChan := make(chan string, 3)
	barberClientChan := make(chan string, 3)
	clientBarberChan := make(chan string, 3)

	shop := NewBarberShop(clientChan, doneChan)
	client := NewClient("Joy", shop)
	barber := NewBarber("Handsome")

	color.Green("The shop is open for the day")

	wg.Add(3)

	go func(s *BarberShop, sc chan<- string, cs <-chan string) {

		for {
			color.Green("[SHOP] Shop status is %s", s.ShopStatus)
			time.Sleep(1 * time.Second)
			color.Green("length of channel %d", len(cs))
			if len(cs) > 0 {

				clientInquiry := <-cs
				color.Green("[SHOP] Get client inquiry %s", clientInquiry)
				sc <- s.ShopStatus
			}
			// break

		}

	}(shop, shopClientChan, clientShopChan)

	go func(c *Client, sc chan string, cs chan string, bc chan string, cb chan string) {

		for {
			color.Red("[Client] ================ NEW CLIENT =============")

			for {

				if s := c.CheckShopOpen(sc, cs); s == "OPEN" {
					time.Sleep(100 * time.Millisecond)
					if cl := c.CheckBarber(bc, cb); cl == "SLEEPING" {
						time.Sleep(100 * time.Millisecond)
						color.Red("Barber is %s", cl)
						if s1 := c.WakeUpBarber(bc, cb); s1 == "WAKEDUP" {
							c.Status = "GETTINGHAIRCUT"
							color.Red("[Client] %s", c.Status)
							time.Sleep(cutDuration)
							c.Status = "DONE AND LEFT"
							color.Red("[Client] %s", c.Status)
							break
						}

					}
				}
				time.Sleep(1 * time.Second)
				color.Red("[Client] %s", c.Status)
			}
		}

	}(client, shopClientChan, clientShopChan, barberClientChan, clientBarberChan)

	go func(b *Barber, bc chan string, cb chan string) {

		for {
			time.Sleep(100 * time.Millisecond)

			if len(cb) > 0 {

				if clientInquiry := <-cb; clientInquiry == "Is sleeping?" {
					color.Blue("[BB] %s", clientInquiry)
					bc <- b.Status
					color.Blue("[BB] status %s", b.Status)
				}
				if clientInquiry := <-cb; clientInquiry == "Wake Up!" {
					color.Blue("[BB] %s", clientInquiry)
					b.Status = "WAKEDUP"
					bc <- b.Status
					color.Blue("[BB] status %s", b.Status)
					time.Sleep(cutDuration)
					b.Status = "SLEEPING"

				}
			}

			time.Sleep(100 * time.Millisecond)

		}

	}(barber, barberClientChan, clientBarberChan)

	wg.Wait()
	// close(shopClientChan)
	// close(doneChan)
	// close(clientChan)
	// close(barberClientChan)
	// close(clientShopChan)
	// wg.Done()

}
