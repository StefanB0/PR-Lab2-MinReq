package aggregator

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type DeliveryMan struct {
	id           int
	producerUrl  string
	finOrderChan chan Order
}

func NewDeliveryMan(_id int, _producerUrl string, _finOrderChan chan Order) *DeliveryMan {
	return &DeliveryMan{id: _id, producerUrl: _producerUrl, finOrderChan: finOrderChannel}
}

func (d *DeliveryMan) Run() {
	for {
		finOrder := <-finOrderChannel
		d.deliverFood(finOrder, d.producerUrl + "/order/finished")
	}
}

func (d *DeliveryMan) deliverFood(order Order, url string) {
	payloadBuffer := new(bytes.Buffer)
	json.NewEncoder(payloadBuffer).Encode(order)

	req, _ := http.NewRequest("POST", url, payloadBuffer)
	client := &http.Client{}
	client.Do(req)
}
