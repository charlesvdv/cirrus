package identity

import (
	"github.com/charlesvdv/cirrus/backend/db"
	"golang.org/x/net/context"
)

type UserRepositoryPostgres struct{}

func (r *UserRepositoryPostgres) Create(ctx context.Context, tx db.Tx, user User) (User, error) {
	err := tx.(db.PostgresTx).QueryRow(ctx, `
		INSERT INTO identity.user (email, password)
		VALUES ($1, $2)
		RETURNING user_id
	`, user.Email(), string(user.hashedPassword)).Scan(&user.id)

	return user, err
}

func (r *UserRepositoryPostgres) GetUserWithEmail(ctx context.Context, tx db.Tx, email email) (User, error) {
	var user User
	err := tx.(db.PostgresTx).QueryRow(ctx, `
		SELECT user_id, email, password
		FROM identity.user
		WHERE email = $1
	`, email).Scan(&user.id, &user.email, &user.hashedPassword)

	return user, err
}

func (r *UserRepositoryPostgres) GetUserWithID(ctx context.Context, tx db.Tx, userID UserID) (User, error) {
	var user User
	err := tx.(db.PostgresTx).QueryRow(ctx, `
		SELECT user_id, email, password
		FROM identity.user
		WHERE user_id = $1
	`, userID).Scan(&user.id, &user.email, &user.hashedPassword)

	return user, err
}
