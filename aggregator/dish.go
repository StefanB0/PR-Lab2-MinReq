package aggregator

type Dish struct {
	Id               int    `json:"order_id"`
	Name             string `json:"name"`
	Preparation_time int    `json:"prep_time"`
}
