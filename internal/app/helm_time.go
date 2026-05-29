package app

import (
	"fmt"
	"strings"
	"time"
)

const (
	ctLayout  = "2006-01-02 15:04:05.000000000 -0700 MST"
	ctLayout2 = "2006-01-02 15:04:05.000000000 -0700 -0700"
)

var nilTime = (time.Time{}).UnixNano()

type HelmTime struct {
	time.Time
}

func (ht *HelmTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		ht.Time = time.Time{}
		return
	}
	// we need to split the time into parts and make sure milliseconds len = 6, it happens to skip trailing zeros
	updatedFields := strings.Fields(s)
	updatedHour := strings.Split(updatedFields[1], ".")
	milliseconds := updatedHour[1]
	for i := len(milliseconds); i < 9; i++ {
		milliseconds = fmt.Sprintf("%s0", milliseconds)
	}
	s = fmt.Sprintf("%s %s.%s %s %s", updatedFields[0], updatedHour[0], milliseconds, updatedFields[2], updatedFields[3])
	ht.Time, err = time.Parse(ctLayout, s)
	if err != nil {
		ht.Time, err = time.Parse(ctLayout2, s)
	}
	return
}

func (ht *HelmTime) MarshalJSON() ([]byte, error) {
	if ht.Time.UnixNano() == nilTime {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", ht.Time.Format(ctLayout))), nil
}

func (ht *HelmTime) IsSet() bool {
	return ht.UnixNano() != nilTime
}
