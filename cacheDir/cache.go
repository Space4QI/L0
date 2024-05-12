package cacheDir

import (
	"L0/database"
	"encoding/json"
	"log"
)

type AllCache struct {
	orders map[string]database.Orders
}

func NewCache() *AllCache {
	return &AllCache{
		orders: make(map[string]database.Orders),
	}
}

func (c *AllCache) Read(id string) (item []byte, ok bool) {
	order, ok := c.orders[id]
	if ok {
		log.Println("from cacheDir")
		res, err := json.Marshal(order)
		if err != nil {
			log.Fatal(err)
		}
		return res, true
	}
	return nil, false
}

func (c *AllCache) updateCache(orders []database.Orders) {
	c.orders = make(map[string]database.Orders)
	for _, order := range orders {
		c.orders[order.OrderUID] = order
	}
}

func RestoreCacheFromDB(c *AllCache) {
	ordersFromDB := database.ReadFromDB()
	c.updateCache(ordersFromDB)
	log.Println("Кэш успешно восстановлен из базы данных")
}
