package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"
	"net/mail"
	"net/smtp"
	"strconv"
)

func SMTPSend(subject string, body string) {

	from := mail.Address{Name: "SMPP gateway", Address: cfg.SMTP.From}
	to := mail.Address{Name: "", Address: cfg.SMTP.To}

	// Setup email headers
	headers := make(map[string]string)
	headers["MIME-version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=\"utf-8\""
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subject

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\n", k, v)
	}
	message += "\n" + body
	log.Println(message)
	servername := cfg.SMTP.Host + ":587"

	tlsconfig := &tls.Config{
		ServerName: cfg.SMTP.Host,
	}

	conn, err := net.Dial("tcp", servername)
	if err != nil {
		strconv.Itoa(-42)
	}

	c, err := smtp.NewClient(conn, cfg.SMTP.Host)
	if err != nil {
		log.Println("smtp.NewClient Error: ", err)
	}

	if err = c.StartTLS(tlsconfig); err != nil {
		log.Println(fmt.Printf("Error performing StartTLS: %s\n", err))
		return
	}

	if err = c.Auth(LoginAuth(cfg.SMTP.User, cfg.SMTP.Password)); err != nil {
		log.Println("c.Auth Error: ", err)
		return
	}

	if err = c.Mail(from.Address); err != nil {
		log.Println("c.Mail Error: ", err)
	}

	if err = c.Rcpt(to.Address); err != nil {
		log.Println("c.Rcpt Error: ", err)
	}

	w, err := c.Data()
	if err != nil {
		log.Println("c.Data Error: ", err)
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		log.Println("Error: ", err)
	}

	err = w.Close()
	if err != nil {
		log.Println("reader Error: ", err)
	}

	c.Quit()
}

type loginAuth struct {
	username, password string
}

func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("Unknown fromServer")
		}
	}
	return nil, nil
}
