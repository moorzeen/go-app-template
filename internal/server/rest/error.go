package rest

import (
	"database/sql"
	"github.com/rs/zerolog/log"
	"go-app-template/internal/service/auth"
	"go-app-template/internal/service/user"
	"net/http"
	"strings"
)

func errStatusWithJSON(err error) (int, any) {
	switch {

	case strings.Contains(err.Error(), "unknown field"):
		log.Warn().Err(err).Send()
		return errResponse(http.StatusBadRequest, "unknown fields in request")

	case strings.Contains(err.Error(), user.ErrEmptyCredentials.Error()):
		log.Warn().Err(err).Send()
		return errResponse(http.StatusBadRequest, user.ErrEmptyCredentials.Error())

	case strings.Contains(err.Error(), user.ErrWeakPassword.Error()):
		log.Warn().Err(err).Send()
		return errResponse(http.StatusBadRequest, user.ErrWeakPassword.Error())

	case strings.Contains(err.Error(), user.ErrLoginTaken.Error()):
		log.Warn().Err(err).Send()
		return errResponse(http.StatusConflict, user.ErrLoginTaken.Error())

	case strings.Contains(err.Error(), auth.ErrEmptyCredentials.Error()):
		log.Warn().Err(err).Send()
		return errResponse(http.StatusBadRequest, auth.ErrEmptyCredentials.Error())

	case strings.Contains(err.Error(), auth.ErrWrongCredentials.Error()) ||
		strings.Contains(err.Error(), sql.ErrNoRows.Error()):
		log.Warn().Err(err).Send()
		return errResponse(http.StatusBadRequest, auth.ErrWrongCredentials.Error())

	default:
		log.Error().Err(err).Send()
		return errResponse(http.StatusInternalServerError, err.Error())
	}
}

func errResponse(code int, message string) (int, any) {
	return code, &metaResponse{
		StatusCode: code,
		Message:    message,
	}
}
