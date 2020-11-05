package mailer

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"regexp"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
	"github.com/gocraft/work"
	"github.com/nats-io/nats.go"
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
type Context struct{}

var (
	performersLinkRe = regexp.MustCompile(`(https://chaturbate\.com\/[^\/]{3,}/)(?:\n|\]|<||\s|$)`)
	performerNameRe  = regexp.MustCompile(`https://chaturbate\.com\/([^\/]{3,})\/`)
)

// CheckEmail connects to imap server and check new emails
func (ctx *Context) CheckEmail(job *work.Job) error {
	performersUnique := make(map[string]struct{})

	log.Println("Connecting to server...")

	// Connect to server
	c, err := client.DialTLS(imapServer, nil)
	if err != nil {
		log.Printf("Error: %v", err)
		return err
	}

	// Don't forget to logout
	defer c.Logout()

	// Login
	if err := c.Login(viper.GetString("user"), viper.GetString("password")); err != nil {
		log.Printf("Error: %v", err)
		return err
	}
	log.Println("Logged in")

	// Select INBOX
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		log.Printf("Error: %v", err)
		return err
	}
	log.Println("Flags for INBOX:", mbox.Flags)

	if mbox.Messages == 0 {
		log.Printf("Empty mailbox...Exit")
		return nil
	}
	from := uint32(1)
	to := mbox.Messages
	if mbox.Messages > 49 {
		// We're using unsigned integers here, only substract if the result is > 0
		from = mbox.Messages - 49
	}
	seqset := new(imap.SeqSet)
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

			header := mr.Header
			if subject, err := header.Subject(); err == nil {
				matched, _ := regexp.MatchString(
					`(?:someone you follow is chaturbating|broadcasters you follow are chaturbating)$`,
					subject,
				)
				if !matched {
					break
				}
			}

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
					matches := performersLinkRe.FindAllStringSubmatch(string(b), -1)
					if len(matches) == 0 {
						break
					}

					for _, matching := range matches {
						if _, ok := performersUnique[matching[1]]; !ok {
							performersUnique[matching[1]] = struct{}{}
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
		log.Printf("Error: %v", err)
		return err
	}

	log.Println("Done!")

	log.Println("Connecting to NATS...")
	// NATS
	nc, err := nats.Connect(viper.GetString("natsAddress"), nats.NoEcho())
	if err != nil {
		log.Printf("Error: %v", err)
		return err
	}
	defer nc.Drain()

	if len(performersUnique) == 0 {
		log.Println("Nothing here...")
		return nil
	}

	for performer := range performersUnique {
		name := performerNameRe.FindStringSubmatch(performer)
		if len(name) == 0 {
			continue
		}

		message, err := json.Marshal(struct {
			Name string `json:"performer_name"`
		}{
			Name: name[1],
		})
		if err != nil {
			log.Printf("marshaling message failed: %v\n", err)
			return nil
		}

		if err := nc.Publish("downloading", message); err != nil {
			log.Printf("Failed to publush message: %v", err)
			return nil
		}
	}

	return nil
}
