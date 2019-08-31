package snoc

import (
	"errors"
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

	s.conf = &serial.Config{Name: portName, Baud: 9600, Parity: serial.ParityNone, StopBits: serial.Stop1, Size: 8, ReadTimeout: time.Millisecond * 500}
}

func (s *Sigfox) OpenPort() {
	serTmp, err := serial.OpenPort(s.conf)
	if err != nil {
		log.Fatal(err)
	}
	s.ser = serTmp
}

func (s *Sigfox) WaitFor(success, failure string, timeOut int) error {
	iterCount := float32(timeOut) / 0.1
	var currentMsg string
	var err error
	buf := make([]byte, 128)
	for iterCount >= 0 && !strings.Contains(currentMsg, success) && !strings.Contains(currentMsg, failure) {
		time.Sleep(100 * time.Millisecond)
		n, _ := s.ser.Read(buf)
		currentMsg = fmt.Sprintf(currentMsg+"%q", buf[:n])
		iterCount -= 1
	}
	if strings.Contains(currentMsg, success) {
		return nil
	} else if strings.Contains(currentMsg, failure) {
		//log.Fatal("Erreur (" + strings.ReplaceAll(currentMsg, "\r\n", "") + ")")
		err = errors.New("Error (" + strings.ReplaceAll(currentMsg, "\r\n", "") + ")")
	} else {
		err = errors.New("Timeout (" + strings.ReplaceAll(currentMsg, "\r\n", "") + ")")
	}
	return err
}

func (s *Sigfox) SendMessage(message string) error {

	var err error

	if s.ser != nil {
		s.ser.Close()
	}

	s.OpenPort()

	s.ser.Write([]byte("AT\r"))

	err = s.WaitFor("OK", "ERROR", 3)

	if err == nil {

		messSend := fmt.Sprintf("AT$SF=%s\r", message)

		s.ser.Write([]byte(messSend))
		err = s.WaitFor("OK", "ERROR", 15)
	}

	s.ser.Close()
	return err
}
