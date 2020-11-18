package main

import (
	"fmt"
	"log"
	"os"

	"github.com/nats-io/nats.go"
)

func main() {
	subj := os.Args[1]
	msg := os.Args[2]

	// NATS
	nc, err := nats.Connect("nats://localhost:10222", nats.NoEcho())
	if err != nil {
		log.Panic(err)
	}
	defer nc.Close()

	fmt.Println(subj)
	fmt.Println(msg)

	if err := nc.Publish(subj, []byte(msg)); err != nil {
		log.Panic(err)
	}

	nc.Flush()

	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Published [%s] : '%s'\n", subj, msg)
	}
}
