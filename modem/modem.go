package modem

import (
	"github.com/tarm/serial"
	"log"
	"strings"
	"time"
	"errors"
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

	if err == nil {
		m.initModem()
	}

	return err
}

func (m *GSMModem) initModem() {
	m.SendCommand("ATE0\r\n", true) // echo off
	m.SendCommand("AT+CMEE=1\r\n", true) // useful error messages
	m.SendCommand("AT+WIND=0\r\n", true) // disable notifications
	m.SendCommand("AT+CMGF=1\r\n", true) // switch to TEXT mode
}

func (m *GSMModem) Expect(possibilities []string) (string, error) {
	readMax := 0
	for _, possibility := range possibilities {
		length := len(possibility)
		if length > readMax {
			readMax = length
		}
	}

	readMax = readMax + 2; // we need offset for \r\n sent by modem

	var status string = ""
	buf := make([]byte, readMax)

	for i := 0; i < readMax; i++ {
		// ignoring error as EOF raises error on Linux
		n, _ := m.Port.Read(buf)
		if n > 0 {
			status = string(buf[:n])

			for _, possibility := range possibilities {
				if strings.HasSuffix(status, possibility) {
					log.Println("--- Expect:", m.transposeLog(strings.Join(possibilities, "|")), "Got:", m.transposeLog(status));
					return status, nil
				}
			}
		}
	}

	log.Println("--- Expect:", m.transposeLog(strings.Join(possibilities, "|")), "Got:", m.transposeLog(status), "(match not found!)");
	return status, errors.New("match not found")
}

func (m *GSMModem) Send(command string) {
	log.Println("--- Send:", m.transposeLog(command))
	m.Port.Flush()
	_, err := m.Port.Write([]byte(command))
	if err != nil {
		log.Fatal(err)
	}
}

func (m *GSMModem) Read(n int) string {
	var output string = "";
	buf := make([]byte, n)
	for i := 0; i < n; i++ {
		// ignoring error as EOF raises error on Linux
		c, _ := m.Port.Read(buf)
		if c > 0 {
			output = string(buf[:c])
		}
	}

	log.Printf("--- Read(%d): %v", n, m.transposeLog(output))
	return output
}

func (m *GSMModem) SendCommand(command string, waitForOk bool) string {
	m.Send(command)

	if waitForOk {
		output, _ := m.Expect([]string{"OK\r\n", "ERROR\r\n"}) // we will not change api so errors are ignored for now
		return output
	} else {
		return m.Read(1)
	}
}

func (m *GSMModem) SendSMS(mobile string, message string) string {
	log.Println("--- SendSMS ", mobile, message)

	m.Send("AT+CMGS=\""+mobile+"\"\r") // should return ">"
	m.Read(3)

	// EOM CTRL-Z = 26
	return m.SendCommand(message+string(26), true)
}

func (m *GSMModem) transposeLog(input string) string {
	output := strings.Replace(input, "\r\n", "\\r\\n", -1);
	return strings.Replace(output, "\r", "\\r", -1);
}