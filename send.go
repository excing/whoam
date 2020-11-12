package main

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Send mail server domain
var sesServer string

func init() {
	Flag("ses", &sesServer, "Send mail server domain url")
}

// SendMail 发送短信
func SendMail(to string, subject string, body string) error {
	formData := url.Values{
		"to":      {to},
		"subject": {subject},
		"body":    {body},
	}

	resp, err := http.PostForm(sesServer+"/v1/send/mail", formData)

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
