package http

import "time"

func FormatTime(t time.Time) string {
	return t.Format(time.RFC3339)
}
