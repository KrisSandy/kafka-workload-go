package consumer

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/spiffe/go-spiffe/v2/spiffeid"
	"github.com/spiffe/go-spiffe/v2/spiffetls/tlsconfig"
	"github.com/spiffe/go-spiffe/v2/workloadapi"
)

// func newKafkaReader(kafkaURL string, topic string, groupID string) *kafka.Reader {
// 	brokers := strings.Split(kafkaURL, ",")
// 	clientCertFile := "/Users/sandy/Dev/certs/client.crt"
// 	clientKeyFile := "/Users/sandy/Dev/certs/client.key"
// 	caCertFile := "/Users/sandy/Dev/certs/server.crt"
// 	return kafka.NewReader(kafka.ReaderConfig{
// 		Brokers:  brokers,
// 		GroupID:  groupID,
// 		Dialer:   utils.GetDialer(clientCertFile, clientKeyFile, caCertFile),
// 		Topic:    topic,
// 		MinBytes: 10e3, // 10KB
// 		MaxBytes: 10e6, // 10MB
// 	})
// 	// return kafka.NewReader(kafka.ReaderConfig{
// 	// 	Brokers:  brokers,
// 	// 	GroupID:  groupID,
// 	// 	Topic:    topic,
// 	// 	MinBytes: 10e3, // 10KB
// 	// 	MaxBytes: 10e6, // 10MB
// 	// })
// }

func Consumer(kafkaURL string, topic string) {
	groupID := ""
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	socketPath := "unix:///tmp/spire-agent/public/api.sock"
	x509Source, err := workloadapi.NewX509Source(ctx, workloadapi.WithClientOptions(workloadapi.WithAddr(socketPath)))
	if err != nil {
		log.Fatal("Unable to create X509Source: ", err)
	}
	defer x509Source.Close()

	serverID := spiffeid.RequireFromString("spiffe://example.org/myservice")

	config := tlsconfig.MTLSClientConfig(x509Source, x509Source, tlsconfig.AuthorizeID(serverID))

	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
		TLS:       config,
	}

	brokers := strings.Split(kafkaURL, ",")
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Dialer:   dialer,
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

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
