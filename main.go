package main

import (
	"example/kafka-workload/package/consumer"
	"example/kafka-workload/package/producer"
	"log"
	"os"
)

func main() {
	entry := os.Args[1]
	kafkaURL := "localhost:9093"
	topic := "test-topic"

	if entry == "producer" {
		log.Println("Starting producer")
		producer.Producer(kafkaURL, topic)
	}

	if entry == "consumer" {
		log.Println("Starting consumer")
		consumer.Consumer(kafkaURL, topic)
	}
}
