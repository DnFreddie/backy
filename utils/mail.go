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

func readTheConfig(c Email_Creds) (Email_Creds, error) {

	email_conf, err := GetUser(LOG_DIR + "/email.json")

	if err != nil {
		return c, err
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
			return Email_Creds{}, err

		}
	}

	creds, err := ReadJson(email_conf, &Email_Creds{})
	if err != nil {

		slog.Error("Can't read the creds", err)

		return c, err

	}

	if len(creds) != 1 {
		err := errors.New("Smth is wrong in the config file ")
		return c, err
	}

	return creds[0], nil
}

func SendMessage(body string, creds Email_Creds) error {
	creds, err := readTheConfig(creds)
	if err != nil {
		slog.Error("Can't read the config", err)
		return err
	}
	msg := gomail.NewMessage()
	msg.SetHeader("From", creds.Email)
	msg.SetHeader("To", creds.Email)
	msg.SetBody("text", body)

	n := gomail.NewDialer("smtp.gmail.com", 587, creds.Email, creds.Passwd)

	if err := n.DialAndSend(msg); err != nil {
		slog.Error("Can't send the message", err)
		return err
	}
	slog.Info("The message has been succesfully send")
	return nil
}
