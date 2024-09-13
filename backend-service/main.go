package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/streadway/amqp"
)

type VehicleEvent struct {
	ID           string `json:"id"`
	VehiclePlate string `json:"vehicle_plate"`
	DateTime     string `json:"date_time"`
	EventType    string `json:"event_type"` // "entry" or "exit"
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

	entryQueue, err := ch.QueueDeclare("vehicle_entry_queue", false, false, false, false, nil)
	exitQueue, err := ch.QueueDeclare("vehicle_exit_queue", false, false, false, false, nil)

	msgs, err := ch.Consume(entryQueue.Name, "", true, false, false, false, nil)
	exitMsgs, err := ch.Consume(exitQueue.Name, "", true, false, false, false, nil)

	go func() {
		for entry := range msgs {
			var entryEvent VehicleEvent
			json.Unmarshal(entry.Body, &entryEvent)
			entryEvent.EventType = "entry"
			processEvent(entryEvent)
		}
	}()

	go func() {
		for exit := range exitMsgs {
			var exitEvent VehicleEvent
			json.Unmarshal(exit.Body, &exitEvent)
			exitEvent.EventType = "exit"
			processEvent(exitEvent)
		}
	}()

	forever := make(chan bool)
	<-forever
}

func processEvent(event VehicleEvent) {
	url := "http://python-server:5000/log_summary"
	jsonData, _ := json.Marshal(event)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error sending to server:", err)
		return
	}
	defer resp.Body.Close()
	fmt.Println("Event processed:", event)
}
