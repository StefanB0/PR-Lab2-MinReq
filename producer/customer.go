package producer

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"time"
)

type Customer struct {
	id                 int
	order              Order
	address            chan Order
	producerChan       chan struct{}
	aggregator_address string
	order_counter      *Counter
}

func NewCustomer(_id int, _aggregator_address string, _order_counter *Counter, _producerChan chan struct{}) *Customer {
	customer := &Customer{
		id:                 _id,
		address:            make(chan Order),
		aggregator_address: _aggregator_address,
		order_counter:      _order_counter,
		producerChan:       _producerChan,
	}
	return customer
}

func (c *Customer) Run() {
	for {
		<-c.producerChan
		c.order = GenerateOrder(c.order_counter.Increment(), c.id, c.address)
		c.sendOrder(c.order, c.aggregator_address+"/order/new")
		c.waitOrder()
		OrderPause()
	}
}

func (c *Customer) waitOrder() {
	finished_food := <-c.address
	var score int
	if !MatchOrder(c.order, finished_food) {
		score = 0
	} else {
		score = 5 - (int(time.Now().Sub(finished_food.OrderTime) / RUNTIME) - c.order.Dish.Preparation_time)
	}

	c.sendReview(c.order.Id, score, c.aggregator_address+"/order/review")
}

func (c *Customer) sendOrder(newOrder Order, url string) {
	payloadBuffer := new(bytes.Buffer)
	json.NewEncoder(payloadBuffer).Encode(newOrder)

	req, _ := http.NewRequest("POST", url, payloadBuffer)
	client := &http.Client{}
	client.Do(req)
}

func (c *Customer) sendReview(_id, _score int, url string) {
	review := &OrderReview{Id: _id, Score: _score}

	payloadBuffer := new(bytes.Buffer)
	json.NewEncoder(payloadBuffer).Encode(review)

	req, _ := http.NewRequest("POST", url, payloadBuffer)
	client := &http.Client{}
	client.Do(req)
}

func SendToCustomer(order Order) {
	CostumerMap[order.Client_id] <- order
}

func OrderPause() {
	variation := rand.Intn(DELAY_VARIATION)
	time.Sleep(RUNTIME * time.Duration(variation+DELAY_LOWER_BOUND))
}
