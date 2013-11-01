package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"encoding/json"
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

		var student Student
		err = json.Unmarshal(n[1].([]byte), &student)
		if err != nil {
			fmt.Println("error extracting", err)
		}
		fmt.Println(i, student.FirstName)
	}

}

