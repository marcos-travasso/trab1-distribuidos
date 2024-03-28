package main

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log/slog"
	"time"
)

func main() {
	//slog.SetLogLoggerLevel(slog.LevelDebug)
	conn := declareRabbitDefaults()
	defer conn.Close()

	//salesCh.QueueDeclare("compradorQ", false, true, false, false, nil) //TODO remover isso
	//salesCh.QueueBind("compradorQ", "comprador123", "vendas", false, nil)

	for i := 0; i < 10; i++ {
		spawnBuyer(conn)
	}

	var forever chan struct{}
	<-forever
}

func spawnBuyer(conn *amqp.Connection) {
	farm := getRandomFarm()
	slog.Debug(fmt.Sprintf("created farm:\n%s", farm.String()))
	msgs, salesCh := declareFarmQueues(conn, farm)
	_ = salesCh

	go func() {
		for m := range msgs {
			slog.Info(fmt.Sprintf("received: %s", m.Body), "id", farm.id)
			buy(salesCh, string(m.Body))
		}
	}()

	slog.Info("farm ready", "id", farm.id)
}

func buy(ch *amqp.Channel, value string) {
	slog.Info(fmt.Sprintf("buying: %s", value))
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	err := ch.PublishWithContext(ctx,
		SalesExchange,
		"comprador123",
		false, false,
		amqp.Publishing{
			ContentType:  "text/plain",
			DeliveryMode: amqp.Transient,
			Body:         []byte(fmt.Sprintf("buying: %s", value)),
		})

	if err != nil {
		slog.Error("error publishing buy", "err", err)
	}
}

func declareFarmQueues(conn *amqp.Connection, farm Farm) (<-chan amqp.Delivery, *amqp.Channel) {
	ch, _ := conn.Channel()
	offersQ, _ := ch.QueueDeclare(farm.id, false, true, true, false, nil)

	for _, area := range farm.areas {
		ch.QueueBind(offersQ.Name, area, OffersExchange, false, nil)
		slog.Debug(fmt.Sprintf("%s listening to '%s'", farm.id, area))
	}

	msgs, _ := ch.Consume(offersQ.Name, farm.id, true, false, true, false, nil)

	return msgs, ch
}
