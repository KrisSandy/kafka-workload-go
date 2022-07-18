package consumer

import (
	"context"
	"example/kafka-workload/package/utils"
	"fmt"
	"log"
	"strings"

	"github.com/segmentio/kafka-go"
)

func newKafkaReader(kafkaURL string, topic string, groupID string) *kafka.Reader {
	brokers := strings.Split(kafkaURL, ",")
	clientCertFile := "/Users/sandy/Dev/certs/client.crt"
	clientKeyFile := "/Users/sandy/Dev/certs/client.key"
	caCertFile := "/Users/sandy/Dev/certs/server.crt"
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Dialer:   utils.GetDialer(clientCertFile, clientKeyFile, caCertFile),
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
	// return kafka.NewReader(kafka.ReaderConfig{
	// 	Brokers:  brokers,
	// 	GroupID:  groupID,
	// 	Topic:    topic,
	// 	MinBytes: 10e3, // 10KB
	// 	MaxBytes: 10e6, // 10MB
	// })
}

func Consumer(kafkaURL string, topic string) {
	groupID := ""

	reader := newKafkaReader(kafkaURL, topic, groupID)

	defer reader.Close()

	fmt.Println("Starting consuming...")

	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("message at topic:%v partition:%v offset:%v	%s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
	}
}
