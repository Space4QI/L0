package pubsub

import (
	"L0/database"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"log"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func Publish(clientID string, clusterID string, orders database.Orders) error {
	ordersJSON, err := json.Marshal(orders)
	if err != nil {
		return fmt.Errorf("ошибка при маршалинге JSON: %w", err)
	}

	if err := validateJSON(ordersJSON); err != nil {
		return fmt.Errorf("данные не прошли валидацию: %w", err)
	}

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return fmt.Errorf("ошибка при подключении к NATS: %w", err)
	}
	defer nc.Close()

	sc, err := stan.Connect(clientID, clusterID, stan.NatsConn(nc), stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
		log.Printf("Connection lost, reason: %v", reason)
	}))
	if err != nil {
		return fmt.Errorf("ошибка при подключении к NATS Streaming: %w", err)
	}
	defer sc.Close()

	err = sc.Publish("orders", ordersJSON)
	if err != nil {
		return fmt.Errorf("ошибка при публикации сообщения: %w", err)
	}
	log.Println("Сообщение успешно опубликовано")
	return nil
}

func validateJSON(data []byte) error {
	var order database.Orders
	err := json.Unmarshal(data, &order)
	if err != nil {
		return err
	}

	// Валидация структуры Order на основе тегов json
	err = validate.Struct(order)
	if err != nil {
		return err
	}

	return nil
}
