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
	var timeoutSeconds float64
	var begin, end, h, help bool

	flag.Float64Var(&timeoutSeconds, "timeout", 1, "amount of time between output")
	flag.BoolVar(&begin, "beginEdge", false, "trigger at start")
	flag.BoolVar(&end, "endEdge", true, "trigger at end")
	flag.BoolVar(&h, "h", false, "help for debounce")
	flag.BoolVar(&help, "help", false, "help for debounce")

	flag.Parse()

	if h || help {
		fmt.Println("\n" +
			" debounce          [--beginEdge] [--endEdge] [--timeout 2] [-h|--help]\n" +
			"    --beginEdge     pass this flag to output at the start of a cycle\n" +
			"    --endEdge       pass this flag to output after the end of a cycle\n" +
			"    --timeout       set the timeout in seconds, default is 1 second\n" +
			"\n" +
			"    -h --help       print usage message and exit\n",
		)
		return
	}

	c := make(chan string)
	quit := make(chan struct{})
	err := make(chan error)

	go cat(c, err, quit)

	for {
		var x string
		shouldPrint := false
		select {
		case x = <-c:
			shouldPrint = true
			if begin {
				shouldPrint = false
				fmt.Println(x)
			}
		case x := <-err:
			fmt.Fprintln(os.Stderr, "reading standard input:", x)
		case <-quit:
			return
		}
		timeout := time.After(time.Duration(timeoutSeconds) * time.Second)
	InnerLoop:
		for {
			select {
			case x = <-c:
				shouldPrint = true
				timeout = time.After(time.Duration(timeoutSeconds) * time.Second)
			case x := <-err:
				fmt.Fprintln(os.Stderr, "reading standard input:", x)
			case <-quit:
				return
			case <-timeout:
				if end && shouldPrint {
					fmt.Println(x)
				}
				break InnerLoop
			}
		}
	}
}
