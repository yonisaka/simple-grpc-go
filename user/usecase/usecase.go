package usecase

import (
	models "simple-grpc-go/user"
	"simple-grpc-go/user/repository"
	"strconv"
	"time"
)

type UserUsecase interface {
	Fetch(cursor string, num int64) ([]*models.User, string, error)
	GetByID(id int64) (*models.User, error)
	Update(us *models.User) (*models.User, error)
	Store(*models.User) (*models.User, error)
	Delete(id int64) (bool, error)
}

type userUsecase struct {
	userRepos repository.UserRepository
}

func (a *userUsecase) Fetch(cursor string, num int64) ([]*models.User, string, error) {
	if num == 0 {
		num = 10
	}

	listUser, err := a.userRepos.Fetch(cursor, num)
	if err != nil {
		return nil, "", err
	}
	nextCursor := ""

	if size := len(listUser); size == int(num) {
		lastId := listUser[num-1].ID
		nextCursor = strconv.Itoa(int(lastId))
	}

	return listUser, nextCursor, nil
}

func (a *userUsecase) GetByID(id int64) (*models.User, error) {
	return a.userRepos.GetByID(id)
}

func (a *userUsecase) Store(m *models.User) (*models.User, error) {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	id, err := a.userRepos.Store(m)
	if err != nil {
		return nil, err
	}

	m.ID = id
	return m, nil
}

func (a *userUsecase) Update(m *models.User) (*models.User, error) {
	_, err := a.userRepos.GetByID(m.ID)
	if err != nil {
		return nil, err
	}

	m.UpdatedAt = time.Now()
	return a.userRepos.Update(m)
}

func (a *userUsecase) Delete(id int64) (bool, error) {
	existedUser, _ := a.GetByID(id)

	if existedUser == nil {
		return false, models.NOT_FOUND_ERROR
	}

	return a.userRepos.Delete(id)
}

func NewUserUsecase(u repository.UserRepository) UserUsecase {
	return &userUsecase{u}
}
