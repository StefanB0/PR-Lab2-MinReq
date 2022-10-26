package main

import (
	"PR/Lab2/aggregator"
	"PR/Lab2/consumer"
	"PR/Lab2/producer"
)

func main() {
	go producer.Run()
	go aggregator.Run()
	go consumer.Run()
	for {
	}
}
