
package client

import (
	"github.com/nu7hatch/gouuid"
	"log"
	"bytes"
	"github.com/garyburd/redigo/redis"
)

type Client struct {
	addr string
	conn redis.Conn
	in chan jobInfo
	clientPool *redis.Pool
}


type jobInfo struct {
	handle string
	result *FutureResult
}

func New(addr string) (c *Client, err error) {
	c = &Client {
		addr: addr,
		in: make(chan jobInfo),
	}

	c.conn, err = redis.Dial("tcp", c.addr)
	if err != nil {
		return
	}

	c.clientPool = &redis.Pool{
	        MaxIdle: 3,
	        MaxActive: 25, // max number of connections
	        Dial: func() (redis.Conn, error) {
	                c, err := redis.Dial("tcp", c.addr)
	                if err != nil {
	                        panic(err.Error())
	                }
	                return c, err
	        },
	}



	go c.jobLoop()

	return
}

func (c *Client) Close() {
	c.conn.Close()
	c.clientPool.Close()
}

func (c *Client) Do(jobName string, data []byte) (r *FutureResult) {
	//Create redis list var, listen for reply
	resultListName, err := uuid.NewV4()
	if err != nil {
	    log.Println("error:", err)
	    return
	}

	handle := resultListName.String()

	//Create data structure, including "listen list var"
	d := [][]byte {[]byte(handle), data}
	out := bytes.Join(d, []byte("/"))

	//Push data to redis job list var
	queueName := "gopi:queue:" + jobName
	c.conn.Do("LPUSH", queueName, out)

	r = &FutureResult {
		Success: make(chan Job, 1),
		Failure: make(chan error, 1),
	}
	c.in <- jobInfo { handle, r}

	return
}

func (c *Client) jobLoop() {
	for ji := range c.in {
		queueName := "gopi:queue:" + ji.handle

		go func() {
			clientConn := c.clientPool.Get()
			defer clientConn.Close()

			args := []interface{}{queueName, 0}
			n, err := redis.Values(clientConn.Do("BRPOP", args...))
			if err != nil {
				ji.result.Failure <- err
				return
			}

			ji.result.Success <- Job{Data: n[1].([]byte)}
		}()
	}
}


