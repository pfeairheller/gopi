package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/pfeairheller/gopi"
)

type Student struct {
	FirstName string
	LastName  string
}

func main() {
	waitForGoopie()
}


func waitForGoopie() {
	fmt.Println("Starting")
	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		fmt.Println("Opps", err)
		panic(err)
		// handle error
	}
	defer c.Close()

	i:= 0
	for ;;i++{
		n, err := redis.Values(c.Do("BLPOP", "queue:name", 0))
		if err != nil {
			fmt.Println("Opps", err)
		}

		job := &gopi.Job {n[1].([]byte)}
		workerHandler(job)
		fmt.Println(i)
	}

}

func workerHandler(job *gopi.Job) error {
	var student Student
	err := job.Value(&student)
	if err != nil {
		fmt.Println("error extracting", err)
		return err
	}

	fmt.Println(student)
	return nil;
}

