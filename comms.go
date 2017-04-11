package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/zeromq/gyre"
)

type Msg struct {
	Type string
	Payload string
	Target string
}

var (
	input = make(chan Msg)
	name  = flag.String("name", "Gyreman", "Your name or nick name in the chat session")
	users = make(map[string]string)
  Myname string
)

func Chat(mountpoint string) error {
	node, err := gyre.New()
	if err != nil {
		return err
	}
	defer node.Stop()

	err = node.Start()
	if err != nil {
		return err
	}

	node.Join("/all")

	log.Println("starting comms with mountpoint", mountpoint)
	log.Println("i am", node.Name())
	Myname = node.Name()

	myself := mountpoint + "/all/" + node.Name()
	mf, err := os.Create(myself)
	if err != nil {
    log.Panicln(err)
  }

	defer mf.Close()

	for {
		select {

		// gyre events
		case e := <-node.Events():
			switch e.Type() {
			case gyre.EventEnter:
				log.Println(e.Name()," has entered the room")
				if _, err := os.Create(mountpoint + "/all/" + e.Name()); err != nil {
			    return err
			  }
				users[e.Name()] = e.Sender()

			case gyre.EventExit:
				log.Println(e.Name()," has left the room")
				if err := os.Remove(mountpoint + "/all/" + e.Name()); err != nil {
			    return err
			  }
				delete(users, e.Name())

			case gyre.EventShout:
				fmt.Printf("%c[2K\rSHOUT from %s> %s\n", 27, *name, string(e.Msg()))

			case gyre.EventWhisper:
				fmt.Printf("%c[2K\rWHISPER from %s> %s\n", 27, *name, string(e.Msg()))
				_, err := mf.WriteString(string(e.Msg()) + "\n")
				if err != nil {
					log.Panicln(err)
				}
			}

		// fuse events
		case msg := <-input:
			switch msg.Type {
			case "DIRECT_MSG":
				log.Println(msg.Payload)
				node.Whisper(users[msg.Target], []byte(msg.Payload))

			case "JOIN_GROUP":
				log.Println("Joining group", string(msg.Payload))
				node.Join(msg.Payload)

			case "LEAVE_GROUP":
				log.Println("Leaving group", string(msg.Payload))
				node.Leave(msg.Payload)
			}
		}
	}
	return nil
} // func Chat
