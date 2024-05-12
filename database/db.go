package database

import (
	"database/sql"
	"fmt"
	"log"
)

func UpdateDatabase(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("ошибка начала транзакции: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		DELETE FROM orders;
		DELETE FROM delivery;
		DELETE FROM payment;
		DELETE FROM items;
	`)
	if err != nil {
		return fmt.Errorf("ошибка выполнения скрипта обновления базы данных: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("ошибка фиксации транзакции: %w", err)
	}

	return nil
}

func ReadFromDB() []Orders {
	psqlInfo := fmt.Sprintf("host=localhost port=5432 user=postgres password=itsme dbname=mypostgres sslmode=disable")
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных:", err)
	}
	defer db.Close()

	request, err := db.Query(`
        SELECT 
            o.order_uid, o.track_number, o.entry, o.locale, o.internal_signature,
            o.customer_id, o.delivery_service, o.shardkey, o.sm_id, o.date_created, o.oof_shard,
            d.name, d.phone, d.zip, d.city, d.address, d.region, d.email,
            p.transaction, p.request_id, p.currency, p.provider, p.amount, p.payment_dt, p.bank, p.delivery_cost, p.goods_total, p.custom_fee,
            i.chrt_id, i.track_number, i.price, i.rid, i.name, i.sale, i.size, i.total_price, i.nm_id, i.brand, i.status
        FROM 
            orders o
        JOIN 
            delivery d ON o.order_uid = d.order_uid
        JOIN 
            payment p ON o.order_uid = p.order_uid
        JOIN 
            items i ON o.order_uid = i.order_uid
    `)
	if err != nil {
		log.Fatal("Ошибка выполнения запроса к базе данных:", err)
	}
	defer request.Close()

	var orders []Orders
	for request.Next() {
		var order Orders
		var item Item

		err := request.Scan(
			&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature,
			&order.CustomerID, &order.DeliveryService, &order.Shardkey, &order.SmID, &order.DateCreated, &order.OofShard,
			&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip, &order.Delivery.City, &order.Delivery.Address,
			&order.Delivery.Region, &order.Delivery.Email,
			&order.Payment.Transaction, &order.Payment.RequestID, &order.Payment.Currency,
			&order.Payment.Provider, &order.Payment.Amount, &order.Payment.PaymentDt, &order.Payment.Bank, &order.Payment.DeliveryCost,
			&order.Payment.GoodsTotal, &order.Payment.CustomFee,

			&item.ChrtID, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale, &item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status,
		)
		if err != nil {
			log.Fatal("Ошибка сканирования поля массива Items:", err)
		}
		order.Items = append(order.Items, item)

		orders = append(orders, order)

	}
	if err := request.Err(); err != nil {
		log.Fatal("Ошибка получения следующей строки результата:", err)
	}

	return orders
}

func SaveToDB(db *sql.DB, orders Orders) error {
	// Первый INSERT - orders
	_, err := db.Exec(`
	INSERT INTO orders (
		order_uid, track_number, entry, locale, internal_signature,
		customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
	)
`, orders.OrderUID, orders.TrackNumber, orders.Entry, orders.Locale,
		orders.InternalSignature, orders.CustomerID, orders.DeliveryService,
		orders.Shardkey, orders.SmID, orders.DateCreated, orders.OofShard)
	if err != nil {
		log.Println("Ошибка при выполнении INSERT в таблицу orders:", err)
		return err
	}

	// Второй INSERT - delivery
	_, err = db.Exec(`
	INSERT INTO delivery (
		order_uid, name, phone, zip, city, address, region, email
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8
	)`, orders.OrderUID, orders.Delivery.Name, orders.Delivery.Phone, orders.Delivery.Zip,
		orders.Delivery.City, orders.Delivery.Address, orders.Delivery.Region,
		orders.Delivery.Email)
	if err != nil {
		log.Println("Ошибка при выполнении первого INSERT в таблицу delivery:", err)
		return err
	}

	// Третий INSERT - payment
	_, err = db.Exec(`
	INSERT INTO payment (
		order_uid, transaction, request_id, currency, provider, amount, payment_dt, bank,
		delivery_cost, goods_total, custom_fee
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
	)
`, orders.OrderUID, orders.Payment.Transaction, orders.Payment.RequestID,
		orders.Payment.Currency, orders.Payment.Provider, orders.Payment.Amount,
		orders.Payment.PaymentDt, orders.Payment.Bank, orders.Payment.DeliveryCost,
		orders.Payment.GoodsTotal, orders.Payment.CustomFee)
	if err != nil {
		log.Println("Ошибка при выполнении INSERT в таблицу payment:", err)
		return err
	}

	// Четвертый INSERT - items
	_, err = db.Exec(`
	INSERT INTO items (
		order_uid, chrt_id, track_number, price, rid, name, sale, size,
		total_price, nm_id, brand, status
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
	)
`, orders.OrderUID, orders.Items[0].ChrtID, orders.Items[0].TrackNumber, orders.Items[0].Price,
		orders.Items[0].Rid, orders.Items[0].Name, orders.Items[0].Sale,
		orders.Items[0].Size, orders.Items[0].TotalPrice, orders.Items[0].NmID,
		orders.Items[0].Brand, orders.Items[0].Status)
	if err != nil {
		log.Println("Ошибка при выполнении INSERT в таблицу items:", err)
		return err
	}
	return nil
}
