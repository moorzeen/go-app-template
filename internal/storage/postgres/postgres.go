package postgres

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"template/internal/service/user/model"
	"time"
)

type DB struct {
	*sqlx.DB
}

type Config struct {
	Host      string
	Port      string
	User      string
	Password  string
	Name      string
	Migration bool
}

func NewPGStorage(cfg *Config) (*DB, error) {
	url := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name,
	)

	db, err := sqlx.Connect("postgres", url)
	if err != nil {
		return nil, err
	}

	if cfg.Migration {
		err = Migration(url)
		if err != nil {
			return nil, err
		}
	}

	return &DB{
		db,
	}, nil
}

func (db *DB) InsertUser(ctx context.Context, user *model.User) (*model.User, error) {
	query := `insert into users (login, password_hash, created_at, updated_at, last_login_at)
				values ($1, $2, $3, $4, $5)
				returning *`

	row := db.QueryRowContext(ctx, query, user.Login, user.Password, user.CreatedAt, user.UpdatedAt, user.LastLoginAt)
	if err := row.Err(); err != nil {
		return nil, err
	}

	err := row.Scan(&user.ID, &user.Login, &user.Password, &user.Active, &user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (db *DB) GetUserByLogin(ctx context.Context, login string) (*model.User, error) {
	query := `select user_id, login, password_hash, active, created_at, updated_at, last_login_at
				from users
				where login = $1`

	row := db.QueryRowContext(ctx, query, login)
	if err := row.Err(); err != nil {
		return nil, err
	}

	user := &model.User{}

	err := row.Scan(&user.ID, &user.Login, &user.Password, &user.Active, &user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (db *DB) UpdateLoginTime(ctx context.Context, ID uuid.UUID, time time.Time) error {
	query := `update users set last_login_at = $2 where user_id = $1`

	_, err := db.ExecContext(ctx, query, ID, time)
	if err != nil {
		return err
	}

	return nil
}
