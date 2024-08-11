package utils

import (
	"log/slog"
	gomail "gopkg.in/gomail.v2"
)

func SendMessage(body string, email, passw string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", email)
	msg.SetHeader("To", email)
	msg.SetBody("text", body)

	n := gomail.NewDialer("smtp.gmail.com", 587, email, passw)

	if err := n.DialAndSend(msg); err != nil {
		slog.Error("Can't send the message", err)
		return err
	}
	slog.Info("The message has been succesfully send")
	return nil
}
