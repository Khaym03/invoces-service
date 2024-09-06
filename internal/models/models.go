package models

import "time"

type InvoiceDetails struct {
	Id       int64     `json:"id"`
	Date     time.Time `json:"date"`
	DueDate  time.Time `json:"due_date"`
	TotalDue float64   `json:"total_due"`
}

type CustomerDetails struct {
	FullName    string `json:"full_name"`
	Email       string `json:"email"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
}

type InvoiceInput struct {
	InvoiceDetails  InvoiceDetails  `json:"invoice_details"`
	CustomerDetails CustomerDetails `json:"customer_details"`
}

func (i *InvoiceInput) Format() {}
