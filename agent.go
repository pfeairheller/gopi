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
	"bytes"
)


type agent struct {
	eb *Worker
	inconn redis.Conn
	outconn redis.Conn
	addr string
	in chan  []byte
	out chan []byte
}


// Create the agent of job server.
func newAgent(addr string, eb *Worker) (a *agent, err error) {
    inconn, err := redis.Dial("tcp", addr)
    if err != nil {
        return
    }
	outconn, err := redis.Dial("tcp", addr)
	if err != nil {
     return
 }

	a = &agent{
        inconn:   inconn,
        outconn:   outconn,
        eb: eb,
        addr:   addr,
        in:     make(chan []byte, 8),
        out:    make(chan []byte, 8),
    }
	return
}

func (a *agent) Close() {
	a.inconn.Close()
	a.outconn.Close()
}

func (a *agent) Work() {
	var names []interface{}
	fmt.Println(a.eb.funcs)
	for name, jobfunc := range a.eb.funcs {
		names = append(names, "gopi:queue:" + name)
		for i := 0; i < jobfunc.numberOfWorkers; i++ {
			fmt.Println("listening at ", jobfunc.c)
			go a.listen(jobfunc.c, name, jobfunc.f)
		}

	}
	names = append(names, 0)

	for {
		n, err := redis.Values(a.inconn.Do("BRPOP", names...))
		if err != nil {
			fmt.Println("Opps", err)
			continue
		}
		queuename := string(n[0].([]byte))
		name := strings.Split(queuename, ":")[2]
		a.eb.funcs[name].c <- n[1].([]byte)
	}
}

func (a *agent)listen(c chan []byte, name string, f JobFunc) {
	for data := range c {
		vals := bytes.SplitN(data, []byte("/"), 2)
		job := &Job{Fn: name, Target: string(vals[0]), Data: vals[1]}
		result, err := f(job)

		if err != nil {
			fmt.Println("error", err)
			//handle error
		}

		returnQueueName := "gopi:queue:" + string(vals[0])
		a.outconn.Do("LPUSH", returnQueueName, result)
	}
}



