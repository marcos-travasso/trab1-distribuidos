package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log/slog"
	"time"
	"math/rand"
)

func main() {
	//slog.SetLogLoggerLevel(slog.LevelDebug)
	conn := declareRabbitDefaults()
	defer conn.Close()

	// salesCh.QueueDeclare("compradorQ", false, true, false, false, nil) //TODO remover isso
	// salesCh.QueueBind("compradorQ", "comprador123", "vendas", false, nil)

	for i := 0; i < 10; i++ {
		spawnBuyer(conn)
	}

	var forever chan struct{}
	<-forever
}

func spawnBuyer(conn *amqp.Connection) {
	farm := getRandomFarm()
	slog.Debug(fmt.Sprintf("created farm:\n%s", farm.String()))
	farm.declareQueues(conn)

	go func() {
		for m := range farm.offersCh {
			time.Sleep(time.Duration(rand.Intn(5000) + 500) * time.Millisecond)
			slog.Info(fmt.Sprintf("offer message received: %s", m.Body), "id", farm.id)
			if checkOffer(m.Body) {
				farm.buy(m.Body)
			}
		}
	}()

	go func() {
		for m := range farm.deliveriesCh {
			slog.Info(fmt.Sprintf("delivery message received: %s", m.Body), "id", farm.id)
			handleDelivery(m.Body)
		}
	}()

	slog.Info("farm ready", "id", farm.id)
}
