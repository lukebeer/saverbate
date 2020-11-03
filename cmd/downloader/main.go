package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"saverbate/pkg/downloader"
	"syscall"

	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"

	"github.com/jmoiron/sqlx"
	"github.com/nats-io/nats.go"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"

	_ "github.com/lib/pq"
)

type performer struct {
	Name string `json:"performer_name"`
}

func main() {
	var eventPerformer performer

	flag.String("dbconn", "postgres://postgres:qwerty@saverbate-db:5432/saverbate_records?sslmode=disable", "Database connection string")
	flag.String("natsAddress", "nats://saverbate-nats:4222", "Address to connect to NATS server")
	flag.String("redisAddress", "saverbate-redis:6379", "Address to redis server")
	flag.Parse()

	viper.BindPFlags(flag.CommandLine)

	client := goredislib.NewClient(&goredislib.Options{
		Addr: viper.GetString("redisAddress"),
	})
	redsyncPool := goredis.NewPool(client)
	rs := redsync.New(redsyncPool)
	// Initialize downloads
	dwnl := downloader.New(rs)

	db, err := sqlx.Connect("postgres", viper.GetString("dbconn"))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v\n", err)
	}

	// NATS
	nc, err := nats.Connect(viper.GetString("natsAddress"), nats.NoEcho())
	if err != nil {
		log.Panic(err)
	}

	// Subscribe
	if _, err := nc.QueueSubscribe("downloading", "download", func(m *nats.Msg) {
		if err := json.Unmarshal(m.Data, &eventPerformer); err != nil {
			log.Printf("Unmarshal error: %v", err)
			return
		}

		downloadName := eventPerformer.Name

		go dwnl.Start(downloadName)
	}); err != nil {
		log.Fatal(err)
	}

	// Wait for a signal to quit:
	signalChan := make(chan os.Signal, 1)
	// SIGTERM is called when Ctrl+C was pressed
	signal.Notify(signalChan, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-signalChan

	nc.Drain()
	db.Close()
}
