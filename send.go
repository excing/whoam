package main

import (
	"errors"
	"flag"
	"io/ioutil"
	"net/http"
	"net/url"
)

var sendMailServerDomain string

func init() {
	flag.StringVar(&sendMailServerDomain, "mail", "http://10.0.0.254:19825:", "Send mail server domain url")
}

// SendMail 发送短信
func SendMail(to string, subject string, body string) error {
	formData := url.Values{
		"email":   {to},
		"subject": {subject},
		"body":    {body},
	}

	resp, err := http.PostForm(sendMailServerDomain+"/v1/send/mail", formData)

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
