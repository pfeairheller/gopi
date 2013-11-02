/**
 * Created with IntelliJ IDEA.
 * User: pfeairheller
 * Date: 11/2/13
 * Time: 9:41 AM
 * To change this template use File | Settings | File Templates.
 */
package gopi

import "github.com/garyburd/redigo/redis"


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


