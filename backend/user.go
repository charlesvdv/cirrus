package cirrus

import "time"

// UserID is the id of the user
type UserID string

func (id UserID) String() string {
	return string(id)
}

// User is a cirrus user
type User struct {
	ID        UserID
	CreatedAt time.Time
}
