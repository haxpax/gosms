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
	UUID    string `json:"uuid"`
	Mobile  string `json:"mobile"`
	Body    string `json:"body"`
	Status  int    `json:"status"`
	Retries int    `json:"retries"`
	Device  string `json:"device"`
}

var messages chan SMS
var wakeupMessageLoader chan bool

var messageCountSinceLastWakeup int
var timeOfLastWakeup time.Time

func InitWorker(modems []*GSMModem) {
	log.Println("--- InitWorker")

	numModems := len(modems)
	//TODO: check number of modems > 0

	//need to work on bufferSize and related numbers.
	//if not all modems get connected, these numbers are not quite right
	//still, good to start with
	//may be just let user decide these numbers
	bufferSize := numModems * 5

	messages = make(chan SMS, bufferSize)
	wakeupMessageLoader = make(chan bool, 1)
	wakeupMessageLoader <- true
	messageCountSinceLastWakeup = 0
	timeOfLastWakeup = time.Now().Add(-1 * time.Minute) //older time handles the cold start state of the system

	// its important to init messages channel before starting modems because nil
	// channel is non-blocking

	for i := 0; i < len(modems); i++ {
		modem := modems[i]
		if err := modem.Connect(); err == nil {
			go modem.ProcessMessages()
		}
	}
	go messageLoader(bufferSize, bufferSize/2)
}

func EnqueueMessage(message *SMS, insertToDB bool) {
	log.Println("--- EnqueueMessage: ", message)
	if insertToDB {
		insertMessage(message)
	}
	//wakeup message loader and exit
	go func() {
		//notify the message loader only if its been to too long
		//or too many messages since last notification
		messageCountSinceLastWakeup++
		if messageCountSinceLastWakeup > 10 || time.Now().Sub(timeOfLastWakeup) > time.Duration(1*time.Minute) {
			log.Println("EnqueueMessage: ", "waking up message loader")
			wakeupMessageLoader <- true
			messageCountSinceLastWakeup = 0
			timeOfLastWakeup = time.Now()
		}
		log.Println("EnqueueMessage - anon: count since last wakeup: ", messageCountSinceLastWakeup)
	}()
}

func messageLoader(bufferSize, minFill int) {
	// Load pending messages from database as needed
	for {

		/*
		   - set a fairly long timeout for wakeup
		   - if there are very few number of messages in the system and they failed at first go,
		   and there are no events happening to call EnqueueMessage, those messages might get
		   stalled in the system until someone knocks on the API door
		   - we can afford a really long polling in this case
		*/
		timeout := make(chan bool, 1)
		go func() {
			time.Sleep(10 * time.Minute)
			timeout <- true
		}()
		log.Println("messageLoader: ", "waiting for wakeup call")
		select {
		case <-wakeupMessageLoader:
			log.Println("messageLoader: woken up by channel call")
		case <-timeout:
			log.Println("messageLoader: woken up by timeout")
		}
		if len(messages) >= minFill {
			//if we have sufficient number of messages to process,
			//don't bother hitting the database
			log.Println("messageLoader: ", "I have sufficient messages")
			continue
		}

		countToFetch := bufferSize - len(messages)
		log.Println("messageLoader: ", "I need to fetch more messages", countToFetch)
		pendingMsgs, err := getPendingMessages(countToFetch)
		if err == nil {
			log.Println("messageLoader: ", len(pendingMsgs), " pending messages found")
			for _, msg := range pendingMsgs {
				messages <- msg
			}
		}
	}
}
