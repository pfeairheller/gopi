
package gopi

import (
	"sync"
	"log"
	"encoding/json"
	"io/ioutil"
	"github.com/garyburd/redigo/redis"
	"strings"
)

var (
	DefaultEventBus = NewEventBus()
)

func RegisterHandler(address string, handler Handler) (err error) {
	return DefaultEventBus.RegisterHandler(address, handler)
}

func Send(address string, message interface{}) (c *Call) {
	return DefaultEventBus.Send(address, message)
}

func Join(configFileName string) {
	DefaultEventBus.Join(configFileName)
}

type Call struct {
	Address string
	Message interface{}
	Reply   interface{}
	Error   error
	Done    chan *Call
}

type EventBus struct{
	mu         sync.RWMutex
	handlerMap map[string][]*handlerHolder
	inconn     redis.Conn
	outconn    redis.Conn

	Address string
}

type handlerHolder struct {
	h Handler
	c chan []byte
}

func NewEventBus() (eb *EventBus) {
	eb = new(EventBus)
	eb.handlerMap = make(map[string][]*handlerHolder)
	return
}

func (eb *EventBus) RegisterHandler(address string, handler Handler) (err error) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	if _, present := eb.handlerMap[address]; !present {
		eb.handlerMap[address] = []*handlerHolder{}
	}

	eb.handlerMap[address] = append(eb.handlerMap[address], &handlerHolder {h: handler, c: make(chan []byte, 100)})
	return
}

func (eb *EventBus) Send(address string, message interface {}) (c *Call) {
	c = new(Call)
	return
}

func (eb *EventBus) Notify(address string, message interface {}) {


}

func (eb *EventBus) Publish(address string, message interface {}) {

}

func (eb *EventBus) Join(configFileName string) (err error) {
	data, err := ioutil.ReadFile(configFileName)
	if err != nil {
		log.Fatalln("Can't open Gopi config file")
	}

	err = json.Unmarshal(data, eb)
	if err != nil {
		log.Fatalln("Unable to parse Gopi config file")
	}

	if eb.Address == "" {
		log.Fatalln("Need redis address")
	}

	eb.inconn, err = redis.Dial("tcp", eb.Address)
	if err != nil {
		log.Fatalln("Error connecting to redis", err)
	}

	eb.outconn, err = redis.Dial("tcp", eb.Address)
	if err != nil {
		log.Fatalln("Error connecting to redis", err)
	}

	var names []interface{}
	for name, holders := range eb.handlerMap {
		names = append(names, "gopi:queue:" + name)
		log.Println("listening at", name, )
		for _, holder := range holders {
			go eb.listen(holder.c, name, holder.h)
		}
	}

	names = append(names, 0)

	for {
		n, err := redis.Values(eb.inconn.Do("BRPOP", names...))
		if err != nil {
			log.Println("Opps", err)
			continue
		}
		queuename := string(n[0].([]byte))
		name := strings.Split(queuename, ":")[2]
		for _, handler := range eb.handlerMap[name] {
			handler.c <- n[1].([]byte)
		}
	}

	return
}


func (eb *EventBus) listen(c chan []byte, name string, h Handler) {
	for data := range c {
		var enc Encoding

		err := json.Unmarshal(data, &enc)
		if err != nil {
			log.Println("Unable to parse event", string(data))
			continue
		}

		msg := &Message{ Data: enc.Data }
		msg.reply = make(chan interface{}, 1)
		msg.fail = make(chan error, 1)

		go func() {
			h(msg)

			var result interface{}
			var fault error
			select {
			case result = <- msg.reply:
				data, err := json.Marshal(result)
				if err != nil {
					log.Println("Unmarshallable reply", err)
					return
				}
				log.Println("Pushing", string(data))
				eb.outconn.Do("LPUSH", enc.ReturnQueueName(), string(data))
			case fault = <- msg.fail:
				log.Println(fault)
			default:
				log.Println("No response")
			}
		}()
	}

}


