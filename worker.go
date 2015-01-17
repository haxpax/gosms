package gosms

import (
	"log"
	"time"
)

const (
	SMSPending   = iota // 0
	SMSProcessed        // 1
	SMSError            // 2
)

type SMS struct {
	UUID   string `json:"uuid"`
	Mobile string `json:"mobile"`
	Body   string `json:"body"`
	Status int    `json:"status"`
}

var messages chan SMS

func InitWorker() {
	log.Println("--- InitWorker")
	// Buffered Channel with capacity of 100 Messages
	messages = make(chan SMS, 100)

	// Load pending messages from database
	pendingMsgs, err := getPendingMessages()
	if err == nil {
		log.Println(len(pendingMsgs), " pending messages found")
		for _, msg := range pendingMsgs {
			//should not use EnqueueMessage here
			//EnqueueMessage will try to insert this message again, which will cause
			//integrity error
			messages <- msg
		}
	}

	// Start processing messages from queue
	go processMessages()
}

func EnqueueMessage(message *SMS) {
	log.Println("Queuing: ", message)
	messages <- *message
	insertMessage(message)
}

func processMessages() {
	for {
		message := <-messages
		log.Println("processing: "+message.UUID, time.Now())
		message.Status = SendSMS(message.Mobile, message.Body)
		updateMessageStatus(message)
	}
}
