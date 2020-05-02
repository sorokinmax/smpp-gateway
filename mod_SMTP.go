package main

import (
	"log"

	mail "github.com/xhit/go-simple-mail"
)

// SendMail - send email
func SendMail(host string, port int, user string, password string, from string, to []string, subject string, body string, attach string) {

	server := mail.NewSMTPClient()

	server.Host = host
	server.Port = port
	server.Username = user
	server.Password = password
	server.Authentication = mail.AuthLogin
	server.Encryption = mail.EncryptionTLS
	server.KeepAlive = false
	smtpClient, err := server.Connect()

	if err != nil {
		log.Fatal(err)
	}

	email := mail.NewMSG()

	email.SetFrom(from)
	email.SetSubject(subject)
	email.SetBody(mail.TextHTML, body)

	for _, v := range to {
		email.AddTo(v)
	}

	err = email.Send(smtpClient)

	if err != nil {
		log.Println(err)
	} else {
		log.Println("Email Sent")
	}
}
