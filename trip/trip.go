package trip

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/DnFreddie/backy/utils"
	gomail "gopkg.in/gomail.v2"
)

type Email_Creds struct {
	Email  string
	Passwd string
}

const (
	MAIL    = "Test"
	E_PASSW = "test"
)

func readTheConfig() (Email_Creds, error) {

	email_conf, err := utils.GetUser(utils.LOG_DIR + "/email.json")

	if err != nil {
		return Email_Creds{}, err
	}

	_, err = os.Stat(email_conf)

	if os.IsNotExist(err) || err != nil {
		fmt.Println("creating the email_conf ", err)

		newCreds := Email_Creds{
			Email:  MAIL,
			Passwd: E_PASSW,
		}
		jr, err := json.Marshal(&newCreds)

		if err != nil {
			slog.Error("cant unmrashall the creds", err)
		}

		err = os.WriteFile(email_conf, jr, 0666)

		if err != nil {
			slog.Error("Can't create a Email Creds file", err)
			return Email_Creds{}, err

		}
	}

	creds, err := utils.ReadJson(email_conf, &Email_Creds{})
	if err != nil {

		slog.Error("Can't read the creds", err)

		return Email_Creds{}, err

	}

	if len(creds) != 1 {
		err := errors.New("Smth is wrong in the config file ")
		return Email_Creds{}, err
	}

	return creds[0], nil
}

func SendMessage(body string) error {
	creds, err := readTheConfig()
	if err != nil {
		slog.Error("Can't read the config", err)
		return err
	}
	msg := gomail.NewMessage()
	msg.SetHeader("From", MAIL)
	msg.SetHeader("To", MAIL)
	msg.SetBody("text", body)

	n := gomail.NewDialer("smtp.gmail.com", 587, creds.Email, creds.Passwd)

	if err := n.DialAndSend(msg); err != nil {
		slog.Error("Can't send the message", err)
		return err
	}
	slog.Info("The message has been succesfully send")
	return nil
}
