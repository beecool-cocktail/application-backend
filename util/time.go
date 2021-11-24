package util

import "time"

func GetFormatTime(myTime time.Time, location string) string {
	loc, _ := time.LoadLocation(location)
	return myTime.In(loc).Format("2006-01-02 15:04:05")
}