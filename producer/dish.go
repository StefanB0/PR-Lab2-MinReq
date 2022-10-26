package producer

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
)

type Dish struct {
	Id               int    `json:"order_id"`
	Name             string `json:"name"`
	Preparation_time int    `json:"prep_time"`
}

func GenerateDish() Dish {
	l := len(Menu)
	r := rand.Intn(l)
	dish := Menu[r]
	return dish
}

func ParseMenu(path string) []Dish {
	jsonfile, err := os.Open(path)
	defer jsonfile.Close()

	if err != nil {
		log.Println(err)
	}

	bytevalue, _ := ioutil.ReadAll(jsonfile)
	newMenu := []Dish{}
	json.Unmarshal(bytevalue, &newMenu)

	return newMenu
}
