package main

import (
	"github.com/tarm/goserial"
	"io"
	"log"
	"time"
)

type GSMModem struct {
	port   string
	baud   int
	status bool
	conn   io.ReadWriteCloser
}

var modem *GSMModem

func (m *GSMModem) Connect() error {
	c := &serial.Config{Name: m.port, Baud: m.baud}
	s, err := serial.OpenPort(c)
	if err == nil {
		m.status = true
		m.conn = s
	}
	return err
}

func InitModem(comport string) error {
	modem = &GSMModem{port: comport, baud: 115200}
	return modem.Connect()
}

func (m *GSMModem) sendCommand(command string) {
	time.Sleep(time.Duration(500 * time.Millisecond))
	_, err := m.conn.Write([]byte(command + "\r"))
	if err != nil {
		log.Fatal(err)
	}
}

func SendSMS(mobile string, message string) {
	// Fire and Forget
	// and hope that the SMS gets delivered

	// Put Modem in SMS Text Mode
	modem.sendCommand("AT+CMGF=1")

	modem.sendCommand("AT+CMGS=\"" + mobile + "\"")
	modem.sendCommand(message)
	// EOM CTRL-Z
	modem.sendCommand(string(26))
}
