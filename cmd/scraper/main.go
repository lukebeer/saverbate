package main

import (
	"os"
	"os/signal"
	"saverbate/pkg/crawler"
	"syscall"

	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"github.com/robfig/cron/v3"

	flag "github.com/spf13/pflag"

	"github.com/spf13/viper"
)

func main() {
	flag.String("pydbconn", "pq://postgres:qwerty@saverbate-pg:5432/saverbate_records", "Database connection string")
	flag.String("redisAddress", "saverbate-redis:6379", "Address to redis server")
	viper.BindPFlags(flag.CommandLine)

	client := goredislib.NewClient(&goredislib.Options{
		Addr: viper.GetString("redisAddress"),
	})
	redsyncPool := goredis.NewPool(client)
	rs := redsync.New(redsyncPool)

	ctx := crawler.New(rs)

	c := cron.New()
	c.AddFunc("0 */6 * * *", func() { ctx.Crawl("cam_scrapper") })
	c.Start()

	// Wait for a signal to quit:
	signalChan := make(chan os.Signal, 1)
	// SIGTERM is called when Ctrl+C was pressed
	signal.Notify(signalChan, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-signalChan

	c.Stop()
	ctx.Close()
}
