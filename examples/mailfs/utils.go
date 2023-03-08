package main

import (
	"io"
	"log"
	"os"

	"github.com/emersion/go-imap"
	id "github.com/emersion/go-imap-id"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/charset"
	"github.com/emersion/go-message/mail"
)

func login() *client.Client {
	imap.CharsetReader = charset.Reader

	log.Println("Connecting to server...")

	// Connect to server
	c, err := client.DialTLS(os.Getenv("MAIL_IMAP_ADDR"), nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected")

	// Don't forget to logout
	// defer c.Logout()

	// ID
	idClient := id.NewClient(c)
	cid, err := idClient.ID(id.ID{
		id.FieldName:    "Dev",
		id.FieldVersion: "0.0.1",
		id.FieldVendor:  "ArkJit",
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("ID Info: %s", cid)

	// Login
	if err := c.Login(
		os.Getenv("MAIL_USERNAME"), os.Getenv("MAIL_PASSWORD"),
	); err != nil {
		log.Fatal(err)
	}
	log.Println("Logged in")

	return c
}

func parseMsg(c *client.Client, seqNum uint32) {
	seqSet := new(imap.SeqSet)
	seqSet.AddNum(seqNum)

	// Get the whole message body
	var section imap.BodySectionName
	items := []imap.FetchItem{section.FetchItem()}

	messages := make(chan *imap.Message, 1)
	go func() {
		if err := c.Fetch(seqSet, items, messages); err != nil {
			log.Fatal(err)
		}
	}()

	msg := <-messages
	if msg == nil {
		log.Fatal("Server didn't returned message")
	}

	r := msg.GetBody(&section)
	if r == nil {
		log.Fatal("Server didn't returned message body")
	}

	// Create a new mail reader
	mr, err := mail.CreateReader(r)
	if err != nil {
		log.Fatal(err)
	}

	// Print some info about the message
	header := mr.Header
	if date, err := header.Date(); err == nil {
		log.Println("Date:", date)
	}
	if from, err := header.AddressList("From"); err == nil {
		log.Println("From:", from)
	}
	if to, err := header.AddressList("To"); err == nil {
		log.Println("To:", to)
	}
	if subject, err := header.Subject(); err == nil {
		log.Println("Subject:", subject)
	}

	// Process each message's part
	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		switch h := p.Header.(type) {
		case *mail.InlineHeader:
			// This is the message's text (can be plain-text or HTML)
			b, _ := io.ReadAll(p.Body)
			log.Printf("Got text: %v\n", string(b))
		case *mail.AttachmentHeader:
			// This is an attachment
			filename, _ := h.Filename()
			body, err := io.ReadAll(p.Body)
			if err != nil {
				log.Printf("Got attachment body err: %s\n", err)
			}
			log.Printf("Got attachment: %v, size: %d\n", filename, len(body))
		}
	}
}
