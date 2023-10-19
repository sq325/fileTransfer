package service

import (
	"time"

	str2duration "github.com/xhit/go-str2duration/v2"
)

type DateTimer interface {
	DateTime(duration, layout string) (string, error)
}

type datetimer struct{}

func NewDateTimer() DateTimer {
	return &datetimer{}
}

func (dt datetimer) DateTime(duration, layout string) (string, error) {
	d, err := str2duration.ParseDuration(duration)
	if err != nil {
		return "", err
	}
	now := time.Now()
	t := now.Add(d)
	return t.Format(layout), nil
}
