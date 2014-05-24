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
	var timeoutSeconds = flag.Float64("timeout", 1, "amount of time between output")
	var begin = flag.Bool("beginEdge", false, "trigger at start")
	var end = flag.Bool("endEdge", true, "trigger at end")
	flag.Parse()

	c := make(chan string)
	quit := make(chan struct{})
	err := make(chan error)

	go cat(c, err, quit)

	for {
		var x string
		select {
		case x = <-c:
			if *begin {
				fmt.Println(x)
			}
		case x := <-err:
			fmt.Fprintln(os.Stderr, "reading standard input:", x)
		case <-quit:
			return
		}
		timeout := time.After(time.Duration(*timeoutSeconds) * time.Second)
		shouldPrint := false
	InnerLoop:
		for {
			select {
			case x = <-c:
				shouldPrint = true
				timeout = time.After(time.Duration(*timeoutSeconds) * time.Second)
			case x := <-err:
				fmt.Fprintln(os.Stderr, "reading standard input:", x)
			case <-quit:
				return
			case <-timeout:
				if *end && shouldPrint {
					fmt.Println(x)
				}
				break InnerLoop
			}
		}
	}
}
