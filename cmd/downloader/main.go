package main

import (
	"log"
	"time"

	"github.com/nats-io/nats.go"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	flag.String("natsAddress", "nats://localhost:10222", "Address to connect to NATS server")
	flag.Parse()

	viper.BindPFlags(flag.CommandLine)

	// NATS
	nc, err := nats.Connect(viper.GetString("natsAddress"), nats.NoEcho())
	if err != nil {
		log.Panic(err)
	}
	defer nc.Close()

	// Subscribe
	if _, err := nc.QueueSubscribe("downloading", "download", func(m *nats.Msg) {
		// Use the response
		log.Printf("Reply: %s", m.Data)

		// some long task
		time.Sleep(5 * time.Second)
	}); err != nil {
		log.Fatal(err)
	}

	select {}
}
