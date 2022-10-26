package aggregator

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
