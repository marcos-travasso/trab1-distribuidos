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
	conn, _ := amqp.Dial("amqp://guest:guest@localhost:5672/")
	ch, _ := conn.Channel()
	ch.ExchangeDeclare("ofertas", "topic", false, false, false, false, nil)

	offersQ, _ := ch.QueueDeclare("", false, true, true, false, nil)
	ch.QueueBind(offersQ.Name, area, "ofertas", false, nil)
	msgs, _ := ch.Consume(offersQ.Name, farmName, true, false, true, false, nil)

	ch2, _ := conn.Channel()
	ch2.ExchangeDeclare("vendas", "direct", false, false, false, false, nil)

	return conn, msgs, ch2
}

func failOnError(msg string, err error) {
	if err != nil {
		slog.Error(fmt.Sprintf("error %s", msg), "err", err)
		panic(err)
	}
}
