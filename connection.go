package twitch

import (
	"net"
)

type connection struct {
	isActive  tAtomBool
	Conn      *net.Conn
	Channels  []string
	reconnect chanCloser
	// read is the incoming messages channel, normally buffered with ReadBufferSize
	read chan string
	// write is the outgoing messages channel, normally buffered with WriteBufferSize
	write chan string
	// messageReceived is listened to by the pinger go-routine to interrupt the idle ping interval
	messageReceived chan bool
	disconnect      chanCloser
}

func newConnection() *connection {
	conn := &connection{
		read:            make(chan string, ReadBufferSize),
		write:           make(chan string, WriteBufferSize),
		messageReceived: make(chan bool),
	}

	conn.isActive.set(false)

	return conn
}

func (c *connection) StartParser(handleLine func(line string) error) error {
	for {
		// reader
		select {
		case msg := <-c.read:
			if err := handleLine(msg); err != nil {
				return err
			}

		case <-c.reconnect.channel:
			return errReconnect

		case <-c.disconnect.channel:
			return ErrClientDisconnected
		}
	}
}