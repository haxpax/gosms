package gosms

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
	log.Println("--- Connect")
	c := &serial.Config{Name: m.port, Baud: m.baud}
	s, err := serial.OpenPort(c)
	if err == nil {
		m.status = true
		m.conn = s
	}
	return err
}

func (m *GSMModem) SendCommand(command string) {
	log.Println("--- SendCommand: ", command)
	time.Sleep(time.Duration(500 * time.Millisecond))
	_, err := m.conn.Write([]byte(command + "\r"))
	if err != nil {
		log.Fatal(err)
	}
}

func InitModem(comport string) error {
	log.Println("--- InitModem ", comport)
	modem = &GSMModem{port: comport, baud: 115200}
	return modem.Connect()
}

func SendSMS(mobile string, message string) {
	// Fire and Forget
	// and hope that the SMS gets delivered

	log.Println("--- SendSMS ", mobile, message)

	// Put Modem in SMS Text Mode
	modem.SendCommand("AT+CMGF=1")

	modem.SendCommand("AT+CMGS=\"" + mobile + "\"")
	modem.SendCommand(message)
	// EOM CTRL-Z
	modem.SendCommand(string(26))
}
