package common

import (
	"time"
)

func GetSemaforizationColorFromDate(expiresAt time.Time) string {
	if expiresAt.Unix() < time.Now().AddDate(0, 6, 0).Unix() {
		return "red"
	} else if expiresAt.Unix() >= time.Now().AddDate(0, 6, 0).Unix() && expiresAt.Unix() <= time.Now().AddDate(0, 12, 0).Unix() {
		return "ambar"
	} else {
		return "green"
	}
}
