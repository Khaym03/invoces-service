package emailsender

import (
	"bytes"
	"context"
	"fmt"
	"net/smtp"

	"github.com/Khaym03/invoces-service/internal/components"
)

// var (
// 	from       = "gopher@example.net"
// 	msg        = []byte("dummy message")
// 	recipients = []string{"foo@example.com"}
// )

type emailsender struct{}

func Service() *emailsender {
	return &emailsender{}
}

func (es *emailsender) Send() {
	hostname := "smtp.gmail.com"
	user := "daniellarosa20003@gmail.com"
	password := "dxjbrigxfosjhnpy"

	auth := smtp.PlainAuth("", user, password, hostname)

	var buf bytes.Buffer

	comp := components.Header("pepe")
	comp.Render(context.Background(), &buf)

	fmt.Println(buf.String())

	to := []string{user}
	msg := []byte("To: daniellarosa20003@gmail.com\r\n" +
		"Subject: Account Status\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n\r\n" + buf.String())

	err := smtp.SendMail(hostname+":587", auth, user, to, msg)
	if err != nil {
		panic(err)
	}
}
