package formatter

import (
	"fmt"
	"time"
)

func DateWithoutTimestamp(t time.Time) string {
	return t.Format("2006-01-02")
}


func DateWithTimestamp (t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func DateStringToTime(date string) (time.Time, error){
	layout := "2006-01-02"
	birthday, errParse := time.Parse(layout, date)
	if errParse != nil {
		return time.Time{}, fmt.Errorf("error while parsing date string to time : %v", errParse)
	}

	return birthday, nil
}