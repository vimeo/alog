package alog

import "time"

// Entry is the struct passed to user-supplied formatters.
type Entry struct {
	Time  time.Time
	Tags  [][2]string
	STags []STag
	File  string
	Line  int
	Msg   string
}
