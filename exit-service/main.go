package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/streadway/amqp"
)

type ExitEvent struct {
	ID           string `json:"id"`
	VehiclePlate string `json:"vehicle_plate"`
	ExitDateTime string `json:"exit_date_time"`
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"vehicle_exit_queue", // queue name
		false,                // durable
		false,                // delete when unused
		false,                // exclusive
		false,                // no-wait
		nil,                  // arguments
	)
	if err != nil {
		panic(err)
	}

	for {
		event := ExitEvent{
			ID:           fmt.Sprintf("%d", rand.Int()),
			VehiclePlate: fmt.Sprintf("ABC%03d", rand.Intn(1000)),
			ExitDateTime: time.Now().UTC().Format(time.RFC3339),
		}

		body, err := json.Marshal(event)
		if err != nil {
			panic(err)
		}

		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        body,
			})
		if err != nil {
			panic(err)
		}

		fmt.Printf("Generated vehicle exit: %+v\n", event)
		time.Sleep(3 * time.Second) // simulate event generation delay
	}
}
