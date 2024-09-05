package models

import "time"

type InvoiceDescription struct {
	Id       int64     `json:"id"`
	Date     time.Time `json:"date"`
	DueDate  time.Time `json:"due_date"`
	TotalDue float64   `json:"total_due"`
}
