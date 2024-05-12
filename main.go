package main

import (
	"L0/cacheDir"
	"L0/database"
	"L0/handlers"
	"L0/pubsub"
	"context"
	"database/sql"
	"fmt"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
)

const (
	Host     = "localhost"
	Port     = 5432
	User     = "postgres"
	Password = "itsme"
	Dbname   = "mypostgres"
)

func main() {
	clusterId := "test-cluster"
	clientId := "my-client"

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", Host, Port, User, Password, Dbname)
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных:", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Ошибка проверки подключения к базе данных:", err)
	}
	fmt.Println("Подключение к базе данных успешно установлено")

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal("Ошибка подключения к NATS:", err)
	}
	defer nc.Close()
	fmt.Println("Подключение к NATS успешно установлено")

	sc, err := stan.Connect(clusterId, clientId)
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Println("Ошибка начала транзакции", err)
		return
	}
	defer tx.Rollback()

	//Выполнение скрипта обновления базы данных
	if err := database.UpdateDatabase(db); err != nil {
		log.Fatal("Ошибка обновления базы данных: ", err)
	}

	fmt.Println("База данных успешно обновлена")

	sqlScript, err := ioutil.ReadFile("createDB.sql")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(string(sqlScript))
	if err != nil {
		log.Fatal(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Println("Ошибка при совершении транзакции", err)
		return
	}

	pubsub.Publish(clientId, clusterId, database.Orders{})

	ctx, cancelFunc := context.WithCancel(context.Background())

	pubsub.Wg.Add(1)

	go pubsub.Subscribe(ctx, db, nc)

	pubsub.Wg.Wait()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		cancelFunc()
		os.Exit(0)
	}()

	// Создание кэша
	cache := cacheDir.NewCache()

	// Восстановление кэша из базы данных
	cacheDir.RestoreCacheFromDB(cache)

	router := httprouter.New()
	router.GET("/order/:order_uid", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		handlers.GetOrder(w, p, cache)
	})

	// Блокирующий вызов ListenAndServe
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
	}

}
