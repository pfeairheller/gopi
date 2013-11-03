package gopi


import (
	"github.com/garyburd/redigo/redis"
	"fmt"
)


func ListenAndWork(numberWorkers int, worker *Worker) {


}

func listen(c chan []byte, f JobFunc) {
	for data := range c {
		job := &Job {data}
		f(job)
	}
}

