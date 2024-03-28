package main

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	OffersExchange   = "ofertas"
	SalesExchange    = "vendas"
	DeliveryExchange = "entregas"
)

func declareRabbitDefaults() *amqp.Connection {
	conn, _ := amqp.Dial("amqp://guest:guest@localhost:5672/")

	ch, _ := conn.Channel()
	ch.ExchangeDeclare(OffersExchange, "topic", false, false, false, false, nil)
	ch.ExchangeDeclare(SalesExchange, "direct", false, false, false, false, nil)
	ch.ExchangeDeclare(DeliveryExchange, "direct", false, false, false, false, nil)
	ch.Close()

	return conn
}
