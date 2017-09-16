package util

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"neo/config"
	"neo/core"
	"net"
	"net/smtp"
	"strings"
)

var (
	receivers = strings.Split(config.GetConfig().Mail.Receiver, ",")
)

const (
	// smtpServer SMTP server
	smtpServer string = "smtp.gmail.com"
	smtpPort   string = "465"
)

// Mail mail with sender, receivers, subject and body
type Mail struct {
	senderID string
	toIds    []string
	subject  string
	body     string
}

// SMTPServer Smtp server with host and port
type SMTPServer struct {
	host     string
	port     string
	user     string
	password string
}

// Auth plain auth
func (s *SMTPServer) Auth() smtp.Auth {
	return smtp.PlainAuth("", s.user, s.password, smtpServer)
}

// ServerName resolve server name (host:port)
func (s *SMTPServer) ServerName() string {
	return s.host + ":" + s.port
}

// buildMessage build message body
func (mail *Mail) buildMessage() string {
	message := ""
	message += fmt.Sprintf("From: %s\r\n", mail.senderID)
	if len(mail.toIds) > 0 {
		message += fmt.Sprintf("To: %s\r\n", strings.Join(mail.toIds, ";"))
	}

	message += fmt.Sprintf("Subject: %s\r\n", mail.subject)
	message += "\r\n" + mail.body

	return message
}

// SendMail send mail when fails
func SendMail(rate *core.Rate, config *config.Config) (err error) {
	mail := Mail{}
	mail.senderID = config.Mail.Sender.Account
	mail.toIds = receivers
	mail.subject = fmt.Sprintf("AgentSmith: Neo price %g (%g)", rate.Neo.Usd, rate.Neo.UsdRate)
	body, err := json.Marshal(rate)
	if err != nil {
		return err
	}
	mail.body = string(body)
	messageBody := mail.buildMessage()
	smtpServer := SMTPServer{host: smtpServer, port: smtpPort, user: config.Sender.Account, password: config.Sender.Password}
	serverName := smtpServer.ServerName()

	// TLS config
	host, _, _ := net.SplitHostPort(serverName)
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	conn, err := tls.Dial("tcp", serverName, tlsconfig)
	defer conn.Close()
	if err != nil {
		log.Panic(err)
	}

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		log.Panic(err)
	}

	// Auth
	if err = client.Auth(smtpServer.Auth()); err != nil {
		log.Panic(err)
	}

	// add all from and to
	if err = client.Mail(mail.senderID); err != nil {
		log.Panic(err)
	}

	for _, k := range mail.toIds {
		if err = client.Rcpt(k); err != nil {
			log.Panic(err)
		}
	}

	// Data
	w, err := client.Data()
	if err != nil {
		log.Panic(err)
	}

	_, err = w.Write([]byte(messageBody))
	if err != nil {
		log.Panic(err)
	}

	err = w.Close()
	if err != nil {
		log.Panic(err)
	}

	return
}
