package main

import (
	"time"

	"github.com/fatih/color"
)

type BarberShopWaitingClient struct {
	ClientName string
}

type BarberShop struct {
	ShopCapacity    int
	HairCutDuration time.Duration
	NumberOfBarbers int
	BarberDoneChan  chan bool
	ClientChan      chan string
	ShopStatus      string
	WaitingClient   *[]BarberShopWaitingClient
}

func NewBarberShop(cChan chan string, dChan chan bool) *BarberShop {
	w := make([]BarberShopWaitingClient, 0)

	return &BarberShop{
		ShopCapacity:    seatingCapacity,
		HairCutDuration: cutDuration,
		NumberOfBarbers: 0,
		ClientChan:      cChan,
		BarberDoneChan:  dChan,
		ShopStatus:      "OPEN",
		WaitingClient:   &w,
	}
}

func (shop *BarberShop) ClientWaitingQueueLength() int {
	return len(shop.ClientChan)
}

func (shop *BarberShop) addBarber(barberName string) {
	shop.NumberOfBarbers++
	go func() {
		barber := NewBarber(barberName)
		color.Yellow("%s goes to the waiting room to check for clients.", barber.Name)

		for {
			if q := shop.ClientWaitingQueueLength(); q == 0 {
				color.Yellow("Cuurent queue length is %d \nCheck customer queue and got zero, so barber %s take a nap!", q, barber.Name)
				barber.TakeANap()
			}
		}

	}()
}
