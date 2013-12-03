package gopi

import (
	"encoding/json"
)

type Message struct {
	Data  []byte
	reply chan interface{}
	fail  chan error
}

type Encoding struct {
	ReturnQueue string
	Data []byte
}

func (e *Encoding) ReturnQueueName() string {
	return "gopi:queue:" + e.ReturnQueue
}

func (m *Message) Reply(reply interface{}) {
	m.reply <- reply
}

func (m *Message) Fail(err error) {
	m.fail <- err
}


func (m *Message) Body(target interface{}) (err error) {
	return json.Unmarshal(m.Data, target)
}
