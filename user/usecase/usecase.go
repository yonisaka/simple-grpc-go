package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"simple-grpc-go/config"
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

var expiredCache = 1 * time.Minute

func getCache(ctx context.Context, key string) ([]byte, error) {
	get := config.RedisClient.Get(ctx, key)
	err := get.Err()
	if err != nil {
		return nil, err
	}

	res, err := get.Result()
	if err != nil {
		return nil, err
	}

	return []byte(res), nil
}

func (a *userUsecase) Fetch(cursor string, num int64) ([]*models.User, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if num == 0 {
		num = 10
	}
	nextCursor := ""
	var result []*models.User
	key := "users_"+fmt.Sprintf("%d", num)+cursor
	get, err := getCache(ctx, key)
	if err != nil {
		log.Println(err)
	}

	err = json.Unmarshal(get, &result)
	if err != nil {
		log.Println(err)
	}

	if size := len(result); size == int(num) {
		lastId := result[num-1].ID
		nextCursor = strconv.Itoa(int(lastId))
	}

	if get != nil {
		return result, nextCursor, nil
	}
	
	result, err = a.userRepos.Fetch(cursor, num)
	if err != nil {
		return nil, "", err
	}

	data, err := json.Marshal(result)
	if err != nil {
		return nil, "", err
	}
	set := config.RedisClient.Set(ctx, key, string(data), expiredCache).Err()
	if set != nil {
		return nil, "", set
	}
	
	if size := len(result); size == int(num) {
		lastId := result[num-1].ID
		nextCursor = strconv.Itoa(int(lastId))
	}

	return result, nextCursor, nil
}

func (a *userUsecase) GetByID(id int64) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	var result *models.User
	key := "user_"+fmt.Sprintf("%d", id)
	get, err := getCache(ctx, key)
	if err != nil {
		log.Println(err)
	}

	err = json.Unmarshal(get, &result)
	if err != nil {
		log.Println(err)
	}

	if get != nil {
		return result, nil
	}

	result, err = a.userRepos.GetByID(id)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	set := config.RedisClient.Set(ctx, key, string(data), expiredCache).Err()
	if set != nil {
		return nil, set
	}

	return result, nil
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
