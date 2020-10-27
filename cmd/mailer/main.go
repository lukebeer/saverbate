package main

import (
	"os"
	"os/signal"

	"saverbate/pkg/mailer"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const concurrency = 1

func main() {
	flag.String("user", "saverbate@gmail.com", "Username for imap server")
	flag.String("password", "sUpervised896491", "Password for imap server")
	flag.String("redisAddress", "saverbate-redis:6379", "Address to redis server")
	flag.String("natsAddress", "nats://saverbate-nats:4222", "Address to connect to NATS server")
	flag.Parse()

	viper.BindPFlags(flag.CommandLine)

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
	pool.PeriodicallyEnqueue("0 */15 * * * *", mailer.JobName)

	pool.Job(mailer.JobName, (*mailer.Context).CheckEmail)

	// Start processing jobs
	pool.Start()

	// Wait for a signal to quit:
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan

	// Stop the pool
	pool.Stop()
}
