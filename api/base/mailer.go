package kpbatApi

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"net/smtp"
)

func getAddress() string {
	var config = GetConfig()
	return fmt.Sprintf("%s:%d", config.Mail.Hostname, config.Mail.Port)
}
func getTlSConfig() *tls.Config {
	return &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         GetConfig().Mail.Hostname,
	}
}

func getMailClient() (*smtp.Client, error) {
	var config = GetConfig()
	conn, err := tls.Dial("tcp", getAddress(), getTlSConfig())
	if err != nil {
		return nil, err
	}
	client, err := smtp.NewClient(conn, config.Mail.Hostname)
	if err != nil {
		return nil, err
	}
	return client, err
}
func authorizeClient(client *smtp.Client) error {
	var config = GetConfig()
	auth := smtp.PlainAuth("", config.Mail.User, config.Mail.Password, config.Mail.Hostname)
	return client.Auth(auth)
}
func parseTemplate(templateName string, data any) (string, error) {
	t, _ := template.ParseFiles("templates/" + templateName)
	var body bytes.Buffer
	if err := t.Execute(&body, data); err != nil {
		return "", err
	}
	return body.String(), nil
}

func SendMail(to string, subject string, templateName string, data any) error {
	var config = GetConfig()
	from := config.Mail.User

	client, err := getMailClient()
	if err := authorizeClient(client); err != nil {
		return err
	}

	message, err := parseTemplate(templateName, data)
	if err != nil {
		return err
	}

	msg := "From: " + from + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-version: 1.0;\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\r\n\r\n" +
		message

	if err := client.Mail(from); err != nil {
		return err
	}

	if err := client.Rcpt(to); err != nil {
		return err
	}

	w, err := client.Data()
	if err != nil {
		return err
	}
	if _, err := w.Write([]byte(msg)); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	defer client.Quit()
	return nil
}
