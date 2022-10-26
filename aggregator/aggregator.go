package aggregator

import (
	"time"
)

const (
	RECEPTION_NR        = 5
	DELIVERY_NR         = 5
	DELIVERY_TIME       = 30
	CHAIN_DELIVERY_TIME = 10
	PRODUCER_PORT       = ":8001"
	LOCAL_PORT          = ":8002"
	CONSUMER_PORT       = ":8003"
	PRODUCER_ADDRESS    = "http://localhost:8001"
	CONSUMER_ADDRESS    = "http://localhost:8003"
	RUNTIME             = time.Millisecond * 1
)

var (
	newOrderChannel     = make(chan Order, 100)
	orderRequestChannel = make(chan OrderRequest, 100)
	orderReviewChannel  = make(chan OrderReview, 100)
	finOrderChannel     = make(chan Order, 100)
)

func Run() {

	for i := 0; i < RECEPTION_NR; i++ {
		r := newReceptionist(i, PRODUCER_ADDRESS, CONSUMER_ADDRESS, newOrderChannel, orderRequestChannel, orderReviewChannel)
		go r.Run()
	}

	for i := 0; i < DELIVERY_NR; i++ {
		d := NewDeliveryMan(i, PRODUCER_ADDRESS, finOrderChannel)
		go d.Run()
	}

	StartServer(LOCAL_PORT, newOrderChannel, orderRequestChannel, orderReviewChannel, finOrderChannel)
}
