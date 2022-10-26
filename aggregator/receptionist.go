package aggregator

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Receptionist struct {
	id              int
	inputChannel    chan Order
	consumerRequest chan OrderRequest
	reviewChannel   chan OrderReview
	producerUrl     string
	consumerUrl     string
}

func newReceptionist(
	_id int,
	_producerUrl string,
	_consumerUrl string,
	_inputChannel chan Order,
	_consumerRequest chan OrderRequest,
	_reviewChannel chan OrderReview,
) *Receptionist {
	return &Receptionist{
		id:              _id,
		producerUrl:     _producerUrl,
		consumerUrl:     _consumerUrl,
		inputChannel:    _inputChannel,
		consumerRequest: _consumerRequest,
		reviewChannel:   _reviewChannel,
	}
}

func (r *Receptionist) Run() {
	for {
		select {
		case order := <-r.inputChannel:
			r.forwardOrder(order, r.consumerUrl+"/order/new")
		case request := <-r.consumerRequest:
			r.forwardRequest(request, r.producerUrl+"/order/request")
		case review := <-r.reviewChannel:
			r.forwardReview(review, r.consumerUrl+"/order/review")
		}
	}
}

func (r *Receptionist) forwardOrder(order Order, url string) {
	payloadBuffer := new(bytes.Buffer)
	json.NewEncoder(payloadBuffer).Encode(order)

	req, _ := http.NewRequest("POST", url, payloadBuffer)
	client := &http.Client{}
	client.Do(req)
}
func (r *Receptionist) forwardRequest(request OrderRequest, url string) {
	payloadBuffer := new(bytes.Buffer)
	json.NewEncoder(payloadBuffer).Encode(request)

	req, _ := http.NewRequest("POST", url, payloadBuffer)
	client := &http.Client{}
	client.Do(req)
}
func (r *Receptionist) forwardReview(review OrderReview, url string) {
	payloadBuffer := new(bytes.Buffer)
	json.NewEncoder(payloadBuffer).Encode(review)

	req, _ := http.NewRequest("POST", url, payloadBuffer)
	client := &http.Client{}
	client.Do(req)
}
