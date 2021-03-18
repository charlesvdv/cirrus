package sqlite

import (
	"fmt"
	"strconv"
	"time"

	cirrus "github.com/charlesvdv/cirrus/backend"
	"github.com/rs/zerolog/log"
)

func formatTime(t time.Time) string {
	return t.Format(time.RFC3339)
}

func getTime(text string) time.Time {
	t, err := time.Parse(time.RFC3339, text)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get time from database")
	}
	return t
}

func formatUserID(id cirrus.UserID) int64 {
	return formatID(id)
}

func getUserID(id int64) cirrus.UserID {
	return cirrus.UserID(getID(id))
}

func formatObjectID(id cirrus.ObjectID) int64 {
	return formatID(id)
}

func getObjectID(id int64) cirrus.ObjectID {
	return cirrus.ObjectID(getID(id))
}

func formatID(id fmt.Stringer) int64 {
	rowid, err := strconv.ParseInt(id.String(), 10, 64)
	if err != nil {
		log.Fatal().Err(err).Msg("invalid user id")
	}
	return rowid
}

func getID(id int64) string {
	return strconv.FormatInt(id, 10)
}
