package domain

import (
	"context"
	"time"
)

type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Age       int64     `json:"age"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserUsecase interface {
	Fetch(ctx context.Context, cursor string, num int64) ([]User, error)
	GetByID(ctx context.Context, id int64) (User, error)
	Update(ctx context.Context, ar *User) error
	Store(context.Context, *User) error
	Delete(ctx context.Context, id int64) error
}

type UserRepository interface {
	Fetch(ctx context.Context, cursor string, num int64) (res []User, nextCursor string, err error)
	GetByID(ctx context.Context, id int64) (User, error)
	Update(ctx context.Context, ar *User) error
	Store(ctx context.Context, a *User) error
	Delete(ctx context.Context, id int64) error
}
