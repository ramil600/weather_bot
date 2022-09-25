package main

import (
	"fmt"
	"regexp"
	"strconv"
)

//regExp has to match for time from user input
const regExp = `^([0-1]?[0-9]|2[0-3]):[0-5][0-9]$`

func parseTime(time string) (int, error) {

	var validTime = regexp.MustCompile(regExp)

	if !validTime.MatchString(time) {
		return 0, fmt.Errorf("time selected is not millitary format")
	}
	timeStr := fmt.Sprintf("%s%s", time[0:2], time[3:])

	timeInt, _ := strconv.Atoi(timeStr)
	return timeInt, nil

}
