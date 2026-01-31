package models

import "time"

type LabResultItem struct {
	ID         string
	Subject    string
	ReportDate time.Time
}
