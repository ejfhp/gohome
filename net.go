package gohome

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/pkg/errors"
)

var ErrNoConnection = errors.New("NO CONNECTION")
var ErrNoData = errors.New("NO DATA")
var ErrNAK = errors.New("NAK")
var ErrServerNotFound = errors.New("SERVER NOT FOUND")
var ErrConnectionFailed = errors.New("CONNECTION FAILED")

//HomeError wraps OWN errors
type HomeError struct {
	code     string
	incoming string
	answer   string
}

// SystemMessages contains the OpenWebNet codes for various system messages
var SystemMessages = map[string]Message{
	"ACK":                   "*#*1##",
	"NACK":                  "*#*0##",
	"OPEN_COMMAND_SESSION":  "*99*0##", // OpenWebNet command to ask for a command session
	"OPEN_EVENT_SESSION":    "*99*1##",
	"OPEN_SCENARIO_SESSION": "*99*9##",
}

type Cable struct {
	address string
}

//Home is a Btcino MyHome plant that can be controlled with a OpenWebNet enabled device (F452 ecc)
type Home struct {
	cable *Cable
	plant *Plant
}

//NewHome creates a new Home connected through the given Cable
func NewHome(plant *Plant) *Home {
	log.Printf("NewHome")
	address := plant.ServerAddress()
	cable := newCable(address)
	return &Home{cable: cable, plant: plant}
}

//Do some action with your home
func (h *Home) Do(command Message) error {
	log.Printf("Home.Do")
	return h.cable.sendCommand(command)
}

//Ask the system
func (h *Home) Ask(request Message) ([]Message, error) {
	log.Printf("Home.Ask")
	return h.cable.sendRequest(request)
}

func (h *Home) Listen() (<-chan Message, chan<- struct{}, <-chan error) {
	msgChan := make(chan Message, 1)
	signChan := make(chan struct{})
	errChan := make(chan error)
	go h.cable.listen(msgChan, signChan, errChan)
	return msgChan, signChan, errChan
}

func newCable(address string) *Cable {
	c := Cable{address: address}
	return &c
}

func (c *Cable) connect() (*net.TCPConn, error) {
	log.Printf("Cable.connect address:%s", c.address)
	timeout := 1 * time.Second
	conn, err := net.DialTimeout("tcp", c.address, timeout)
	if conn == nil {
		return nil, errors.Wrapf(ErrNoConnection, "no dial to server address: %s", c.address)
	}
	if err != nil {
		return nil, errors.Wrapf(err, "no dial to server address: %s", c.address)
	}
	connTCP := conn.(*net.TCPConn)
	if err != nil {
		return nil, errors.Wrapf(ErrNoConnection, "no dial to server address: %s", c.address)
	}
	if err = c.acked(connTCP); err != nil {
		return nil, errors.Wrapf(ErrConnectionFailed, "NAK from server address: %s", c.address)
	}
	return connTCP, nil
}

func (c *Cable) sendCommand(message Message) error {
	log.Printf("Cable.SendCommmand message:%s", message)
	conn, err := c.connect()
	if err != nil {
		return errors.Wrap(err, "cannot connect")
	}
	defer conn.Close()
	log.Printf("Cable.open_command_session")
	if c.send(conn, SystemMessages["OPEN_COMMAND_SESSION"]) != nil {
		return errors.Wrap(err, "cannot open session")
	}
	log.Printf("Cable.check ack")
	if err = c.acked(conn); err != nil {
		return errors.Wrapf(ErrConnectionFailed, "NAK from server address: %s", c.address)
	}
	log.Printf("Sending nessage")
	if c.send(conn, message) != nil {
		return errors.Wrapf(err, "cannot send message %s, ", message)
	}
	log.Printf("Cable.check ack")
	if err = c.acked(conn); err != nil {
		return errors.Wrapf(ErrConnectionFailed, "NAK from server address: %s", c.address)
	}
	return nil
}

func (c *Cable) sendRequest(message Message) ([]Message, error) {
	log.Printf("Cable.SendRequest request:%s", message)
	conn, err := c.connect()
	if err != nil {
		return nil, errors.Wrap(err, "cannot connect")
	}
	defer conn.Close()
	if c.send(conn, SystemMessages["OPEN_COMMAND_SESSION"]) != nil {
		return nil, errors.Wrapf(err, "cannot open command session")
	}
	if err = c.acked(conn); err != nil {
		return nil, errors.Wrapf(ErrConnectionFailed, "NAK from server address: %s", c.address)
	}
	if c.send(conn, message) != nil {
		return nil, errors.Wrapf(err, "cannot send request %s, ", message)
	}
	answers := make([]Message, 0, 10)
	for {
		a, err := c.receive(conn, false)
		if err != nil {
			return answers, errors.Wrapf(err, "failed to receive answer for request: %s", message)
		}
		if a == SystemMessages["ACK"] {
			break
		}
		answers = append(answers, a)
	}
	return answers, nil
}

func (c *Cable) listen(out chan<- Message, in <-chan struct{}, errs chan<- error) {
	log.Printf("Cable.Listen")
	conn, err := c.connect()
	if err != nil {
		errs <- errors.Wrap(err, "cannot connect")
		close(out)
		log.Printf("Cable.Listen channel closed, return")
		return
	}
	defer conn.Close()
	if c.send(conn, SystemMessages["OPEN_EVENT_SESSION"]) != nil {
		errs <- errors.Wrapf(err, "cannot open event session")
		close(out)
		close(errs)
		return
	}
	if err = c.acked(conn); err != nil {
		errs <- errors.Wrapf(ErrConnectionFailed, "NAK from server address: %s", c.address)
		close(out)
		close(errs)
		return
	}
	for {
		select {
		case <-in:
			fmt.Printf("Closing.. \n")
			break
		default:
			fmt.Printf("Receiving.. \n")
			a, err := c.receive(conn, true)
			fmt.Printf("Done receiving.\n")
			if err != nil {
				errs <- errors.Wrapf(err, "failed to receive events")
			}
			msg := Message(a)
			out <- msg
		}
	}
}

func (c *Cable) send(conn *net.TCPConn, message Message) error {
	log.Printf("Cable.send message:%s", message)
	_, err := conn.Write([]byte(message))
	if err != nil {
		return errors.Wrap(ErrConnectionFailed, "failed to send")
	}
	return nil
}

func (c *Cable) acked(conn *net.TCPConn) error {
	log.Printf("Cable.acked")
	msg, err := c.receive(conn, false)
	if err != nil {
		return errors.Wrap(err, "cannot check ACK")
	}
	return isAck(msg)
}

//returns answer, ok
func (c *Cable) receive(conn *net.TCPConn, noTimeout bool) (Message, error) {
	log.Printf("Cable.receive")
	if conn == nil {
		return "", errors.Wrap(ErrNoConnection, "cannot receive from nil connection")
	}
	msg := make([]byte, 0, 20)
	b := make([]byte, 1)
	var end bool
	ok := true
	for !end {
		conn.SetReadDeadline(time.Now().Add(time.Second * 1))
		n, err := conn.Read(b)
		if err, ok := err.(net.Error); ok && err.Timeout() && noTimeout {
			continue
		}
		if err != nil {
			return "", errors.Wrap(err, "cannot read from connection")
		}
		if n == 0 {
			return "", ErrNoData
		}
		msg = append(msg, b[0])
		//TODO sostituire con regexp
		if len(msg) > 1 && msg[len(msg)-1] == '#' && msg[len(msg)-2] == '#' {
			end = true
		}
	}
	ans := Message(msg)
	log.Printf("Cable.receive: msg:%s, ok:%t", ans, ok)
	return ans, nil
}

func isAck(m Message) error {
	log.Printf("isAck")
	if m != SystemMessages["ACK"] {
		return ErrNAK
	}
	return nil
}
