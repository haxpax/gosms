package gosms

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
	UUID   string `json:"uuid"`
	Mobile string `json:"mobile"`
	Body   string `json:"body"`
	Status int    `json:"status"`
}

var messages chan SMS

func InitWorker() {

	// Buffered Channel with capacity of 100 Messages
	messages = make(chan SMS, 100)

	// Load pending messages from database
	pendingMsgs, err := getPendingMessages()
	if err == nil {
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
	fmt.Println("Queuing " + message.UUID)
	messages <- *message
	insertMessage(message)
}

func processMessages() {
	for {
		message := <-messages
		fmt.Println("processing: " + message.UUID)
		SendSMS(message.Mobile, message.Body)
		/*
		   TODO: modify this block to check result of SendSMS and change status
		   code accordingly
		*/
		message.Status = SMSProcessed
		updateMessageStatus(message)

		time.Sleep(5 * time.Second)
	}
}
