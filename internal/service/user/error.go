package user

import (
	"errors"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
)

var (
	ErrEmptyCredentials = errors.New("login and password must be filled in")
	ErrWeakPassword     = errors.New("password required at least 8 characters")
	ErrLoginTaken       = errors.New("login is already taken")
)

func errToResponse(err error) error {
	switch {
	case errors.Is(err, ErrEmptyCredentials):
		log.Warn().Err(err).Send()
		return status.Error(codes.InvalidArgument, ErrEmptyCredentials.Error())
	case errors.Is(err, ErrWeakPassword):
		log.Warn().Err(err).Send()
		return status.Error(codes.InvalidArgument, err.Error())
	case strings.Contains(err.Error(), "violates unique constraint \"users_login_key\""):
		log.Warn().Err(err).Send()
		return status.Error(codes.AlreadyExists, ErrLoginTaken.Error())
	default:
		log.Error().Err(err).Send()
		return status.Error(codes.Internal, err.Error())
	}
}
