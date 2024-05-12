package pubsub

import (
	"L0/database"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"sync"
	"time"
)

var Wg sync.WaitGroup

const maxRetries = 3
const retryInterval = 5 * time.Second

func Subscribe(ctx context.Context, db *sql.DB, nc *nats.Conn) {
	fmt.Println("Подписка на канал 'orders'")

	done := make(chan struct{})

	sub, err := nc.Subscribe("orders", func(m *nats.Msg) {
		retryCount := 0
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if !json.Valid(m.Data) {
					log.Println("Получено невалидное JSON сообщение:", string(m.Data))
					return
				}
				log.Println("Получено сообщение:")

				// Преобразование среза байтов в строку
				jsonData := string(m.Data)
				fmt.Println("Содержимое сообщения:", jsonData)

				// Парсинг сообщения и сохранение в БД
				orders, err := parseOrdersJSON(jsonData)
				if err != nil {
					log.Println("Error parsing JSON:", err)
					return
				}
				fmt.Println("Parsed orders:", orders)

				// Начало транзакции
				tx, err := db.Begin()
				if err != nil {
					log.Println("Ошибка начала транзакции:", err)
					return
				}
				defer tx.Rollback()

				// Сохранение данных в БД
				err = database.SaveToDB(db, orders)
				if err != nil {
					log.Println("Ошибка сохранения данных в БД:", err)
					if retryCount < maxRetries {
						retryCount++
						log.Printf("Перезапуск обработки сообщения, попытка %d", retryCount)
						time.Sleep(retryInterval)
						continue
					}
					return
				}

				// Фиксация транзакции
				err = tx.Commit()
				if err != nil {
					log.Println("Ошибка фиксации транзакции:", err)
					return
				}

				fmt.Println("Данные успешно сохранены в БД")

				// Закрываем канал только после обработки сообщения
				close(done)
				return
			}
		}
	})
	if err != nil {
		log.Fatal("Ошибка подписки на канал:", err)
	}

	select {
	case <-done:
		fmt.Println("Подписка на канал 'orders' успешно остановлена")
	case <-ctx.Done():
		fmt.Println("Прошло время ожидания, подписка на канал 'orders' завершена")
	}

	sub.Unsubscribe()

	defer Wg.Done()
}

func parseOrdersJSON(jsonData string) (database.Orders, error) {
	var orders database.Orders
	err := json.Unmarshal([]byte(jsonData), &orders)
	if err != nil {
		return database.Orders{}, err
	}
	return orders, nil
}
