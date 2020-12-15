package mailer

import (
	"io"
	"io/ioutil"
	"log"
	"regexp"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
	"github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
)

const (
	imapServer = "imap.gmail.com:993"
	// NamespaceDownloads is redis namespace
	NamespaceDownloads = "saverbate_downloads"
	// JobName is job name for periodically enqueuing
	JobName = "check_email"
)

// Context is context of mailer service
type Context struct {
	nc      *nats.Conn
	Metrics Metrics
}

var (
	performersLinkRe    = regexp.MustCompile(`(https:\/\/chaturbate\.com\/[^\/]{3,}/)(?:\n|\]|<||\s|$)`)
	performerIsOnlineRe = regexp.MustCompile(`strong>([^\s]{3,}) is now online`)
	performerNameRe     = regexp.MustCompile(`https:\/\/chaturbate\.com\/([^\/]{3,})\/`)
	subjectRe           = regexp.MustCompile(`(?:someone you follow is chaturbating|broadcasters you follow are chaturbating)$`)
)

func New(nc *nats.Conn) *Context {
	c := &Context{nc: nc}
	c.Metrics.parsedPerformers = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "saverbate_mailer_parsed_performers_total",
			Help: "The total number of processed from mailbox performers who online",
		},
		[]string{"status"},
	)
	prometheus.MustRegister(c.Metrics.parsedPerformers)
	return c
}

// CheckEmail connects to imap server and check new emails
func (ctx *Context) CheckEmail() error {
	var seqset *imap.SeqSet

	performersUnique := make(map[string]struct{})

	log.Println("Connecting to IMAP server...")

	// Connect to server
	c, err := client.DialTLS(imapServer, nil)
	if err != nil {
		return err
	}

	// Login
	if err := c.Login(viper.GetString("user"), viper.GetString("password")); err != nil {
		return err
	}

	defer c.Logout()

	// Select INBOX
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		return err
	}

	if mbox.Messages == 0 {
		log.Println("Empty mailbox...Exit")
		return nil
	}
	from := uint32(1)
	to := mbox.Messages
	if mbox.Messages > 49 {
		// We're using unsigned integers here, only substract if the result is > 0
		from = mbox.Messages - 49
	}
	seqset = new(imap.SeqSet)
	seqset.AddRange(from, to)

	var section imap.BodySectionName
	items := []imap.FetchItem{section.FetchItem()}

	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.Fetch(seqset, items, messages)
	}()

	if messages == nil {
		log.Println("Server didn't returned message")
		return nil
	}

	for msg := range messages {
		for s := range msg.Body {
			r := msg.GetBody(s)
			if r == nil {
				log.Println("No body for given section")
				break
			}

			mr, err := mail.CreateReader(r)
			if err != nil {
				log.Printf("Error: %v", err)
				break
			}

			ctx.Metrics.parsedPerformers.WithLabelValues("messages_got").Inc()

			header := mr.Header
			if subject, err := header.Subject(); err == nil {
				matched := subjectRe.MatchString(subject)
				if !matched {
					break
				}
			}

			ctx.Metrics.parsedPerformers.WithLabelValues("messages_matched").Inc()

			// Process each message's part
			for {
				p, err := mr.NextPart()
				if err == io.EOF {
					break
				} else if err != nil {
					log.Printf("Error: %v", err)
					return nil
				}

				switch p.Header.(type) {
				case *mail.InlineHeader:
					// This is the message's text (can be plain-text or HTML)
					b, _ := ioutil.ReadAll(p.Body)
					strBody := string(b)

					matches := performersLinkRe.FindAllStringSubmatch(strBody, -1)

					for _, matching := range matches {
						if _, ok := performersUnique[matching[1]]; !ok {
							performersUnique[matching[1]] = struct{}{}
							ctx.Metrics.parsedPerformers.WithLabelValues("performers_parsed").Inc()
						}
					}

					matches = performerIsOnlineRe.FindAllStringSubmatch(strBody, -1)

					for _, matching := range matches {
						m := "https://chaturbate.com/" + matching[1] + "/"
						if _, ok := performersUnique[m]; !ok {
							performersUnique[m] = struct{}{}
							ctx.Metrics.parsedPerformers.WithLabelValues("performers_parsed").Inc()
						}
					}
				case *mail.AttachmentHeader:
					// This is an attachment
					// filename, _ := h.Filename()
					// log.Printf("Got attachment: %v", filename)
				}
			}
		}
	}

	if err := <-done; err != nil {
		return err
	}

	log.Println("Clean up...")

	// Delete messages
	storeItems := imap.FormatFlagsOp(imap.AddFlags, true)

	if err := c.Store(seqset, storeItems, []interface{}{imap.DeletedFlag}, nil); err != nil {
		return err
	}
	if err := c.Expunge(nil); err != nil {
		return err
	}

	log.Println("Done!")

	for performer := range performersUnique {
		name := performerNameRe.FindStringSubmatch(performer)
		if len(name) == 0 {
			continue
		}

		if err := ctx.nc.Publish("downloading", []byte(name[1])); err != nil {
			return err
		}

		if err := ctx.nc.Flush(); err != nil {
			return err
		}
	}

	return nil
}
