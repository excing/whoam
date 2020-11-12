package main

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

// SesServer Send mail server config
type SesServer struct {
	Ses string `flag:"Send mail server domain url"`
}

var ses SesServer

func init() {
	ses = SesServer{}
	FlagVar(&ses)
}

// SendMail 发送短信
func SendMail(to string, subject string, body string) error {
	formData := url.Values{
		"to":      {to},
		"subject": {subject},
		"body":    {body},
	}

	resp, err := http.PostForm(ses.Ses+"/v1/send/mail", formData)

	if err != nil {
		return err
	}

	if http.StatusNoContent == resp.StatusCode {
		return nil
	}

	bytes, err := ioutil.ReadAll(resp.Body)

	if err == nil {
		err = errors.New(string(bytes))
	}

	return err
}
