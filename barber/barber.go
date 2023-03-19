package main

type Barber struct {
	Status string
	Name   string
}

func NewBarber(name string) *Barber {
	return &Barber{
		Status: "SLEEPING",
		Name:   name,
	}
}

func (barber *Barber) TakeANap() {
	barber.Status = "SLEEPING"
}
