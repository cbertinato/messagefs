package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/zeromq/gyre"
)

var (
	input = make(chan string)
	name  = flag.String("name", "Gyreman", "Your name or nick name in the chat session")
	users = make(map[string]string)
)

func chat() {
	node, err := gyre.New()
	if err != nil {
		log.Fatalln(err)
	}
	defer node.Stop()

	err = node.Start()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("I am", node.Name())
	node.Join("/all")

	for {
		select {
		case e := <-node.Events():
			switch e.Type() {
			case gyre.EventEnter:
				users[e.Name()] = e.Sender()
				fmt.Printf("%s has entered the room", e.Name())
			case gyre.EventExit:
				delete(users, e.Name())
				fmt.Printf("%s has left the room", e.Name())
			case gyre.EventShout:
				fmt.Printf("%c[2K\rSHOUT from %s> %s%s> ", 27, e.Name(), string(e.Msg()), *name)
			case gyre.EventWhisper:
				fmt.Printf("%c[2K\rWHISPER from %s> %s%s> ", 27, e.Name(), string(e.Msg()), *name)
			}
		case msg := <-input:
			text := strings.Fields(string(msg))

			if len(text) > 0 {
				switch text[0] {
				case "SHOUT":
					payload := text[1]
					node.Shout("/all", []byte(payload))
				case "WHISPER":
					target := text[1]
					payload := strings.Join(text[2:len(text)], " ")
					node.Whisper(users[target], []byte(payload))
				default:
					fmt.Println("I don't recognize that command.")
				}
			}
		}
	}
}

func main() {
	flag.Parse()

	go chat()

	fmt.Printf("%s> ", *name)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input <- fmt.Sprintf("%s\n",scanner.Text())
		fmt.Printf("%s> ", *name)
	}
	if err := scanner.Err(); err != nil {
		log.Fatalln("reading standard input:", err)
	}
}
