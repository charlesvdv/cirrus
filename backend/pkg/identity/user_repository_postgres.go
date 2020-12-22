package identity

import (
	"errors"

	"github.com/charlesvdv/cirrus/backend/db"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func (r *PostgresRepository) Create(ctx context.Context, tx db.Tx, user User) (User, error) {
	err := tx.(db.PostgresTx).QueryRow(ctx, `
		INSERT INTO identity.user
		(email, password)
		VALUES($1, $2)
		RETURNING user_id
	`, user.Email(), string(user.hashedPassword)).Scan(&user.id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.DuplicateColumn {
				return User{}, errors.New("User already exists")
			}
		}
		log.Debug().Err(err).Msg("Failed to create user")
		return User{}, errors.New("Failed to create user")
	}
	return user, nil
}
