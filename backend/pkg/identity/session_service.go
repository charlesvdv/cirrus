package identity

import (
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"

	"github.com/charlesvdv/cirrus/backend/db"
)

type SessionRepository interface {
	CreateAuthClient(ctx context.Context, tx db.Tx, client authClient) (authClient, error)
	GetAuthClient(ctx context.Context, tx db.Tx, userID UserID, clientReference string) (authClient, error)
	DeleteTokensWithClientID(ctx context.Context, tx db.Tx, clientID uint64) error
	CreateAccessToken(ctx context.Context, tx db.Tx, clientID uint64, token Token) (Token, error)
	CreateRefreshToken(ctx context.Context, tx db.Tx, clientID uint64, token Token) (Token, error)
	GetUserIDFromAccessToken(ctx context.Context, tx db.Tx, token string) (UserID, Token, error)
}

type SessionUserFetcher interface {
	GetUserWithEmail(ctx context.Context, tx db.Tx, email email) (User, error)
}

func NewSessionService(txProvider db.TxProvider, repository SessionRepository, userFetcher SessionUserFetcher) SessionService {
	return SessionService{
		txProvider:  txProvider,
		repository:  repository,
		userFetcher: userFetcher,
	}
}

type SessionService struct {
	txProvider  db.TxProvider
	repository  SessionRepository
	userFetcher SessionUserFetcher
}

type AuthenticationCredential struct {
	ClientReference string
	Email           string
	Password        string
}

type AuthenticationTokens struct {
	RefreshToken    Token
	AccessToken     Token
	ClientReference string
}

func (s *SessionService) CheckBearerToken(ctx context.Context, tokenValue string) (UserID, error) {
	var userID UserID
	var token Token
	err := s.txProvider.WithTransaction(ctx, func(tx db.Tx) error {
		var err error
		userID, token, err = s.repository.GetUserIDFromAccessToken(ctx, tx, tokenValue)
		return err
	})
	if err != nil {
		if db.IsErrNoRows(err) {
			return userID, ErrUnauthorized
		}
		log.Ctx(ctx).Warn().Err(err).Msg("Check bearer token failed")
		return userID, ErrInternal
	}

	if token.IsExpired() {
		log.Ctx(ctx).Warn().Msg("Usage of expired access token")
		return userID, ErrUnauthorized
	}

	return userID, nil
}

func (s *SessionService) Authenticate(ctx context.Context, credential AuthenticationCredential) (AuthenticationTokens, error) {
	var tokens AuthenticationTokens
	err := s.txProvider.WithTransaction(ctx, func(tx db.Tx) error {
		authenticatedUser, err := s.getAuthenticatedUser(ctx, tx, credential.Email, credential.Password)
		if err != nil {
			return err
		}

		authClient, err := s.getAuthClient(ctx, tx, authenticatedUser.ID(), credential.ClientReference)
		if err != nil {
			return err
		}

		err = s.invalidateTokens(ctx, tx, authClient)
		if err != nil {
			return err
		}

		tokens, err = s.createAuthenticationTokens(ctx, tx, authClient)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		if errors.Is(err, ErrInvalidUsernameOrPassword) {
			log.Ctx(ctx).Debug().Err(err).Msg("Invalid user name or password")
			return tokens, ErrInvalidUsernameOrPassword
		}
		log.Ctx(ctx).Warn().Err(err).Msg("Authentication failed")
		return tokens, ErrInternal
	}

	return tokens, nil
}

func (s *SessionService) getAuthenticatedUser(ctx context.Context, tx db.Tx, rawEmail, plaintextPassword string) (User, error) {
	email, err := newEmail(rawEmail)
	if err != nil {
		return User{}, fmt.Errorf("%w: %v", ErrInvalidUsernameOrPassword, err)
	}

	user, err := s.userFetcher.GetUserWithEmail(ctx, tx, email)
	if err != nil {
		return User{}, fmt.Errorf("%w: %v", ErrInvalidUsernameOrPassword, err)
	}

	err = user.verifyPassword(plaintextPassword)
	if err != nil {
		return User{}, fmt.Errorf("%w: %v", ErrInvalidUsernameOrPassword, err)
	}

	return user, nil
}

func (s *SessionService) getAuthClient(ctx context.Context, tx db.Tx, userID UserID, clientReference string) (authClient, error) {
	if clientReference == "" {
		log.Ctx(ctx).Debug().Msg("No client reference.. Creating a new one")
		// TODO: the security could be increased here by controlling a bit the client creation (sending email before accepting the
		// login, ...)
		client, err := generateNewAuthClient(userID)
		if err != nil {
			return client, err
		}

		client, err = s.repository.CreateAuthClient(ctx, tx, client)
		if err != nil {
			return client, err
		}

		return client, nil
	}

	client, err := s.repository.GetAuthClient(ctx, tx, userID, clientReference)
	if err != nil {
		return client, err
	}

	return client, err
}

func (s *SessionService) invalidateTokens(ctx context.Context, tx db.Tx, client authClient) error {
	return s.repository.DeleteTokensWithClientID(ctx, tx, client.id)
}

func (s *SessionService) createAuthenticationTokens(ctx context.Context, tx db.Tx, client authClient) (AuthenticationTokens, error) {
	refreshToken, err := generateRefreshToken()
	if err != nil {
		return AuthenticationTokens{}, err
	}
	accessToken, err := generateAccessToken()
	if err != nil {
		return AuthenticationTokens{}, err
	}

	tokens := AuthenticationTokens{
		ClientReference: client.clientReference,
	}
	tokens.RefreshToken, err = s.repository.CreateRefreshToken(ctx, tx, client.id, refreshToken)
	if err != nil {
		return tokens, err
	}

	tokens.AccessToken, err = s.repository.CreateAccessToken(ctx, tx, client.id, accessToken)
	if err != nil {
		return tokens, err
	}

	return tokens, nil
}
