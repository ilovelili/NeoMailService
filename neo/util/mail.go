package util

import (
	"encoding/json"
	"fmt"
	"log"
	"neo/config"
	"neo/core"
	"net/smtp"
	"strings"
)

var (
	receivers = strings.Split(config.GetConfig().Mail.Receiver, ",")
)

const (
	// smtpServer SMTP server
	smtpServer string = "smtp.gmail.com"
	smtpPort   string = "587"
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
	mail.subject = fmt.Sprintf("AgentSmith: Neo price is $%g now (%g percent %s)", rate.Neo.Usd, rate.Neo.UsdRate, resolveUpOrDown(rate.Neo.UsdRate))
	body, err := json.Marshal(rate)
	if err != nil {
		return err
	}
	mail.body = string(body)
	messageBody := mail.buildMessage()
	smtpServer := SMTPServer{host: smtpServer, port: smtpPort, user: config.Sender.Account, password: config.Sender.Password}

	log.Println("sending mail")
	return smtp.SendMail(smtpServer.ServerName(), smtpServer.Auth(), mail.senderID, mail.toIds, []byte(messageBody))
}

func resolveUpOrDown(rate float32) string {
	if rate > 0 {
		return "up"
	}
	return "down"
}
