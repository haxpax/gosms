package gosms

import (
	"github.com/haxpax/goserial"
	"io"
	"log"
	"runtime"
	"strings"
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
	// Setting ReadTimeout to 1secs
	var readTimeout uint32 = 10
	if runtime.GOOS == "windows" {
		readTimeout = 1000
	}
	c := &serial.Config{Name: m.port, Baud: m.baud, NonBlockingRead: true, ReadTimeout: readTimeout}
	s, err := serial.OpenPort(c)
	if err == nil {
		m.status = true
		m.conn = s
	}
	return err
}

func (m *GSMModem) SendCommand(command string, waitForOk bool) string {
	log.Println("--- SendCommand: ", command)
	var status string = ""
	_, err := m.conn.Write([]byte(command))
	if err != nil {
		log.Fatal(err)
	}
	buf := make([]byte, 32)
	var loop int = 1
	if waitForOk {
		loop = 10
	}
	for i := 0; i < loop; i++ {
		// ignoring error as EOF raises error on Linux
		n, _ := m.conn.Read(buf)
		if n > 0 {
			status = string(buf[:n])
			log.Printf("rcvd %d bytes: %s\n", n, status)
		}
	}
	return status
}

func InitModem(comport string) error {
	log.Println("--- InitModem ", comport)
	modem = &GSMModem{port: comport, baud: 115200}
	return modem.Connect()
}

func SendSMS(mobile string, message string) int {
	log.Println("--- SendSMS ", mobile, message)

	// Put Modem in SMS Text Mode
	modem.SendCommand("AT+CMGF=1\r", false)

	modem.SendCommand("AT+CMGS=\""+mobile+"\"\r", false)

	// EOM CTRL-Z = 26
	status := modem.SendCommand(message+string(26), true)
	if strings.HasSuffix(status, "OK\r\n") {
		return SMSProcessed
	} else if strings.HasSuffix(status, "ERROR\r\n") {
		return SMSError
	} else {
		return SMSPending
	}

}
