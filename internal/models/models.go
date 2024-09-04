package models

import "time"

type InvoiceDescription struct {
	Id       int64
	Date     time.Time
	DueDate  time.Time
	TotalDue float64
}
