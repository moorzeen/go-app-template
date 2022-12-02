package model

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID          uuid.UUID
	Login       string
	Password    []byte
	Active      bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	LastLoginAt time.Time
}
