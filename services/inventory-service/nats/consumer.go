package nats

import (
	"encoding/json"
	"log"
	"strconv"
	
	"github.com/nats-io/nats.go"
	"inventory-service/internal/usecase"
)

type Order struct {
	ID        string `json:"id"`
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

func InitNATSConsumer(usecase *usecase.ProductUsecase) {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	log.Println("Inventory connected to NATS")

	_, err = nc.Subscribe("order.created", func(m *nats.Msg) {
		var order Order
		if err := json.Unmarshal(m.Data, &order); err != nil {
			log.Println("Failed to parse order:", err)
			return
		}

		log.Println("Received order:", order.ID)

		productID, err := strconv.Atoi(order.ProductID)
		if err != nil {
			log.Println("Invalid product ID:", err)
			return
		}

		err = usecase.DecreaseStock(productID, int32(order.Quantity))
		if err != nil {
			log.Println("Failed to decrease stock:", err)
		} else {
			log.Println("Stock updated for product:", productID)
		}
	})
	if err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
	}
}
