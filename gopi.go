package gopi


import (
	"github.com/garyburd/redigo/redis"
	"fmt"
)


func ListenAndWork(numberWorkers int, worker *Worker) {

	for name, f := range worker.funcs{
		c, err := redis.Dial("tcp", worker.)
		if err != nil {
			fmt.Println("Opps", err)
			panic(err)
			// handle error
		}
		defer c.Close()

		c := make(chan []byte)
		for i := 0; i < numberWorkers; i++ {
			go listen(c, f)
		}

		for {
			n, err := redis.Values(c.Do("BLPOP", name, 0))
			if err != nil {
				fmt.Println("Opps", err)
				continue
			}
			c <- n[1].([]byte)
		}
	}

}

func listen(c chan []byte, f JobFunc) {
	for data := range c {
		job := &Job {data}
		f(job)
	}
}

