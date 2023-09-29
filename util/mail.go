package util

import (
	"crypto/tls"
	"errors"
	"html/template"
	"strconv"
	"strings"

	gomail "gopkg.in/mail.v2"
)

func getBody(email string, code string) (string, error) {
	var firstname string = strings.Split(email, ".")[0]
	firstname = strings.ToUpper(string(firstname[0])) + string(firstname[1:])
	type Data struct {
		Firstname string
		Code      string
	}
	tmpl, err := template.ParseFiles(Get("TEMPLATE_PATH"))
	if err != nil {
		return "", errors.New("file :" + Get("TEMPLATE_PATH") + " not found")
	}

	var buf strings.Builder

	err = tmpl.Execute(&buf, Data{
		Firstname: firstname,
		Code:      code,
	})

	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func SendConfirmationMail(email string, code string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", Get("MAIL_SMTP"))
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Hack The Crous : confirmation d'inscription")

	body, err := getBody(email, code)

	if err != nil {
		return err
	}

	m.SetBody("text/html", body)

	port, err := strconv.Atoi(Get("SMTP_PORT"))

	if err != nil {
		return err
	}

	d := gomail.NewDialer(Get("SMTP_HOST"), port, Get("MAIL_SMTP"), Get("MAIL_PASSWORD"))

	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
