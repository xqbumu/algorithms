package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-message/mail"
)

// Select INBOX 已发送 草稿箱 MailDisk
const boxName = "MailDisk"

func TestClient(t *testing.T) {
	c := login()
	// Don't forget to logout
	defer c.Logout()

	// List mailboxes
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.List("", "*", mailboxes)
	}()

	log.Println("Mailboxes:")
	for m := range mailboxes {
		// Mailbox Status
		mbox, err := c.Select(m.Name, false)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("* %s: %d\n", m.Name, mbox.Messages)
	}

	if err := <-done; err != nil {
		log.Fatal(err)
	}
	log.Println("List MailBox Done")

	mbox, err := c.Select(boxName, false)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Flags for %s: %s", boxName, mbox.Flags)

	// Get the last {cnt} messages
	var cnt uint32 = 20
	from := uint32(1)
	to := mbox.Messages
	if mbox.Messages > cnt {
		// We're using unsigned integers here, only subtract if the result is > 0
		from = mbox.Messages - cnt
	}
	seqSet := new(imap.SeqSet)
	seqSet.AddRange(from, to)

	messages := make(chan *imap.Message, 10)
	go func() {
		if err := c.Fetch(seqSet, []imap.FetchItem{imap.FetchFull}, messages); err != nil {
			log.Fatal(err)
		}
	}()

	log.Printf("Last %d messages:\n", cnt)
	var seqNum uint32 = 0
	for msg := range messages {
		log.Printf("Current Subject(%d): %s\n", msg.SeqNum, msg.Envelope.Subject)
		// if msg.Envelope.Subject == "Undeliverable: 岗位应聘  文案/策划" {
		seqNum = msg.SeqNum
		// }
	}

	if seqNum > 0 {
		parseMsg(c, seqNum)
	}

	log.Println("Done!")
}

func TestSTMP(t *testing.T) {
	c := login()
	// Don't forget to logout
	defer c.Logout()

	b := generateMail()

	// Append it to Drafts
	if err := c.Append(boxName, nil, time.Now(), &b); err != nil {
		log.Fatal(err)
	}
}

func generateMail() bytes.Buffer {
	// Write the message to a buffer
	var b bytes.Buffer

	from := []*mail.Address{{Name: "Dev", Address: "dev@example.org"}}
	to := []*mail.Address{{Name: "Customer", Address: "customer@example.org"}}

	// Create our mail header
	var h mail.Header
	h.SetDate(time.Now())
	h.SetAddressList("From", from)
	h.SetAddressList("To", to)
	// Set Subject
	h.SetSubject(fmt.Sprintf("No.%d\r\n", rand.Int()))
	h.Add("X-Mailfs-Enbale", "true")
	h.Add("X-Mailfs-Type", "file")

	// Create a new mail writer
	mw, err := mail.CreateWriter(&b, h)
	if err != nil {
		log.Fatal(err)
	}

	// Create a text part
	tw, err := mw.CreateInline()
	if err != nil {
		log.Fatal(err)
	}
	var th mail.InlineHeader
	th.Set("Content-Type", "text/plain")
	w, err := tw.CreatePart(th)
	if err != nil {
		log.Fatal(err)
	}
	io.WriteString(w, "Who are you?")
	w.Close()
	tw.Close()

	// Create first attachment
	var ah mail.AttachmentHeader
	ah.Set("Content-Type", "test/plain")
	ah.SetFilename("main.go")
	w, err = mw.CreateAttachment(ah)
	if err != nil {
		log.Fatal(err)
	}
	// TODO: write a Go file to w
	fp, err := os.Open("main.go")
	if err != nil {
		log.Fatal(err)
	}
	written, err := io.Copy(w, fp)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("written %d", written)
	w.Close()

	// Create second attachment
	ah.Set("Content-Type", "test/plain")
	ah.SetFilename("utils.go")
	w, err = mw.CreateAttachment(ah)
	if err != nil {
		log.Fatal(err)
	}
	// TODO: write a Go file to w
	fp, err = os.Open("utils.go")
	if err != nil {
		log.Fatal(err)
	}
	written, err = io.Copy(w, fp)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("written %d", written)
	w.Close()

	mw.Close()

	return b
}
