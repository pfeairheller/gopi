
package gopi

import "fmt"

func main() {
	num := 8

	c := make(chan []byte)
	for i := 0; i < num; i++ {
		go listen(c, nil)
	}


	for x := 0; x < 100000; x++ {
		data := fmt.Sprintf("Data %d", x)
		c <- []byte(data)
	}

}

