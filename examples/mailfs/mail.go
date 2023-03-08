package main

// import (
// 	"log"
// 	"net/mail"

// 	imap "github.com/emersion/go-imap"
// 	"github.com/emersion/go-imap/client"
// )

// type ServerGmail struct {
// 	user   string
// 	pass   string
// 	erro   string
// 	client *client.Client
// }

// func NewServerGmail() *ServerGmail {
// 	serverGmail := &ServerGmail{}
// 	serverGmail.user = "xxxxxx@gmail.com"
// 	serverGmail.pass = "xxxxx"
// 	serverGmail.erro = ""

// 	return serverGmail
// }

// func (serverGmail *ServerGmail) Connect() {
// 	// Connect to server
// 	cliente, erro := client.DialTLS("smtp.gmail.com:993", nil)
// 	if erro != nil {
// 		serverGmail.erro = erro.Error()
// 	}
// 	log.Println("Connected")

// 	serverGmail.client = cliente

// }

// func (serverGmail *ServerGmail) Login() {
// 	// Login
// 	if erro := serverGmail.client.Login(serverGmail.user, serverGmail.pass); erro != nil {
// 		serverGmail.erro = erro.Error()
// 	}
// 	log.Println("Logged")

// }

// func (serverGmail *ServerGmail) setLabelBox(label string) *imap.MailboxStatus {
// 	mailbox, erro := serverGmail.client.Select(label, true)
// 	if erro != nil {
// 		serverGmail.erro = erro.Error()
// 	}
// 	return mailbox
// }

// func (serverGmail *ServerGmail) ListUnseenMessages() {
// 	// set mailbox to INBOX
// 	serverGmail.setLabelBox("INBOX")
// 	// criteria to search for unseen messages
// 	criteria := imap.NewSearchCriteria()
// 	criteria.WithoutFlags = []string{"\\Seen"}

// 	uids, err := serverGmail.client.UidSearch(criteria)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	seqSet := new(imap.SeqSet)
// 	seqSet.AddNum(uids...)
// 	section := &imap.BodySectionName{}
// 	items := []imap.FetchItem{imap.FetchEnvelope, imap.FetchFlags, imap.FetchInternalDate, section.FetchItem()}
// 	messages := make(chan *imap.Message)
// 	go func() {
// 		if err := serverGmail.client.UidFetch(seqSet, items, messages); err != nil {
// 			log.Fatal(err)
// 		}
// 	}()

// 	for message := range messages {

// 		log.Println(message.Uid)

// 		if message == nil {
// 			log.Fatal("Server didn't returned message")
// 		}

// 		r := message.GetBody(section)
// 		if r == nil {
// 			log.Fatal("Server didn't returned message body")
// 		}

// 		// Create a new mail reader
// 		mr, err := mail.CreateReader(r)
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		// Print some info about the message
// 		header := mr.Header

// 		if date, err := header.Date(); err == nil {
// 			log.Println("Date:", date)
// 		}
// 		if from, err := header.AddressList("From"); err == nil {
// 			log.Println("From:", from)
// 		}
// 		if to, err := header.AddressList("To"); err == nil {
// 			log.Println("To:", to)
// 		}
// 		if subject, err := header.Subject(); err == nil {
// 			log.Println("Subject:", subject)
// 		}

// 		// MARK "SEEN" ------- STARTS HERE  ---------

// 		seqSet.Clear()
// 		seqSet.AddNum(message.Uid)
// 		item := imap.FormatFlagsOp(imap.AddFlags, true)
// 		flags := []interface{}{imap.SeenFlag}
// 		erro := serverGmail.client.UidStore(seqSet, item, flags, nil)
// 		if erro != nil {
// 			panic("error!")
// 		}
// 	}
// }
