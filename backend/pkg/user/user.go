package user

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/goware/emailx"
	"github.com/muesli/crunchy"
)

type ID = uint64

type hashedPassword []byte

func newHashedPassword(password string) (hashedPassword, error) {
	if len(password) <= 8 {
		return hashedPassword{}, errors.New("Password too short")
	}
	if len(password) > 64 {
		return hashedPassword{}, errors.New("Password is too long")
	}
	validator := crunchy.NewValidatorWithOpts(crunchy.Options{
		CheckHIBP: false,
	})
	err := validator.Check(password)
	if err != nil {
		return hashedPassword{}, err
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return hashedPassword{}, err
	}

	return hashedPwd, nil
}

func (p hashedPassword) compareWithPlainText(plaintextPassword string) error {
	return bcrypt.CompareHashAndPassword(p, []byte(plaintextPassword))
}

type email string

func newEmail(rawEmail string) (email, error) {
	normalizedEmail := emailx.Normalize(rawEmail)
	return email(normalizedEmail), emailx.ValidateFast(normalizedEmail)
}

func newSignupUser(rawEmail, rawPassword string) (User, error) {
	email, err := newEmail(rawEmail)
	if err != nil {
		return User{}, fmt.Errorf("Invalid email '%s': %w", rawEmail, err)
	}

	hashedPassword, err := newHashedPassword(rawPassword)
	if err != nil {
		return User{}, errors.New("Invalid password")
	}

	return User{
		id:             0,
		email:          email,
		hashedPassword: hashedPassword,
	}, nil
}

type User struct {
	id             ID
	email          email
	hashedPassword hashedPassword
}

func (u *User) ID() ID {
	return u.id
}

func (u *User) Email() string {
	return string(u.email)
}

func (u *User) verifyPassword(plaintextPwd string) error {
	err := u.hashedPassword.compareWithPlainText(plaintextPwd)
	if err != nil {
		return errors.New("Password doesn't match")
	}
	return nil
}
