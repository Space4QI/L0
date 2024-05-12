package handlers

import (
	"L0/cacheDir"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func GetOrder(w http.ResponseWriter, p httprouter.Params, cache *cacheDir.AllCache) {
	orderUID := p.ByName("order_uid")

	// Проверяем наличие данных в кэше
	orderBytes, ok := cache.Read(orderUID)
	if ok {
		// Данные найдены в кэше, отправляем их в ответ
		w.Header().Set("Content-Type", "application/json")
		w.Write(orderBytes)
		return
	}

	fmt.Fprintf(w, "Данные для заказа с ID %s не найдены в кэше", orderUID)
}
