package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	gomail "gopkg.in/gomail.v2"
	"log/slog"
	"os"
)

type Email_Creds struct {
	Email  string
	Passwd string
}

func (c *Email_Creds) readTheConfig() error {

	email_conf, err := GetUser(LOG_DIR + "/email.json")

	if err != nil {
		return err
	}

	_, err = os.Stat(email_conf)

	if os.IsNotExist(err) || err != nil {
		fmt.Println("creating the email_conf ", err)
		jr, err := json.Marshal(&c)

		if err != nil {
			slog.Error("cant unmrashall the creds", err)
		}

		err = os.WriteFile(email_conf, jr, 0666)

		if err != nil {
			slog.Error("Can't create a Email Creds file", err)
			return err

		}
	}

	var credsA []Email_Creds
	err = ReadJson(email_conf, &credsA)

	if err != nil {

		slog.Error("Can't read the creds", err)

		return err

	}

	if len(credsA) != 1 {
		err := errors.New("Smth is wrong in the config file ")
		return err
	}

	newCreds := credsA[0]
	c.Email = newCreds.Email
	c.Passwd = newCreds.Passwd

	return nil
}

func SendMessage(body string, email ,passw  string ) error {
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
