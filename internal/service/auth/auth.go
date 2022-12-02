package auth

import (
	"context"
	"crypto/hmac"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"go-app-template/internal/service/auth/model"
	pbAuth "go-app-template/internal/service/auth/proto"
	"go-app-template/internal/storage"
	"go-app-template/pkg"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"strings"

	"time"
)

type Service struct {
	storage storage.Storage
	redis   *redis.Client
	key     string
}

func NewAuthService(st storage.Storage, rc *redis.Client, key string) (*Service, error) {
	return &Service{
		storage: st,
		redis:   rc,
		key:     key,
	}, nil
}

func (s Service) Authorize(ctx context.Context, credentials *pbAuth.Credentials) (*pbAuth.Tokens, error) {
	if credentials.Login == "" || credentials.Password == "" {
		return nil, errToResponse(ErrEmptyCredentials)
	}

	credentials.Login = strings.ToLower(credentials.Login)

	user, err := s.storage.GetUserByLogin(ctx, credentials.Login)
	if err != nil {
		return nil, errToResponse(err)
	}

	passHash := pkg.GenerateHash(credentials.Password, s.key)
	if !hmac.Equal(passHash, user.Password) {
		return nil, errToResponse(ErrWrongCredentials)
	}

	// access token
	accessClaims := model.CustomClaims{
		UserID: user.ID,
		Login:  user.Login,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(s.key))
	if err != nil {
		return nil, errToResponse(err)
	}

	// refresh token
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshClaims := refreshToken.Claims.(jwt.MapClaims)
	refreshClaims["refreshTokenId"] = uuid.New()
	refreshClaims["userId"] = user.ID
	refreshClaims["fingerprints"] = credentials.Fingerprints
	refreshTokenString, err := refreshToken.SignedString([]byte(s.key))
	if err != nil {
		return nil, errToResponse(err)
	}

	key := fmt.Sprintf("%s:%s", refreshClaims["userId"], refreshClaims["fingerprints"])

	err = s.redis.Get(ctx, key).Err()
	if err != redis.Nil {
		err = s.redis.Del(ctx, key).Err()
		if err != nil {
			return nil, errToResponse(err)
		}
	}
	err = s.redis.Set(ctx, key, nil, time.Hour*168).Err()
	if err != nil {
		return nil, errToResponse(err)
	}

	err = s.storage.UpdateLoginTime(ctx, user.ID, time.Now().UTC())
	if err != nil {
		err = fmt.Errorf(ErrUpdateLoginTime.Error()+": %w", err)
		log.Error().Err(err).Send()
	}

	return &pbAuth.Tokens{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}, nil
}

func (s Service) Refresh(ctx context.Context, token *pbAuth.RefreshToken) (*pbAuth.Tokens, error) {
	return nil, status.Error(codes.Unimplemented, "unimplemented method")
}

func (s Service) Validate(ctx context.Context, token *pbAuth.AccessToken) (*emptypb.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "unimplemented method")
}

func (s Service) Logout(ctx context.Context, token *pbAuth.RefreshToken) (*emptypb.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "unimplemented method")
}
