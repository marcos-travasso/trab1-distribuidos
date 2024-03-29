package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log/slog"
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
			slog.Info(fmt.Sprintf("received: %s", m.Body), "id", farm.id)
			if checkOffer(m.Body) {
				farm.buy(m.Body)
			}
		}
	}()

	slog.Info("farm ready", "id", farm.id)
}
