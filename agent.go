/**
 * Created with IntelliJ IDEA.
 * User: pfeairheller
 * Date: 11/2/13
 * Time: 9:41 AM
 * To change this template use File | Settings | File Templates.
 */
package gopi

import (
	"github.com/garyburd/redigo/redis"
	"fmt"
	"strings"
)


type agent struct {
	worker *Worker
	conn redis.Conn
	addr string
	in chan []byte
	out chan *Job
}

// Create the agent of job server.
func newAgent(addr string, worker *Worker) (a *agent, err error) {
    conn, err := redis.Dial("tcp", addr)
    if err != nil {
        return
    }

	a = &agent{
        conn:   conn,
        worker: worker,
        addr:   addr,
        in:     make(chan []byte, 8),
        out:    make(chan *Job, 8),
    }
	return
}

func (a *agent) Close() {
	a.conn.Close()
}

func (a *agent) Work() {
	c, err := redis.Dial("tcp", a.addr)
	if err != nil {
		fmt.Println("Opps", err)
		panic(err)
		// handle error
	}
	defer c.Close()

	var names []interface{}
	fmt.Println(a.worker.funcs)
	for name, jobfunc := range a.worker.funcs {
		names = append(names, "gopi:queue:" + name)
		for i := 0; i < jobfunc.numberOfWorkers; i++ {
			fmt.Println("listening at ", jobfunc.c)
			go listen(jobfunc.c, name, jobfunc.f)
		}

	}
	names = append(names, 0)

	for {
		n, err := redis.Values(c.Do("BLPOP", names...))
		if err != nil {
			fmt.Println("Opps", err)
			continue
		}
		queuename := string(n[0].([]byte))
		name := strings.Split(queuename, ":")[2]
		a.worker.funcs[name].c <- n[1].([]byte)
	}
}

func listen(c chan []byte, name string, f JobFunc) {
	for data := range c {
		job := &Job{Fn: name, Data: data}
		f(job)
	}
}



