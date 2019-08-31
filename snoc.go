package snoc

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/tarm/serial"
)

type Sigfox struct {
	conf *serial.Config
	ser  *serial.Port
}

func (s *Sigfox) Init(port string) {
	portName := port

	fmt.Println("Serial port : " + portName)
	s.conf = &serial.Config{Name: portName, Baud: 9600, Parity: serial.ParityNone, StopBits: serial.Stop1, Size: 8, ReadTimeout: time.Millisecond * 500}
}

func (s *Sigfox) OpenPort() {
	serTmp, err := serial.OpenPort(s.conf)
	if err != nil {
		log.Fatal(err)
	}
	s.ser = serTmp
}

func (s *Sigfox) WaitFor(success, failure string, timeOut int) bool {
	return s.ReceiveUntil(success, failure, timeOut) != ""
}

func (s *Sigfox) ReceiveUntil(success, failure string, timeOut int) string {
	iterCount := float32(timeOut) / 0.1
	var currentMsg string
	buf := make([]byte, 128)
	for iterCount >= 0 && !strings.Contains(currentMsg, success) && !strings.Contains(currentMsg, failure) {
		time.Sleep(100 * time.Millisecond)
		n, _ := s.ser.Read(buf)
		currentMsg = fmt.Sprintf(currentMsg+"%q", buf[:n])
		iterCount -= 1
	}
	if strings.Contains(currentMsg, success) {
		return currentMsg
	} else if strings.Contains(currentMsg, failure) {
		fmt.Println("Erreur (" + strings.ReplaceAll(currentMsg, "\r\n", "") + ")")
	} else {
		fmt.Println("Délai de réception dépassé (" + strings.ReplaceAll(currentMsg, "\r\n", "") + ")")
	}
	return ""
}

func (s *Sigfox) SendMessage(message string) {
	fmt.Println("Sending SigFox Message...")

	if s.ser != nil {
		s.ser.Close()
	}

	s.OpenPort()

	s.ser.Write([]byte("AT\r"))

	if s.WaitFor("OK", "ERROR", 3) {
		fmt.Println("SigFox Modem OK")

		messSend := fmt.Sprintf("AT$SF=%s\r", message)

		s.ser.Write([]byte(messSend))
		fmt.Println("Envoi des données ...")
		if s.WaitFor("OK", "ERROR", 15) {
			fmt.Println("OK Message envoyé")
		}
	} else {
		fmt.Println("Erreur Modem SigFox")
	}

	fmt.Println("ClosePort")
	s.ser.Close()
}
