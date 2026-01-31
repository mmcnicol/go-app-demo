package models

import "time"

type RecentGopherItem struct {
	GopherId         string
	Name             string
	DateOfBirth      time.Time
	DateLastAccessed time.Time
}
