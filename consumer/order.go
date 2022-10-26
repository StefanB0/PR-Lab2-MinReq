package consumer

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

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

func SendRequest(request OrderRequest, url string) {
	payloadBuffer := new(bytes.Buffer)
	json.NewEncoder(payloadBuffer).Encode(request)

	req, _ := http.NewRequest("POST", url, payloadBuffer)
	client := &http.Client{}
	client.Do(req)
}
