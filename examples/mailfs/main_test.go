package main

import (
	"log"
	"testing"

	"github.com/emersion/go-imap"
)

// Select INBOX 已发送 草稿箱 MailDisk
const boxName = "Drafts"

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

	// // Write the message to a buffer
	// var b bytes.Buffer
	// // b.WriteString("From: <...@gmail.com>\r\n")
	// // b.WriteString("To: <...@gmail.com>\r\n")
	// fmt.Fprintf(&b, "Subject: No.%d\r\n", rand.Int())
	// b.WriteString("\r\n")
	// // Message body
	// b.WriteString("Append test using IMAP and Draft folder")

	// // Append it to Drafts
	// if err := c.Append(boxName, nil, time.Now(), &b); err != nil {
	// 	log.Fatal(err)
	// }

	// // Store it to Drafts
	// header := mail.AttachmentHeader{
	// 	Header: message.Header{},
	// }
	// header.SetFilename("README.md")
	// attachment, err := message.New(header.Header, strings.NewReader("Hello Go"))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	updates := make(chan *imap.Message, 1)
	// if msg, err := message.NewMultipart(message.Header{}, []*message.Entity{attachment}); err != nil {
	// 	imap.NewMessage(999, []imap.FetchItem{})
	// 	updates <- msg
	// }

	seqSet := new(imap.SeqSet)
	seqSet.AddRange(1, 1)
	_, err := c.Select(boxName, false)
	if err != nil {
		log.Fatal(err)
	}
	if err := c.Store(seqSet, imap.DraftFlag, []any{"body", "foobar"}, updates); err != nil {
		log.Fatal(err)
	}

	for m := range updates {
		log.Println(m)
	}
}
