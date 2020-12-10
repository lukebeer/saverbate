package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"github.com/volatiletech/authboss/v3"

	appMiddleware "saverbate/pkg/middleware"

	flag "github.com/spf13/pflag"

	goredislib "github.com/go-redis/redis/v8"

	"saverbate/pkg/handler"
	"saverbate/pkg/user"
	"saverbate/pkg/utils"

	_ "github.com/lib/pq"
	_ "github.com/volatiletech/authboss/v3/auth"
	_ "github.com/volatiletech/authboss/v3/logout"
	_ "github.com/volatiletech/authboss/v3/register"
	"github.com/volatiletech/authboss/v3/remember"
)

func main() {
	flag.String("dbconn", "postgres://postgres:qwerty@localhost:10532/saverbate_records?sslmode=disable", "Database connection string")
	flag.String("listen", "0.0.0.0:8085", "Listening address for http server")

	flag.String("cookieStoreKey", "", "Secret key for cookies storage")
	flag.String("sessionStoreKey", "", "Secret key for session storage")
	flag.Bool("useTLS", false, "Whether use TLS")

	flag.String("staticHostURL", "http://localhost:8086", "Static server host URL")
	flag.String("rootURL", "http://localhost:8085", "Root URL")
	flag.String("redisAddress", "localhost:6379", "Address to redis server")
	flag.Parse()

	viper.BindPFlags(flag.CommandLine)

	redis := goredislib.NewClient(&goredislib.Options{
		Addr: viper.GetString("redisAddress"),
	})
	err := redis.Ping(context.Background()).Err()
	if err != nil {
		log.Fatalf("Failed to connect to redis server: %v\n", err)
	}

	db, err := sqlx.Connect("postgres", viper.GetString("dbconn"))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v\n", err)
	}

	var ab *authboss.Authboss
	ab, err = user.InitAuthBoss(db, redis)
	if err != nil {
		log.Fatalf("Failed to init authbos: %v\n", err)
	}

	r := chi.NewRouter()
	// Some basic middlewares
	r.Use(
		middleware.RealIP,
		middleware.Logger,
		appMiddleware.Nosurfing,

		ab.LoadClientStateMiddleware,
		appMiddleware.CurrentUserDataInject(ab),
		appMiddleware.ConfigDataInject(),
		remember.Middleware(ab),
		middleware.Recoverer,
	)
	// Homepage
	r.Method("GET", "/", handler.NewHomepageHandler(db))
	r.Method("GET", "/records/{uuid}", handler.NewShowHandler(db))

	r.Group(func(r chi.Router) {
		r.Use(ab.LoadClientStateMiddleware, authboss.ModuleListMiddleware(ab))
		r.Mount("/auth", http.StripPrefix("/auth", ab.Config.Core.Router))
	})

	quit := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	// Serve static assets
	// serves files from web/static dir
	staticDir, err := utils.StaticDir()
	if err != nil {
		log.Fatal(err)
	}
	r.Method("GET", utils.StaticPrefix+"*", http.StripPrefix(utils.StaticPrefix, http.FileServer(http.Dir(staticDir))))

	// favicon, robots.txt etc.

	r.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		if err := utils.ServeStaticFile(staticDir+"/favicon.ico", w); err != nil {
			log.Println(err)
		}
	})

	r.Get("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		if err := utils.ServeStaticFile(staticDir+"/robots.txt", w); err != nil {
			log.Println(err)
		}
	})

	r.Get("/site.webmanifest", func(w http.ResponseWriter, r *http.Request) {
		if err := utils.ServeStaticFile(staticDir+"/site.webmanifest", w); err != nil {
			log.Println(err)
		}
	})

	// Handle 404
	r.NotFound(handler.NewNotFoundHandler().ServeHTTP)

	signal.Notify(quit, os.Interrupt)

	// Configure the HTTP server
	server := &http.Server{
		Addr:              viper.GetString("listen"),
		Handler:           r,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      10 * time.Second,
	}

	// Handle shutdown
	server.RegisterOnShutdown(func() {
		log.Println("Close db connection...")
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
