package producer

import (
	"context"
	"example/kafka-workload/package/utils"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

func newKafkaWriter(kafkaURL string, topic string) *kafka.Writer {
	clientCertFile := "/Users/sandy/Dev/certs/client.crt"
	clientKeyFile := "/Users/sandy/Dev/certs/client.key"
	caCertFile := "/Users/sandy/Dev/certs/server.crt"
	return kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{kafkaURL},
		Topic:    topic,
		Balancer: &kafka.Hash{},
		Dialer:   utils.GetDialer(clientCertFile, clientKeyFile, caCertFile),
	})
	// return &kafka.Writer{
	// 	Addr:     kafka.TCP(kafkaURL),
	// 	Topic:    topic,
	// 	Balancer: &kafka.LeastBytes{},
	// }
}

func Producer(kafkaURL string, topic string) {
	writer := newKafkaWriter(kafkaURL, topic)
	defer writer.Close()
	log.Println("Producer started")
	for i := 0; ; i++ {
		key := fmt.Sprintf("Key-%d", i)
		msg := kafka.Message{
			Key:   []byte(key),
			Value: []byte(fmt.Sprint(uuid.New())),
		}
		err := writer.WriteMessages(context.Background(), msg)
		if err != nil {
			log.Fatal("Writing message to producer failed", err)
		} else {
			log.Println("Produced key - ", key)
		}
		time.Sleep(5 * time.Second)
	}
}
