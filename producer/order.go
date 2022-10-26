package producer

import "time"

type Order struct {
	Id        int       `json:"order_id"`
	Client_id int       `json:"client_id"`
	Cooked    bool      `json:"cook_value"`
	Dish      Dish      `json:"dish"`
	OrderTime time.Time `json:"time"`
}

type OrderRequest struct {
	Quantity int `json:"quantity"`
}

type OrderReview struct {
	Id    int `json:"id"`
	Score int `json:"score"`
}

func GenerateOrder(_id, _client_id int, _address chan Order) Order {
	order := &Order{
		Id:        _id,
		Client_id: _client_id,
		Cooked:    false,
		OrderTime: time.Now(),
	}

	order.Dish = GenerateDish()

	return *order
}

func ProcessOrderRequest(r OrderRequest, producerChan chan struct{}) {
	for i := 0; i < r.Quantity; i++ {
		producerChan <- struct{}{}
	}
}

func MatchOrder(order_a, order_b Order) bool {
	var match bool
	match = order_a.Id == order_b.Id
	return match
}
