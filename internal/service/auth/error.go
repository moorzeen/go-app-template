package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrEmptyCredentials = errors.New("login and password must be filled in")
	ErrWrongCredentials = errors.New("wrong login or password")
	ErrUpdateLoginTime  = errors.New("failed to update last login time")
)

func errToResponse(err error) error {
	switch {

	case errors.Is(err, ErrEmptyCredentials):
		log.Warn().Err(err).Send()
		return status.Error(codes.InvalidArgument, ErrEmptyCredentials.Error())

	case errors.Is(err, sql.ErrNoRows):
		log.Warn().Err(fmt.Errorf(ErrWrongCredentials.Error()+": %w", err)).Send()
		return status.Error(codes.Unauthenticated, ErrWrongCredentials.Error())

	case errors.Is(err, ErrWrongCredentials):
		log.Warn().Err(fmt.Errorf(ErrWrongCredentials.Error()+": %w", err)).Send()
		return status.Error(codes.Unauthenticated, ErrWrongCredentials.Error())

	default:
		log.Error().Err(err).Send()
		return status.Error(codes.Unknown, err.Error())

	}
}
