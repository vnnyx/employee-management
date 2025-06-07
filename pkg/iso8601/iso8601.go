package iso8601

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"time"
)

const (
	dayNs   = 24 * time.Hour
	weekNs  = dayNs * 7
	yearNs  = dayNs * 365
	monthNs = yearNs / 12
)

var pattern = regexp.MustCompile(`^P((?P<year>\d+)Y)?((?P<month>\d+)M)?((?P<week>\d+)W)?((?P<day>\d+)D)?(T((?P<hour>\d+)H)?((?P<minute>\d+)M)?((?P<second>\d+)S)?)?$`)

func MustParse(from string) time.Duration {
	duration, err := Parse(from)
	if err != nil {
		panic(fmt.Sprintf("failed to parse duration string %q: %v", from, err))
	}
	return duration
}

func Parse(from string) (time.Duration, error) {
	if !pattern.MatchString(from) {
		return 0, errors.New("could not parse duration string")
	}

	match := pattern.FindStringSubmatch(from)
	var duration time.Duration

	for i, name := range pattern.SubexpNames() {
		if i == 0 || name == "" || match[i] == "" {
			continue
		}

		val, err := strconv.ParseInt(match[i], 10, 64)
		if err != nil {
			return 0, err
		}

		switch name {
		case "year":
			duration += time.Duration(val) * yearNs
		case "month":
			duration += time.Duration(val) * monthNs
		case "week":
			duration += time.Duration(val) * weekNs
		case "day":
			duration += time.Duration(val) * dayNs
		case "hour":
			duration += time.Duration(val) * time.Hour
		case "minute":
			duration += time.Duration(val) * time.Minute
		case "second":
			duration += time.Duration(val) * time.Second
		default:
			return 0, fmt.Errorf("unknown field %s", name)
		}
	}

	return duration, nil
}

func ToString(from time.Duration) string {
	if from == 0 {
		return "PT0S"
	}

	var (
		negative bool
		result   = "P"
		hasTime  bool
	)

	if from < 0 {
		from = -from
		negative = true
	}

	appendPart := func(value float64, unit string, isTime bool) {
		if !hasTime && isTime {
			result += "T"
			hasTime = true
		}
		result += strconv.FormatFloat(value, 'f', -1, 64) + unit
	}

	units := []struct {
		duration time.Duration
		unit     string
		isTime   bool
	}{
		{yearNs, "Y", false},
		{monthNs, "M", false},
		{weekNs, "W", false},
		{dayNs, "D", false},
		{time.Hour, "H", true},
		{time.Minute, "M", true},
		{time.Second, "S", true},
	}

	for _, u := range units {
		if from >= u.duration {
			value := math.Floor(from.Seconds() / u.duration.Seconds())
			appendPart(value, u.unit, u.isTime)
			from -= time.Duration(value) * u.duration
		}
	}

	if negative {
		return "-" + result
	}
	return result
}
