package producer

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/spiffe/go-spiffe/v2/spiffeid"
	"github.com/spiffe/go-spiffe/v2/spiffetls/tlsconfig"
	"github.com/spiffe/go-spiffe/v2/workloadapi"
)

// func newKafkaWriter(ctx context.Context, kafkaURL string, topic string) *kafka.Writer {
// clientCertFile := "/Users/sandy/Dev/certs/client.crt"
// clientKeyFile := "/Users/sandy/Dev/certs/client.key"
// caCertFile := "/Users/sandy/Dev/certs/server.crt"
// return kafka.NewWriter(kafka.WriterConfig{
// 	Brokers:  []string{kafkaURL},
// 	Topic:    topic,
// 	Balancer: &kafka.Hash{},
// 	// Dialer:   utils.GetDialer(clientCertFile, clientKeyFile, caCertFile),
// 	Dialer: utils.GetSpiffeDialer(ctx),
// })
// return &kafka.Writer{
// 	Addr:     kafka.TCP(kafkaURL),
// 	Topic:    topic,
// 	Balancer: &kafka.LeastBytes{},
// }
// }

func Producer(kafkaURL string, topic string) {
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

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{kafkaURL},
		Topic:    topic,
		Balancer: &kafka.Hash{},
		Dialer:   dialer,
	})
	defer writer.Close()
	log.Println("Producer started")
	for i := 0; ; i++ {
		key := fmt.Sprintf("Key-%d", i)
		msg := kafka.Message{
			Key:   []byte(key),
			Value: []byte(fmt.Sprint(uuid.New())),
		}
		err := writer.WriteMessages(ctx, msg)
		if err != nil {
			log.Fatal("Writing message to producer failed : ", err)
		} else {
			log.Println("Produced key - ", key)
		}
		time.Sleep(5 * time.Second)
	}
}
