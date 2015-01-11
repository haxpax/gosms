package main

import (
	"fmt"
	"time"
)

const (
	SMSPending   = iota // 0
	SMSProcessed        // 1
	SMSError            // 2
)

type SMS struct {
	uuid   string
	mobile string
	body   string
	status int
}

var messages chan SMS

func InitWorker() {

	// Buffered Channel with capacity of 100 Messages
	messages = make(chan SMS, 100)

	// Start processing messages from queue
	go processMessages()
}

func EnqueueMessage(message *SMS) {
	fmt.Println("Queuing " + message.uuid)
	messages <- *message
}

func processMessages() {
	for {
		message := <-messages
		fmt.Println("processing: " + message.uuid)
		SendSMS(message.mobile, message.body)
		time.Sleep(5 * time.Second)
	}
}
