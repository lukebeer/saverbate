package main

import (
	"log"
	"os"
	"os/signal"
	"saverbate/pkg/thumbnails"
	"syscall"

	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"

	"github.com/nats-io/nats.go"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"

	_ "github.com/lib/pq"
)

func main() {
	flag.String("natsAddress", "nats://saverbate-nats:4222", "Address to connect to NATS server")
	flag.String("redisAddress", "saverbate-redis:6379", "Address to redis server")
	flag.Parse()

	viper.BindPFlags(flag.CommandLine)

	client := goredislib.NewClient(&goredislib.Options{
		Addr: viper.GetString("redisAddress"),
	})
	redsyncPool := goredis.NewPool(client)
	rs := redsync.New(redsyncPool)

	// NATS
	nc, err := nats.Connect(viper.GetString("natsAddress"), nats.NoEcho())
	if err != nil {
		log.Panic(err)
	}

	// Run main loop of thumbnails
	t := thumbnails.New(rs, nc)
	go t.Run()

	// Wait for a signal to quit:
	signalChan := make(chan os.Signal, 1)
	// SIGTERM is called when Ctrl+C was pressed
	signal.Notify(signalChan, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-signalChan

	<-t.Close()

	if err := nc.Drain(); err != nil {
		log.Printf("ERROR: drain NATS connections failed: %v", err)
	}
}
