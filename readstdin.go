package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	fmt.Println("Agent Communicator")
	fmt.Println("---------------------")

	inputChan := make(chan string, 100)
	quitChan := make(chan bool)

	go readInput(inputChan, quitChan)
	go agent(inputChan)

	<-quitChan
}

func readInput(output chan string, quit chan bool) {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("send a message to the agent-> ")
		text, err := reader.ReadString('\n')

		if err != nil {
			panic(err)
		}

		// CLRF to LF stupidness
		text = strings.Replace(text, "\n", "", -1)

		if shouldQuit(text) {
			fmt.Print("Exiting")
			quit <- true
			close(output)
			return
		}

		output <- text
	}
}

func shouldQuit(msg string) bool {
	return msg == "quit" || msg == "exit"
}

func agent(input chan string) {
	for {
		select {
		case s := <-input:
			send(s)
		default:
		}
	}
}

func send(msg string) {
	time.Sleep(time.Duration(3) * time.Second)
	fmt.Println("\n\nreceived msg", msg)
}
