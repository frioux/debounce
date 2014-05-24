package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"time"
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
	var timeoutSeconds = flag.Int("timeout", 1, "amount of time between output")
	flag.Parse()
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
		timeout := time.After(time.Duration(*timeoutSeconds) * time.Second)
	InnerLoop:
		for {
			select {
			case <-c:
				timeout = time.After(time.Duration(*timeoutSeconds) * time.Second)
			case x := <-err:
				fmt.Fprintln(os.Stderr, "reading standard input:", x)
			case <-quit:
				return
			case <-timeout:
				break InnerLoop
			}
		}
	}
}
