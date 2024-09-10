package emailsender

import (
	"bytes"
	"context"
	"fmt"
	"sync"

	"net/smtp"

	"github.com/Khaym03/invoces-service/internal/components"
	"github.com/Khaym03/invoces-service/internal/models"
)

var once sync.Once
var intance *emailsender

type emailsender struct{}

func Service() *emailsender {
	once.Do(func() {
		intance = &emailsender{}
	})
	return intance
}

func (es *emailsender) Send(customer models.CustomerDetails, pdfURL string) error {
	hostname := "smtp.gmail.com"
	laMediaDigitalEmail := "daniellarosa20003@gmail.com"
	password := "dxjbrigxfosjhnpy"

	auth := smtp.PlainAuth("", laMediaDigitalEmail, password, hostname)

	var buf bytes.Buffer

	comp := components.Header("pepe", pdfURL)
	comp.Render(context.Background(), &buf)

	fmt.Println(buf.String())

	formatt := fmt.Sprintf("To: %s\r\n", customer.Email)

	to := []string{customer.Email}
	msg := []byte(formatt +
		"Subject: Factura\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n\r\n" + buf.String())

	err := smtp.SendMail(hostname+":587", auth, laMediaDigitalEmail, to, msg)
	if err != nil {
		panic(err)
	}

	return nil
}
