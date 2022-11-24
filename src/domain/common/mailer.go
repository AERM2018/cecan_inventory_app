package common

import (
	"net/smtp"
	"os"
)

func SendPlainEmail(to string, subject string, body string) error {
	password := os.Getenv("EMAIL_APP_PASSWD")
	from := os.Getenv("EMAIL_SENDER")
	toEmailAddresses := []string{to}
	host := "smtp.gmail.com"
	port := "587"
	address := host + ":" + port
	message := []byte("From:" + from + "\r\n" +
		"To:" + to + "\r\n" +
		"Subject:" + subject + "\r\n\r\n" +
		body + "\r\n")
	auth := smtp.PlainAuth("", from, password, host)
	err := smtp.SendMail(address, auth, from, toEmailAddresses, message)
	if err != nil {
		return err
	}
	return nil
}
