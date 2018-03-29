package alog

import "time"

// Entry is the struct passed to user-supplied formatters.
//
// The File and Line members are only populated if the WithFile() Option was
// passed to New.
type Entry struct {
	Time time.Time
	Tags [][2]string
	File string
	Line int
	Msg  string
}
