package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"github.com/volatiletech/authboss/v3"

	nats "github.com/nats-io/nats.go"

	appMiddleware "saverbate/pkg/middleware"

	flag "github.com/spf13/pflag"

	"saverbate/pkg/handler"
	"saverbate/pkg/user"

	_ "github.com/lib/pq"
	_ "github.com/volatiletech/authboss/v3/auth"
	_ "github.com/volatiletech/authboss/v3/logout"
	_ "github.com/volatiletech/authboss/v3/register"
)

func main() {
	flag.String("dbconn", "postgres://postgres:qwerty@localhost:10532/saverbate_records?sslmode=disable", "Database connection string")
	flag.String("listen", "0.0.0.0:8085", "Listening address for http server")

	flag.String("cookieStoreKey", "", "Secret key for cookies storage")
	flag.String("sessionStoreKey", "", "Secret key for session storage")

	flag.String("rootURL", "http://localhost:8085", "Root URL")
	flag.String("natsAddress", "nats://localhost:10222", "Address to connect to NATS server")
	flag.Parse()

	viper.BindPFlags(flag.CommandLine)

	db, err := sqlx.Connect("postgres", viper.GetString("dbconn"))
	if err != nil {
		log.Fatalf("Failed to connect to database: %v\n", err)
	}

	var ab *authboss.Authboss
	ab, err = user.InitAuthBoss(db)
	if err != nil {
		log.Fatalf("Failed to init authbos: %v\n", err)
	}

	// NATS
	nc, err := nats.Connect(viper.GetString("natsAddress"), nats.NoEcho())
	if err != nil {
		log.Panic(err)
	}
	r := chi.NewRouter()
	// Some basic middlewares
	r.Use(
		middleware.RealIP,
		middleware.Logger,
		ab.LoadClientStateMiddleware,
		appMiddleware.CurrentUserDataInject(ab),
		//appMiddleware.ConfigDataInject(),
		middleware.Recoverer,
	)
	// Homepage
	r.Method("GET", "/", handler.NewHomepageHandler(db))
	// Registration
	//r.Method("GET", "/register", handler.NewApplicationHandler("register"))
	// Login
	//r.Method("GET", "/login", handler.NewApplicationHandler("login"))
	// Show
	//r.Method("GET", "/show", handler.NewApplicationHandler("show"))

	//r.Group(func(r chi.Router) {
	//	r.Use(ab.LoadClientStateMiddleware, authboss.Middleware2(ab, authboss.RequireFullAuth, authboss.RespondRedirect))

	// Broadcast
	// r.Method("GET", "/broadcast", handler.NewApplicationHandler("broadcast"))
	//})

	r.Group(func(r chi.Router) {
		r.Use(ab.LoadClientStateMiddleware, authboss.ModuleListMiddleware(ab))
		r.Mount("/auth", http.StripPrefix("/auth", ab.Config.Core.Router))
	})

	quit := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	// Serve static assets
	// serves files from web/static dir
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	staticPrefix := "/static/"
	staticDir := path.Join(cwd, "web", staticPrefix)
	r.Method("GET", staticPrefix+"*", http.StripPrefix(staticPrefix, http.FileServer(http.Dir(staticDir))))

	// Favicon
	r.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		if err := serveStaticFile(staticDir+"/favicon.ico", w); err != nil {
			log.Println(err)
		}
	})

	// Handle 404
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)

		if err := serveStaticFile(staticDir+"/404.html", w); err != nil {
			log.Println(err)
		}
	})

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

		log.Println("Close NATS connection...")
		nc.Drain()

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

func serveStaticFile(filePath string, w io.Writer) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}

	buf := make([]byte, 4*1024) // 4Kb
	if _, err = io.CopyBuffer(w, f, buf); err != nil {
		return err
	}

	return nil
}
