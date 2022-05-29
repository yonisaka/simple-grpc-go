package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	models "simple-grpc-go/user"
	"time"

	"github.com/go-redis/redis/v8"
)

type mysqlUserRepository struct {
	Conn *sql.DB
}

func NewMysqlUserRepository(Conn *sql.DB) UserRepository {
	return &mysqlUserRepository{Conn}
}

var cache = redis.NewClient(&redis.Options{
    Addr: "localhost:6379",
})

var expiredCache = 1 * time.Second

func (m *mysqlUserRepository) Fetch(cursor string, num int64) ([]*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result := make([]*models.User, 0)
	
	key := "users_"+fmt.Sprintf("%d", num)+cursor
	get := cache.Get(ctx, key)
	err := get.Err()
	if err != nil {
		fmt.Println(err)
	} 

	cacheResult, err := get.Result()
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal([]byte(cacheResult), &result)
	if err != nil {
		fmt.Println(err)
	}

	if cacheResult != "" {
		return result, nil
	}

	query := `SELECT id,name,email,age,created_at,updated_at
				FROM users WHERE ID > ? LIMIT ?`

	rows, err := m.Conn.QueryContext(ctx, query, cursor, num)
	if err != nil {
		log.Fatal(err)
		return nil, models.INTERNAL_SERVER_ERROR
	}
	defer rows.Close()

	for rows.Next() {
		t := new(models.User)
		err = rows.Scan(
			&t.ID,
			&t.Name,
			&t.Email,
			&t.Age,
			&t.UpdatedAt,
			&t.CreatedAt,
		)

		if err != nil {
			log.Fatal(err)
			return nil, models.INTERNAL_SERVER_ERROR
		}
		result = append(result, t)
	}

	data, err := json.Marshal(result)
	if err != nil {
		log.Fatal(err)
	}
	set := cache.Set(ctx, key, string(data), expiredCache).Err()
	if set != nil {
		log.Fatal(set)
	}

	return result, nil
}

func (m *mysqlUserRepository) GetByID(id int64) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	u := &models.User{}

	key := "user_"+fmt.Sprintf("%d", id)
	get := cache.Get(ctx, key)
	err := get.Err()
	if err != nil {
		fmt.Println(err)
	} 

	cacheResult, err := get.Result()
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal([]byte(cacheResult), u)
	if err != nil {
		fmt.Println(err)
	}
	
	if cacheResult != "" {
		return u, nil
	}
	
	query := `SELECT id, name, email, age, created_at, updated_at
				FROM users WHERE ID = ?`

	rows, err := m.Conn.QueryContext(ctx, query, id)
	if err != nil {
		log.Fatal(err)
		return nil, models.INTERNAL_SERVER_ERROR
	}

	result := make([]*models.User, 0)
	for rows.Next() {
		t := new(models.User)
		err = rows.Scan(
			&t.ID,
			&t.Name,
			&t.Email,
			&t.Age,
			&t.UpdatedAt,
			&t.CreatedAt,
		)

		if err != nil {
			log.Fatal(err)
			return nil, models.INTERNAL_SERVER_ERROR
		}
		result = append(result, t)
	}
	
	if len(result) > 0 {
		u = result[0]
	} else {
		return nil, models.NOT_FOUND_ERROR
	}

	data, err := json.Marshal(u)
	if err != nil {
		log.Fatal(err)
	}
	set := cache.Set(ctx, key, string(data), expiredCache).Err()
	if set != nil {
		log.Fatal(set)
	}

	return u, nil
}

func (m *mysqlUserRepository) Store(u *models.User) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `INSERT users SET name=?, email=?, age=?, created_at=?, updated_at=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		log.Fatal(err)
		return 0, models.INTERNAL_SERVER_ERROR
	}
	res, err := stmt.Exec(u.Name, u.Email, u.Age, u.CreatedAt, u.UpdatedAt)
	if err != nil {
		log.Fatal(err)
		return 0, models.INTERNAL_SERVER_ERROR
	}
	return res.LastInsertId()
}

func (m *mysqlUserRepository) Update(ur *models.User) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `UPDATE users SET name=?, email=?, age=?, updated_at=? WHERE ID = ?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	res, err := stmt.Exec(ur.Name, ur.Email, ur.Age, ur.UpdatedAt, ur.ID)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	affect, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	if affect < 1 {
		return nil, errors.New("Nothing Affected. Make sure your user is exist in DB")
	}

	return ur, nil
}

func (m *mysqlUserRepository) Delete(id int64) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "DELETE FROM users WHERE id = ?"
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		log.Fatal(err)
		return false, models.INTERNAL_SERVER_ERROR
	}
	res, err := stmt.Exec(id)
	if err != nil {
		log.Fatal(err)
		return false, models.INTERNAL_SERVER_ERROR
	}
	rowsAfected, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
		return false, models.INTERNAL_SERVER_ERROR
	}
	if rowsAfected <= 0 {
		return false, models.INTERNAL_SERVER_ERROR
	}

	return true, nil
}