package main

import (
	"bufio"
	"fmt"
	"os"
)

func cat(c chan string, e chan error, quit chan struct{}) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		c <- scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		e <- err
	}
	quit <- struct{}{}
}

func main() {
	c := make(chan string)
	quit := make(chan struct{})
	err := make(chan error)
	go cat(c, err, quit)
	for {
		select {
		case x := <-c:
			fmt.Println(x)
		case x := <-err:
			fmt.Fprintln(os.Stderr, "reading standard input:", x)
		case <-quit:
			return
		}
	}
}
