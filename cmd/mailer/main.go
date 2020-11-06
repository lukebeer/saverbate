package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"saverbate/pkg/farm"
	"saverbate/pkg/mailer"

	"github.com/jmoiron/sqlx"
	"github.com/nats-io/nats.go"
	"github.com/robfig/cron/v3"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"

	_ "github.com/lib/pq"
)

const serviceName = "mailer"

func main() {
	log.Println("Start mailer service...")

	flag.String("dbconn", "postgres://postgres:qwerty@saverbate-db:5432/saverbate_records?sslmode=disable", "Database connection string")
	flag.String("redisAddress", "saverbate-redis:6379", "Address to redis server")
	flag.String("natsAddress", "nats://saverbate-nats:4222", "Address to connect to NATS server")
	flag.Parse()

	viper.BindPFlags(flag.CommandLine)

	log.Println("Connecting to NATS...")
	// NATS
	nc, err := nats.Connect(viper.GetString("natsAddress"), nats.NoEcho())
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	defer nc.Drain()

	db, err := sqlx.Connect("postgres", viper.GetString("dbconn"))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v\n", err)
	}

	// Check for free farm
	f, err := farm.FindFree(db, serviceName)
	if err != nil {
		log.Fatalf("Failed to retrieve free farm: %v\n", err)
	}
	viper.Set("user", f.Name)
	viper.Set("password", f.Password)

	m := mailer.New(nc)

	c := cron.New()
	c.AddFunc("*/5 * * * *", func() {
		if err := m.CheckEmail(); err != nil {
			log.Printf("ERROR: %v", err)
		}
	})
	c.Start()

	// Wait for a signal to quit:
	signalChan := make(chan os.Signal, 1)
	// SIGTERM is called when Ctrl+C was pressed
	signal.Notify(signalChan, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-signalChan

	c.Stop()

	err = farm.Release(db, serviceName, viper.GetString("user"))
	if err != nil {
		log.Printf("Failed to release farm: %v\n", err)
	}
	db.Close()
}
