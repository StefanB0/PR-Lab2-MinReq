package consumer

import "log"

const (
	COOK_NR            = 10
	MAX_CAPACITY       = COOK_NR
	LOCAL_PORT         = ":8003"
	AGGREGATOR_PORT    = ":8002"
	AGGREGATOR_ADDRESS = "http://localhost:8002"
)

var (
	reviewsNr          float32
	reviewsTotalScore  float32
	avgScore           float32
	newOrderChannel    = make(chan Order, MAX_CAPACITY)
	orderReviewChannel = make(chan OrderReview, 100)
)

func Run() {
	SendRequest(OrderRequest{Quantity: MAX_CAPACITY}, AGGREGATOR_ADDRESS+"/order/forward/request")
	go countScore()
	
	for i := 0; i < COOK_NR; i++ {
		cook := NewCook(i, AGGREGATOR_ADDRESS, newOrderChannel)
		go cook.Run()
	}
	
	StartServer(LOCAL_PORT, newOrderChannel, orderReviewChannel)
}

func countScore() {
	for {
		review := <-orderReviewChannel
		reviewsNr++
		reviewsTotalScore += float32(review.Score)
		avgScore = reviewsTotalScore / reviewsNr
		if int(reviewsNr)%100 == 0 {
			log.Println("avg score", avgScore)
		}
	}
}
