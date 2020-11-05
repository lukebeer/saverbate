package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"saverbate/pkg/farm"
	"syscall"
	"time"

	flag "github.com/spf13/pflag"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"

	_ "github.com/lib/pq"
)

const (
	concurrency = 1
	serviceName = "subscriber"
)

func main() {
	flag.String("dbconn", "postgres://postgres:qwerty@saverbate-db:5432/saverbate_records?sslmode=disable", "Database connection string")
	flag.String("redisAddress", "saverbate-redis:6379", "Address to redis server")

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
	viper.Set("user", f.PortalUsername)
	viper.Set("password", f.PortalPassword)

	allocOpts := chromedp.DefaultExecAllocatorOptions[:]
	allocOpts = append(allocOpts, chromedp.DisableGPU)
	allocOpts = append(allocOpts, chromedp.Flag("headless", false))

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), allocOpts...)
	bctx, _ := chromedp.NewContext(allocCtx)
	defer cancel()

	var performerNodes []*cdp.Node

	err = chromedp.Run(bctx,
		chromedp.Navigate(`https://chaturbate.com/`),
		chromedp.WaitVisible(`#close_entrance_terms`),
		chromedp.Sleep(1*time.Second),
		chromedp.Click(`#close_entrance_terms`),
		chromedp.Sleep(2*time.Second),
		chromedp.Click(`a.login-link`),
		chromedp.WaitVisible(`//input[@name="username"]`),
		chromedp.SendKeys(`//input[@name="username"]`, viper.GetString("user")),
		chromedp.Sleep(1*time.Second),
		chromedp.SendKeys(`//input[@name="password"]`, viper.GetString("password")),
		chromedp.Submit(`//input[@name="password"]`),
		chromedp.WaitVisible(`a#followed_anchor`),
		chromedp.Click(`a#followed_anchor`),
		chromedp.Sleep(1*time.Second),
		chromedp.WaitVisible(`.followedDropdown`),
		chromedp.Click(`//a[contains(text(),'Show All')][1]`),
		chromedp.WaitVisible(`#main`),
		chromedp.Sleep(2*time.Second),
		chromedp.Click(`//a[contains(text(),'Offline Rooms')][1]`),
		chromedp.WaitVisible(`#main`),
		chromedp.Sleep(1*time.Second),
		chromedp.Nodes(`li.room_list_room .details .title a`, &performerNodes),
	)
	if err != nil {
		log.Fatal(err)
	}

	for len(performerNodes) > 0 {
		err = chromedp.Run(bctx,
			chromedp.Click(`//li[@class="room_list_room"][1]/a`),
			chromedp.WaitVisible(`#main`),
			chromedp.Sleep(5*time.Second),
			chromedp.Click(`//span[contains(text(),'- UNFOLLOW')][1]`),
			chromedp.Sleep(7*time.Second),
			chromedp.Click(`//a[@href="/followed-cams/"][1]`),
			chromedp.WaitVisible(`.followedDropdown`),
			chromedp.Click(`//a[contains(text(),'Show All')][1]`),
			chromedp.WaitVisible(`#main`),
			chromedp.Sleep(2*time.Second),
			chromedp.Click(`//a[contains(text(),'Offline Rooms')][1]`),
			chromedp.WaitVisible(`#main`),
			chromedp.Sleep(3*time.Second),
			chromedp.Nodes(`li.room_list_room .details .title a`, &performerNodes),
		)
		if err != nil {
			log.Fatal(err)
		}
	}

	// err = chromedp.Run(bctx,
	// 	chromedp.Navigate(`https://chaturbate.com/auth/login`),
	// 	chromedp.WaitVisible(`//input[@name="username"]`),
	// 	chromedp.SendKeys(`//input[@name="username"]`, viper.GetString("user")),
	// 	chromedp.SendKeys(`//input[@name="password"]`, viper.GetString("password")),
	// 	chromedp.Submit(`//input[@name="password"]`),
	// 	chromedp.WaitVisible(`#close_entrance_terms`),
	// 	chromedp.Click(`#close_entrance_terms`),
	// 	chromedp.Navigate(`https://chaturbate.com/browniezuza/`),
	// 	chromedp.WaitVisible(`//video`),
	// 	chromedp.Click(`//span[contains(text(),'+ FOLLOW')][1]`),
	// )
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Make a redis pool
	// redisPool := &redis.Pool{
	// 	MaxActive: 5,
	// 	MaxIdle:   5,
	// 	Wait:      true,
	// 	Dial: func() (redis.Conn, error) {
	// 		return redis.Dial("tcp", viper.GetString("redisAddress"))
	// 	},
	// }

	// pool := work.NewWorkerPool(subscriber.Context{}, concurrency, subscriber.RedisNamespace, redisPool)
	// pool.PeriodicallyEnqueue("0 0 */4 * * *", subscriber.JobName)

	// pool.Job(subscriber.JobName, (*subscriber.Context).FollowPerformers)

	// Start processing jobs
	// pool.Start()

	// Wait for a signal to quit:
	signalChan := make(chan os.Signal, 1)
	// SIGTERM is called when Ctrl+C was pressed
	signal.Notify(signalChan, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-signalChan

	// Stop the pool
	// pool.Stop()

	err = farm.Release(db, serviceName, viper.GetString("user"))
	if err != nil {
		log.Fatalf("Failed to release farm: %v\n", err)
	}

	db.Close()
}
