package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"saverbate/pkg/farm"
	"saverbate/pkg/mailer"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	"github.com/jmoiron/sqlx"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"

	_ "github.com/lib/pq"
)

const (
	concurrency = 1
	serviceName = "mailer"
)

// TODO: Do not automatically retry if got any error
func main() {
	flag.String("dbconn", "postgres://postgres:qwerty@saverbate-db:5432/saverbate_records?sslmode=disable", "Database connection string")
	flag.String("redisAddress", "saverbate-redis:6379", "Address to redis server")
	flag.String("natsAddress", "nats://saverbate-nats:4222", "Address to connect to NATS server")
	flag.Parse()

	viper.BindPFlags(flag.CommandLine)

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

	// Make a redis pool
	redisPool := &redis.Pool{
		MaxActive: 5,
		MaxIdle:   5,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", viper.GetString("redisAddress"))
		},
	}

	pool := work.NewWorkerPool(mailer.Context{}, concurrency, mailer.NamespaceDownloads, redisPool)
	pool.PeriodicallyEnqueue("0 */5 * * * *", mailer.JobName)

	opts := work.JobOptions{
		SkipDead: true,
	}
	pool.JobWithOptions(mailer.JobName, opts, (*mailer.Context).CheckEmail)

	// Start processing jobs
	pool.Start()

	// Wait for a signal to quit:
	signalChan := make(chan os.Signal, 1)
	// SIGTERM is called when Ctrl+C was pressed
	signal.Notify(signalChan, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-signalChan

	// Stop the pool
	pool.Stop()

	err = farm.Release(db, serviceName, viper.GetString("user"))
	if err != nil {
		log.Fatalf("Failed to release farm: %v\n", err)
	}
	db.Close()
}
