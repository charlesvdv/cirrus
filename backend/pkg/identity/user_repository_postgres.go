package identity

import (
	"github.com/charlesvdv/cirrus/backend/db"
	"github.com/jackc/pgx/v4/pgxpool"
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
		return User{}, db.ConvertPostgresError(err)
	}
	return user, nil
}
