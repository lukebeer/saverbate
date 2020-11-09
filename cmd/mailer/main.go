package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"saverbate/pkg/farm"
	"saverbate/pkg/mailer"

	"github.com/jmoiron/sqlx"
	"github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/robfig/cron/v3"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"

	_ "github.com/lib/pq"
)

const serviceName = "mailer"

func main() {
	log.Println("Start mailer service...")

	quit := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	// SIGTERM is called when Ctrl+C was pressed
	signal.Notify(quit, os.Interrupt, os.Kill, syscall.SIGTERM)

	flag.String("listen", "0.0.0.0:80", "Address for listening")
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

	// Configure the HTTP server
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	server := &http.Server{
		Addr:              viper.GetString("listen"),
		Handler:           mux,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      10 * time.Second,
	}

	// Handle shutdown
	// Register onShutdown func for close streamsHandler (notify all streams and viewers about shutdown)
	// and Websocket hub
	server.RegisterOnShutdown(func() {
		c.Stop()
		nc.Drain()

		err = farm.Release(db, serviceName, viper.GetString("user"))
		if err != nil {
			log.Printf("Failed to release farm: %v\n", err)
		}
		db.Close()

		close(done)
	})

	// Shutdown the HTTP server
	go func() {
		<-quit
		log.Println("Server is going shutting down...")

		// Wait 30 seconds for close http connections
		waitIdleConnCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(waitIdleConnCtx); err != nil {
			log.Fatalf("Cannot gracefully shutdown the server: %v\n", err)
		}
	}()

	// Start HTTP server
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server has been closed immediatelly: %v\n", err)
	}

	<-done
	log.Println("Server stopped")
}
