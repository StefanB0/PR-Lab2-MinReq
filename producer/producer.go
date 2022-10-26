package producer

import (
	"time"
)

const (
	CUSTOMER_NR        = 1000
	LOCAL_PORT         = ":8001"
	AGGREGATOR_PORT    = ":8002"
	AGGREGATOR_ADDRESS = "http://localhost:8002"
	RUNTIME            = time.Millisecond * 1
	DELAY_LOWER_BOUND  = CUSTOMER_NR
	DELAY_VARIATION    = 20
)

var (
	ProducerChan = make(chan struct{}, CUSTOMER_NR)
	CostumerMap  = make(map[int]chan Order, CUSTOMER_NR)
	Menu         []Dish
)

func Run() {
	Menu = ParseMenu("producer/menu.json")

	customer_counter := &Counter{}
	order_counter := &Counter{}

	for i := 0; i < CUSTOMER_NR; i++ {
		customer := NewCustomer(customer_counter.Increment(), AGGREGATOR_ADDRESS, order_counter, ProducerChan)
		CostumerMap[customer.id] = customer.address
		go customer.Run()
	}
	StartServer(LOCAL_PORT, ProducerChan)
}
