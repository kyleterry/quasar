package quasar

import (
	"fmt"
	"log"

	"github.com/garyburd/redigo/redis"
)

type Connection struct {
	r       redis.Conn
	in      <-chan []byte
	out     chan<- string
	pubsub  redis.PubSubConn
	service *Service
	config  Config
	control chan interface{}
}

func NewConn(conf Config, service *Service) *Connection {
	return &Connection{
		config:  conf,
		service: service,
	}
}

func (c *Connection) Start() error {
	c.control = make(chan interface{})
	r, err := c.DialRedis()
	if err != nil {
		return err
	}
	c.r = r
	c.pubsub = redis.PubSubConn{c.r}
	c.in = c.recv()
	c.out = c.send()
	return nil
}

func (c *Connection) Shutdown() {
	close(c.control)
	close(c.out)
	c.pubsub.Close()
}

func (c *Connection) DialRedis() (redis.Conn, error) {
	redisAddr := fmt.Sprintf("%s:%d", c.config.Redis.Host, c.config.Redis.Port)
	r, err := redis.Dial("tcp", redisAddr)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (c *Connection) send() chan<- string {
	chn := make(chan string, 1000)
	go func() {
		for {
			select {
			case msg := <-chn:
				c.publish(c.config.TenyksChannel, msg)
			case <-c.control:
				return
			}
		}
	}()
	return chn
}

func (c *Connection) publish(channel, msg string) {
	r, err := c.DialRedis()
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()
	r.Do("PUBLISH", channel, msg)
}

func (c *Connection) recv() <-chan []byte {
	chn := make(chan []byte, 1000)
	go func() {
		c.pubsub.Subscribe(c.config.ServiceChannel)
		for {
			for {
				switch msg := c.pubsub.Receive().(type) {
				case redis.Message:
					chn <- msg.Data
				}
			}
		}
	}()
	return chn
}
