package repository

import models "simple-grpc-go/user"

type UserRepository interface {
	Fetch(cursor string, num int64) ([]*models.User, error)
	GetByID(id int64) (*models.User, error)
	Update(user *models.User) (*models.User, error)
	Store(u *models.User) (int64, error)
	Delete(id int64) (bool, error)
}
