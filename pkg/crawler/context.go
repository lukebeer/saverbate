package crawler

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os/exec"
	"sync"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/spf13/viper"
)

// Context is context of crawler service
type Context struct {
	rs         *redsync.Redsync
	mutexes    map[string]*redsync.Mutex
	guardMutex *sync.Mutex
}

// New creates new instance of crawler context
func New(rs *redsync.Redsync) *Context {
	return &Context{
		rs:         rs,
		mutexes:    make(map[string]*redsync.Mutex),
		guardMutex: &sync.Mutex{},
	}
}

// Crawl retreives new performers data
func (ctx *Context) Crawl(name string) {
	ctx.guardMutex.Lock()
	if _, ok := ctx.mutexes[name]; !ok {
		ctx.mutexes[name] = ctx.rs.NewMutex(name, redsync.WithExpiry(36*time.Hour))
	}
	ctx.guardMutex.Unlock()

	if err := ctx.mutexes[name].Lock(); err != nil {
		log.Println("Crawler already run")
		return
	}

	defer func() {
		if ok, err := ctx.mutexes[name].Unlock(); !ok || err != nil {
			log.Println("Could not release crawler lock")
			return
		}
	}()

	cmd := exec.Command("scrapy", "crawl", name)
	cmd.Env = append(cmd.Env, "DB_CONN="+viper.GetString("pydbconn"))
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	err = cmd.Start()
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	go copyOutput(stdout)
	go copyOutput(stderr)
	cmd.Wait()
}

func copyOutput(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
