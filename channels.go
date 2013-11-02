
package main

import "fmt"

func main() {
	num := 8

	c := make(chan string)
	for i := 0; i < num; i++ {
		go listen(c, i+1)
	}


	for x := 0; x < 100000; x++ {
		data := fmt.Sprintf("Data %d", x)
		c <- data
	}

}

func listen(c chan string, num int) {
	for data := range c {
		fmt.Println(data, "at proc", num)
	}
}

