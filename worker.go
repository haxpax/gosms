package gosms

import (
	"github.com/haxpax/gosms/modem"
	"log"
	"strings"
	"time"
)

//TODO: should be configurable
const SMSRetryLimit = 3

const (
	SMSPending   = iota // 0
	SMSProcessed        // 1
	SMSError            // 2
)

type SMS struct {
	UUID      string `json:"uuid"`
	Mobile    string `json:"mobile"`
	Body      string `json:"body"`
	Status    int    `json:"status"`
	Retries   int    `json:"retries"`
	Device    string `json:"device"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

var messages chan SMS
var wakeupMessageLoader chan bool

var bufferMaxSize int
var bufferLowCount int
var messageCountSinceLastWakeup int
var timeOfLastWakeup time.Time
var messageLoaderTimeout time.Duration
var messageLoaderCountout int
var messageLoaderLongTimeout time.Duration

func InitWorker(modems []*modem.GSMModem, bufferSize, bufferLow, loaderTimeout, countOut, loaderLongTimeout int) {
	log.Println("--- InitWorker")

	bufferMaxSize = bufferSize
	bufferLowCount = bufferLow
	messageLoaderTimeout = time.Duration(loaderTimeout) * time.Minute
	messageLoaderCountout = countOut
	messageLoaderLongTimeout = time.Duration(loaderLongTimeout) * time.Minute

	messages = make(chan SMS, bufferMaxSize)
	wakeupMessageLoader = make(chan bool, 1)
	wakeupMessageLoader <- true
	messageCountSinceLastWakeup = 0
	timeOfLastWakeup = time.Now().Add((time.Duration(loaderTimeout) * -1) * time.Minute) //older time handles the cold start state of the system

	// its important to init messages channel before starting modems because nil
	// channel is non-blocking

	for i := 0; i < len(modems); i++ {
		modem := modems[i]
		err := modem.Connect()
		if err != nil {
			log.Println("InitWorker: error connecting", modem.DeviceId, err)
			continue
		}
		go processMessages(modem)
	}
	go messageLoader(bufferMaxSize, bufferLowCount)
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
		if messageCountSinceLastWakeup > messageLoaderCountout || time.Now().Sub(timeOfLastWakeup) > messageLoaderTimeout {
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
			time.Sleep(messageLoaderLongTimeout)
			timeout <- true
		}()
		log.Println("messageLoader: ", "waiting for wakeup call")
		select {
		case <-wakeupMessageLoader:
			log.Println("messageLoader: woken up by channel call")
		case <-timeout:
			log.Println("messageLoader: woken up by timeout")
		}
		if len(messages) >= bufferLowCount {
			//if we have sufficient number of messages to process,
			//don't bother hitting the database
			log.Println("messageLoader: ", "I have sufficient messages")
			continue
		}

		countToFetch := bufferMaxSize - len(messages)
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

func processMessages(modem *modem.GSMModem) {
	defer func() {
		log.Println("--- deferring ProcessMessage")
	}()

	//log.Println("--- ProcessMessage")
	for {
		message := <-messages
		log.Println("processing: ", message.UUID, modem.DeviceId)

		status := modem.SendSMS(message.Mobile, message.Body)
		if strings.HasSuffix(status, "OK\r\n") {
			message.Status = SMSProcessed
		} else if strings.HasSuffix(status, "ERROR\r\n") {
			message.Status = SMSError
		} else {
			message.Status = SMSPending
		}
		message.Device = modem.DeviceId
		message.Retries++
		updateMessageStatus(message)
		if message.Status != SMSProcessed && message.Retries < SMSRetryLimit {
			// push message back to queue until either it is sent successfully or
			// retry count is reached
			// I can't push it to channel directly. Doing so may cause the sms to be in
			// the queue twice. I don't want that
			EnqueueMessage(&message, false)
		}
		time.Sleep(5 * time.Microsecond)
	}
}
