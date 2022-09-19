package main

import (
	"log"

	mail "github.com/xhit/go-simple-mail"
)

// SendMails - send emails
//Authentication type: 1 = AuthPlain, 2 = AuthLogin, 3 = AuthCRAMMD5
//Encryption type: 1 = EncryptionNone, 2 = EncryptionSSL, 3 = EncryptionTLS
func SendMails(host string, port int, authenticationType int, encryptionType int, user string, password string, from string, to []string, subject string, body string, attach string) (err error) {

	server := mail.NewSMTPClient()

	server.Host = host
	server.Port = port
	server.Username = user
	server.Password = password
	server.KeepAlive = false

	switch authenticationType {
	case 1:
		server.Authentication = mail.AuthPlain
	case 2:
		server.Authentication = mail.AuthLogin
	case 3:
		server.Authentication = mail.AuthCRAMMD5
	}

	switch encryptionType {
	case 1:
		server.Encryption = mail.EncryptionNone
	case 2:
		server.Encryption = mail.EncryptionSSL
	case 3:
		server.Encryption = mail.EncryptionTLS
	}

	smtpClient, err := server.Connect()
	if err != nil {
		log.Println(err)
		return err
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
		return err
	}

	log.Println("Email Sent")
	return nil
}

// SendMail - send email
//Authentication type: 1 = AuthPlain, 2 = AuthLogin, 3 = AuthCRAMMD5
//Encryption type: 1 = EncryptionNone, 2 = EncryptionSSL, 3 = EncryptionTLS
func SendMail(host string, port int, authenticationType int, encryptionType int, user string, password string, from string, to string, subject string, body string, attach string) (err error) {

	server := mail.NewSMTPClient()

	server.Host = host
	server.Port = port
	server.Username = user
	server.Password = password
	server.KeepAlive = false

	switch authenticationType {
	case 1:
		server.Authentication = mail.AuthPlain
	case 2:
		server.Authentication = mail.AuthLogin
	case 3:
		server.Authentication = mail.AuthCRAMMD5
	}

	switch encryptionType {
	case 1:
		server.Encryption = mail.EncryptionNone
	case 2:
		server.Encryption = mail.EncryptionSSL
	case 3:
		server.Encryption = mail.EncryptionTLS
	}

	smtpClient, err := server.Connect()
	if err != nil {
		log.Fatal(err)
		return err
	}

	email := mail.NewMSG()

	email.SetFrom(from)
	email.AddTo(to)
	email.SetSubject(subject)
	email.SetBody(mail.TextHTML, body)

	err = email.Send(smtpClient)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("Email Sent")
	return nil
}
