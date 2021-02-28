package sqlite

import (
	"strconv"
	"time"

	cirrus "github.com/charlesvdv/cirrus/backend"
	"github.com/rs/zerolog/log"
)

func formatTime(t time.Time) string {
	return t.Format(time.RFC3339)
}

func formatUserID(id cirrus.UserID) int64 {
	rowid, err := strconv.ParseInt(id.String(), 10, 64)
	if err != nil {
		log.Fatal().Err(err).Msg("invalid user id")
	}
	return rowid
}
