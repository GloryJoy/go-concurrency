package main

import "github.com/fatih/color"

// Check if the shop is open, if not leave
//
// Check if the shop has seat available, if not leave. if yes,
// Check if the barbar is sleeping or not, if yes wake the barber
// If not, take a seat and wait

type Client struct {
	Name   string
	Status string // Left, Waiting, GettingHaircut
}

func NewClient(name string, barberShop *BarberShop) *Client {

	return &Client{
		Name: name,
	}
}

func (c *Client) CheckShopOpen(shopClientChan chan string, clientShopChan chan string) string {
	clientShopChan <- "Is Open?"
	color.Red("[Client] Client is checking if the shop is opened...")
	shopStatus := <-shopClientChan
	color.Red("[Client] Get answer back is %s", shopStatus)
	return shopStatus
}

func (c *Client) CheckBarber(barberClient chan string, clientBarber chan string) string {
	clientBarber <- "Is sleeping?"
	color.Red("[Client] asking barber if sleeping...")
	barberStatus := <-barberClient
	color.Red("[Client] got the answer = %s", barberStatus)
	return barberStatus

}

func (c *Client) WakeUpBarber(barberClient chan string, clientBarber chan string) string {
	clientBarber <- "Wake Up!"
	color.Red("[Client] Waked up barder")
	// barberStatus := <-barberClient
	barberStatus := <-barberClient
	color.Red("[Client] barder status %s", barberStatus)
	return barberStatus
}
