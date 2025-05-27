package nats

import (
	"encoding/json"
	"log"
	
	"github.com/nats-io/nats.go"
)
type NatsPublisher struct {
	conn *nats.Conn
}

var natsConn *nats.Conn

func NewPublisher() *NatsPublisher {
	var err error
	natsConn, err = nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	log.Println("Connected to NATS")
	return &NatsPublisher{conn: natsConn}
}

type OrderCreatedEvent struct {
	ID        string `json:"id"`
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}


func (p *NatsPublisher) PublishOrderCreated(event OrderCreatedEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = p.conn.Publish("order.created", data)
	if err != nil {
		log.Println("Failed to publish message:", err)
		return err
	}

	log.Println("Published order.created event:", event.ID)
	return nil
}
