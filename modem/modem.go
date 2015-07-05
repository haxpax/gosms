package modem

import (
	"github.com/tarm/serial"
	"log"
	"strings"
	"time"
)

type GSMModem struct {
	ComPort  string
	BaudRate int
	Port     *serial.Port
	DeviceId string
}

func New(ComPort string, BaudRate int, DeviceId string) (modem *GSMModem) {
	modem = &GSMModem{ComPort: ComPort, BaudRate: BaudRate, DeviceId: DeviceId}
	return modem
}

func (m *GSMModem) Connect() (err error) {
	config := &serial.Config{Name: m.ComPort, Baud: m.BaudRate, ReadTimeout: time.Second}
	m.Port, err = serial.OpenPort(config)
	return err
}

func (m *GSMModem) SendCommand(command string, waitForOk bool) string {
	log.Println("--- SendCommand: ", command)
	var status string = ""
	m.Port.Flush()
	_, err := m.Port.Write([]byte(command))
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
		n, _ := m.Port.Read(buf)
		if n > 0 {
			status = string(buf[:n])
			log.Printf("SendCommand: rcvd %d bytes: %s\n", n, status)
			if strings.HasSuffix(status, "OK\r\n") || strings.HasSuffix(status, "ERROR\r\n") {
				break
			}
		}
	}
	return status
}

func (m *GSMModem) SendSMS(mobile string, message string) string {
	log.Println("--- SendSMS ", mobile, message)

	// Put Modem in SMS Text Mode
	m.SendCommand("AT+CMGF=1\r", false)

	m.SendCommand("AT+CMGS=\""+mobile+"\"\r", false)

	// EOM CTRL-Z = 26
	return m.SendCommand(message+string(26), true)

}
