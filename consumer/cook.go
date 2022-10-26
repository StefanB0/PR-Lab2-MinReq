package consumer

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

type Cook struct {
	id            int
	orderChannel  chan Order
	aggregatorUrl string
}

func NewCook(_id int, _aggregatorUrl string, _order_channel chan Order) *Cook {
	return &Cook{id: _id, aggregatorUrl: _aggregatorUrl, orderChannel: _order_channel}
}

func (c *Cook) Run() {
	for {
		order := <-c.orderChannel
		finished_order := c.prepareOrder(order)
		c.sendOrder(finished_order, c.aggregatorUrl+"/order/finished")
		SendRequest(OrderRequest{Quantity: 1}, AGGREGATOR_ADDRESS+"/order/forward/request")
	}
}

func (c *Cook) prepareOrder(order Order) Order {
	time.Sleep(time.Millisecond * time.Duration(order.Dish.Preparation_time))
	order.Cooked = true
	return order
}

func (c *Cook) sendOrder(order Order, url string) {
	payloadBuffer := new(bytes.Buffer)
	json.NewEncoder(payloadBuffer).Encode(order)

	req, _ := http.NewRequest("POST", url, payloadBuffer)
	client := &http.Client{}
	client.Do(req)
}
