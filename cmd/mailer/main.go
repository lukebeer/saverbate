package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/emersion/go-imap/client"
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

// Make a redis pool
var redisPool = &redis.Pool{
	MaxActive: 5,
	MaxIdle:   5,
	Wait:      true,
	Dial: func() (redis.Conn, error) {
		return redis.Dial("tcp", "saverbate-redis:6379")
	},
}

type Context struct{}

func main() {
	pool := work.NewWorkerPool(Context{}, 10, "saverbate_downloads", redisPool)
	pool.PeriodicallyEnqueue("0 */5 * * * *", "check_email")

	pool.Job("check_email", (*Context).CheckEmail)

	// Start processing jobs
	pool.Start()

	// Wait for a signal to quit:
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan

	// Stop the pool
	pool.Stop()
}

func (ctx *Context) CheckEmail(job *work.Job) error {
	log.Println("Connecting to server...")

	// Connect to server
	c, err := client.DialTLS("imap.gmail.com:993", nil)
	if err != nil {
		log.Printf("Error: %v", err)
		return err
	}

	// Don't forget to logout
	defer c.Logout()

	// Login
	if err := c.Login("saverbate@gmail.com", "sUpervised896491"); err != nil {
		log.Printf("Error: %v", err)
		return err
	}
	log.Println("Logged in")

	return nil
}
