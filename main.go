package main

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log/slog"
	"time"
)

func main() {
	farmName := "fazenda bom dia"
	area := "sul"
	conn, msgs, salesCh := declareRabbit(farmName, area)
	defer conn.Close()

	salesCh.QueueDeclare("compradorQ", false, true, false, false, nil) //TODO remover isso
	salesCh.QueueBind("compradorQ", "comprador123", "vendas", false, nil)

	go func() {
		for m := range msgs {
			slog.Info(fmt.Sprintf("received: %s", m.Body))
			buy(salesCh, string(m.Body))
		}
	}()

	slog.Info(fmt.Sprintf("%s listening to %s", farmName, area))
	var forever chan struct{}
	<-forever
}

func buy(ch *amqp.Channel, value string) {
	slog.Info(fmt.Sprintf("buying: %s", value))
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	err := ch.PublishWithContext(ctx,
		"vendas",
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

func declareRabbit(farmName, area string) (*amqp.Connection, <-chan amqp.Delivery, *amqp.Channel) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError("initializing rabbit", err)

	ch, err := conn.Channel()
	failOnError("initializing offers channel", err)

	err = ch.ExchangeDeclare("ofertas", "topic", false, false, false, false, nil)
	failOnError("declaring offers exchange", err)

	offersQ, err := ch.QueueDeclare("", false, true, true, false, nil)
	failOnError("declaring offers queue", err)

	err = ch.QueueBind(offersQ.Name, area, "ofertas", false, nil)
	failOnError("binding offers queue", err)

	msgs, err := ch.Consume(offersQ.Name, farmName, true, false, true, false, nil)
	failOnError("declaring offers consumer", err)

	ch2, err := conn.Channel()
	failOnError("initializing channel", err)

	err = ch2.ExchangeDeclare("vendas", "direct", false, false, false, false, nil)
	failOnError("declaring sales exchange", err)

	return conn, msgs, ch2
}

func failOnError(msg string, err error) {
	if err != nil {
		slog.Error(fmt.Sprintf("error %s", msg), "err", err)
		panic(err)
	}
}
