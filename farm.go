package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Farm struct {
	id    string
	name  string
	areas []string
	offersCh <-chan amqp.Delivery
	salesCh *amqp.Channel
}

func (f Farm) String() string {
	return fmt.Sprintf("--------------------------------------------\n"+
		"ID: %s\n"+
		"Name: %s\n"+
		"Areas: %+v\n"+
		"--------------------------------------------\n", f.id, f.name, f.areas)
}

func (f Farm) buy(offerPayload []byte) {
	offer := make(map[string]interface{})
	json.Unmarshal(offerPayload, &offer)

	slog.Info(fmt.Sprintf("buying: %s", offer["id"].(string)), "id", f.id)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	buyOrder := map[string]string{
		"id_gado": offer["id"].(string),
		"id_comprador": f.id,
	}
	buyOrderPayload, _ := json.Marshal(buyOrder)

	f.salesCh.PublishWithContext(ctx,
		SalesExchange,
		offer["id_vendedor"].(string),
		false, false,
		amqp.Publishing{
			ContentType:  "text/plain",
			DeliveryMode: amqp.Transient,
			Body:         buyOrderPayload,
	})
	slog.Info(fmt.Sprintf("order sent to %s", offer["id_vendedor"].(string)), "id", f.id)
}

func (f *Farm) declareQueues(conn *amqp.Connection) {
	f.salesCh, _ = conn.Channel()
	offersQ, _ := f.salesCh.QueueDeclare(f.id, false, true, true, false, nil)

	for _, area := range f.areas {
		f.salesCh.QueueBind(offersQ.Name, area, OffersExchange, false, nil)
		slog.Debug(fmt.Sprintf("%s listening to '%s'", f.id, area))
	}

	f.offersCh, _ = f.salesCh.Consume(offersQ.Name, f.id, true, false, true, false, nil)
}