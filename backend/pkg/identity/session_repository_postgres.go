package identity

import (
	"github.com/charlesvdv/cirrus/backend/db"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
)

type SessionRepositoryPostgres struct{}

func (r *SessionRepositoryPostgres) CreateAuthClient(ctx context.Context, tx db.Tx, client authClient) (authClient, error) {
	err := tx.(db.PostgresTx).QueryRow(ctx, `
		INSERT INTO identity.auth_client (user_id, client_reference)
		VALUES ($1, $2)
		RETURNING auth_client_id
	`, client.userID, client.clientReference).Scan(&client.id)

	return client, err
}

func (r *SessionRepositoryPostgres) GetAuthClient(ctx context.Context, tx db.Tx, userID UserID, clientReference string) (authClient, error) {
	var client authClient
	err := tx.(db.PostgresTx).QueryRow(ctx, `
		SELECT auth_client_id, client_reference, user_id
		FROM identity.auth_client
		WHERE user_id = $1 AND client_reference = $2
	`, userID, clientReference).Scan(&client.id, &client.clientReference, &client.userID)

	return client, err
}

func (r *SessionRepositoryPostgres) DeleteTokensWithClientID(ctx context.Context, tx db.Tx, clientID uint64) error {
	commandTag, err := tx.(db.PostgresTx).Exec(ctx, `
		DELETE FROM identity.client_token
		WHERE auth_client_id = $1
	`, clientID)
	if err != nil {
		return err
	}
	log.Ctx(ctx).Debug().Int64("rows affected", commandTag.RowsAffected()).Uint64("client id", clientID).Msg("Delete tokens")
	return nil
}

func (r *SessionRepositoryPostgres) CreateAccessToken(ctx context.Context, tx db.Tx, clientID uint64, token Token) (Token, error) {
	return r.createToken(ctx, tx, clientID, token, "access")
}

func (r *SessionRepositoryPostgres) CreateRefreshToken(ctx context.Context, tx db.Tx, clientID uint64, token Token) (Token, error) {
	return r.createToken(ctx, tx, clientID, token, "refresh")
}

func (r *SessionRepositoryPostgres) GetUserIDFromAccessToken(ctx context.Context, tx db.Tx, tokenValue string) (UserID, Token, error) {
	var userID UserID
	var token Token
	err := tx.(db.PostgresTx).QueryRow(ctx, `
		SELECT c.user_id, t.value, t.expired_at
		FROM identity.auth_client AS c
		INNER JOIN identity.client_token AS t ON (c.auth_client_id = t.auth_client_id)
		WHERE t.value = $1 AND t.type = 'access'
	`, tokenValue).Scan(&userID, &token.value, &token.expiredAt)

	return userID, token, err
}

func (r *SessionRepositoryPostgres) createToken(ctx context.Context, tx db.Tx, clientID uint64, token Token, clientType string) (Token, error) {
	_, err := tx.(db.PostgresTx).Exec(ctx, `
		INSERT INTO identity.client_token (auth_client_id, value, type, expired_at)
		VALUES ($1, $2, $3, $4)
	`, clientID, token.value, clientType, token.expiredAt)

	return token, err
}
