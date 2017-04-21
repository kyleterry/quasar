package quasar

import (
	"fmt"
	"log"

	zmq "github.com/pebbe/zmq4"
)

type pubsub struct {
	ctx      *zmq.Context
	sender   *zmq.Socket
	receiver *zmq.Socket
}

type Connection struct {
	in      <-chan string
	out     chan<- string
	pubsub  *pubsub
	config  *Config
	control chan interface{}
}

func NewConnection(conf *Config) (*Connection, error) {
	conn := &Connection{
		config: conf,
		pubsub: &pubsub{},
	}
	ctx, err := zmq.NewContext()
	if err != nil {
		return nil, err
	}
	sender, err := ctx.NewSocket(zmq.PUB)
	if err != nil {
		return nil, err
	}
	receiver, err := ctx.NewSocket(zmq.SUB)
	if err != nil {
		return nil, err
	}
	receiver.SetSubscribe("")

	conn.pubsub.ctx = ctx
	conn.pubsub.sender = sender
	conn.pubsub.receiver = receiver

	return conn, nil
}

func (c *Connection) start() error {
	c.control = make(chan interface{})
	err := c.pubsub.sender.Connect(c.config.Service.SendAddr)
	if err != nil {
		return err
	}
	err = c.pubsub.receiver.Connect(c.config.Service.RecvAddr)
	if err != nil {
		return err
	}
	c.in = c.recv()
	c.out = c.send()
	return nil
}

func (c *Connection) close() {
	close(c.control)
	close(c.out)
	c.pubsub.sender.Close()
	c.pubsub.receiver.Close()
	c.pubsub.ctx.Term()
}

func (c *Connection) send() chan<- string {
	chn := make(chan string, 1000)
	go func() {
		for {
			select {
			case msg := <-chn:
				c.publish(msg)
			case <-c.control:
				return
			}
		}
	}()
	return chn
}

func (c *Connection) publish(msg string) {
	_, err := c.pubsub.sender.SendMessage(msg)
	if err != nil {
		log.Fatal(err)
	}
}

func (c *Connection) recv() <-chan string {
	chn := make(chan string, 1000)
	go func() {
		for {
			for {
				msgs, err := c.pubsub.receiver.RecvMessage(0)
				if err != nil {
					fmt.Println(err)
					continue
				}
				for _, msg := range msgs {
					chn <- msg
				}
			}
		}
	}()
	return chn
}
