package storage

import (
	"github.com/google/uuid"
	"golang.org/x/net/context"
	"template/internal/service/user/model"
	"time"
)

type Storage interface {
	InsertUser(ctx context.Context, user *model.User) (*model.User, error)
	GetUserByLogin(ctx context.Context, login string) (*model.User, error)
	UpdateLoginTime(ctx context.Context, ID uuid.UUID, time time.Time) error
}
