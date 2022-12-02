package user

import (
	"context"
	"github.com/segmentio/kafka-go"
	"go-app-template/internal/broker"
	"go-app-template/internal/service/user/model"
	pb "go-app-template/internal/service/user/proto"
	"go-app-template/internal/storage"
	"go-app-template/pkg"
	"strings"
	"time"
)

type Service struct {
	storage     storage.Storage
	kafkaWriter *kafka.Writer
	key         string
}

func NewUserService(st storage.Storage, kw *kafka.Writer, key string) (*Service, error) {
	return &Service{
		storage:     st,
		kafkaWriter: kw,
		key:         key,
	}, nil
}

func (s *Service) Register(ctx context.Context, credentials *pb.Credentials) (*pb.User, error) {
	if credentials.Login == "" || credentials.Password == "" {
		return nil, errToResponse(ErrEmptyCredentials)
	}

	credentials.Login = strings.ToLower(credentials.Login)

	if err := passComplexity(credentials.Password); err != nil {
		return nil, errToResponse(err)
	}

	passHash := pkg.GenerateHash(credentials.Password, s.key)

	user := &model.User{
		Login:       credentials.Login,
		Password:    passHash,
		Active:      false,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
		LastLoginAt: *new(time.Time),
	}

	insertedUser, err := s.storage.InsertUser(ctx, user)
	if err != nil {
		return nil, errToResponse(err)
	}

	response := &pb.User{
		Id:          insertedUser.ID.String(),
		Login:       insertedUser.Login,
		Active:      &insertedUser.Active,
		CreatedAt:   insertedUser.CreatedAt.Format(pkg.TimeLayoutISO),
		UpdatedAt:   insertedUser.UpdatedAt.Format(pkg.TimeLayoutISO),
		LastLoginAt: insertedUser.LastLoginAt.Format(pkg.TimeLayoutISO),
	}

	go broker.ProduceNewUser(context.Background(), response, s.kafkaWriter)

	return response, nil
}
