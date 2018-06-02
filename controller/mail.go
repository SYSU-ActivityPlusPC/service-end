package controller

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-gomail/gomail"
)

// SendMail send mail according to the content
func SendMail(send string, to string, content string, sub string) error {
	// Get receivers and password
	toArray := strings.Split(to, ";")
	password := os.Getenv("ADMIN_MAIL_PASS")
	// Set email content
	m := gomail.NewMessage()
	m.SetHeader("From", send)
	m.SetHeader("To", toArray...)
	m.SetHeader("Subject", sub)
	m.SetBody("text/html", content)

	// Dial and send email
	d := gomail.NewDialer("smtp.exmail.qq.com", 25, send, password)
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("成功发送邮件")
	return nil
}