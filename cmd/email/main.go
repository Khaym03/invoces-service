package main

import "github.com/Khaym03/invoces-service/internal/core/service/emailsender"

func main() {
	srv := emailsender.Service()

	srv.Send()
}
