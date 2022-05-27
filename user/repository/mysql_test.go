package repository_test

import (
	"database/sql"
	"log"
	"testing"
	"time"

	models "simple-grpc-go/user"
	userRepo "simple-grpc-go/user/repository"

	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

func TestFetch(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id","name","email","age","created_at","updated_at"}).
		AddRow(1, "test 1", "test1@mail.com", 22, time.Now(), time.Now()).
		AddRow(2, "test 2", "test2@mail.com", 11, time.Now(), time.Now()).
		AddRow(3, "test 3", "test3@mail.com", 33, time.Now(), time.Now())

	query := "SELECT id,name,email,age,created_at,updated_at FROM users WHERE ID > \\? LIMIT \\?"

	mock.ExpectQuery(query).WillReturnRows(rows)
	u := userRepo.NewMysqlUserRepository(db)
	cursor := "tes"
	num := int64(5)
	list, err := u.Fetch(cursor, num)
	assert.NoError(t, err)
	assert.Len(t, list, 3)
}

func TestGetByID(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id","name","email","age","created_at","updated_at"}).
		AddRow(1, "test 1", "test1@mail.com", 22, time.Now(), time.Now())

	query := "SELECT id,name,email,age,created_at,updated_at FROM users WHERE ID = \\?"

	mock.ExpectQuery(query).WillReturnRows(rows)
	u := userRepo.NewMysqlUserRepository(db)

	id := int64(1)
	user, err := u.GetByID(id)
	assert.NoError(t, err)
	assert.NotNil(t, user)
}

func TestCreateUser(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	user := &models.User{
		ID: 1,
		Name: "lacrose",
		Email: "lacrose@gmail.com",
		Age: 22,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	query := "INSERT users SET id = \\?, name = \\?, email = \\?, age = \\?, created_at = \\?, updated_at = \\?"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().
		WithArgs(user.ID, user.Name, user.Email, user.Age, user.CreatedAt, user.UpdatedAt).
		WillReturnResult(sqlmock.NewResult(0, 1))

	u := userRepo.NewMysqlUserRepository(db)

	lastId, err := u.Store(user)
	assert.NoError(t, err)
	assert.Equal(t, int64(15), lastId)
}

func TestUpdateUser(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	user := &models.User{
		ID: 14,
		Name: "lacrose",
		Email: "lacrose@gmail.com",
		Age: 22,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	query := "UPDATE users SET name = \\?, email = \\?, age = \\?, created_at = \\?, updated_at = \\? WHERE id = \\?"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().
		WithArgs(user.Name, user.Email, user.Age, user.CreatedAt, user.UpdatedAt, user.ID).
		WillReturnResult(sqlmock.NewResult(12, 1))

	u := userRepo.NewMysqlUserRepository(db)

	s, err := u.Update(user)
	assert.NoError(t, err)
	assert.NotNil(t, s)
}

func TestDeleteUser(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	query := "DELETE FROM users WHERE id = \\?"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(13).WillReturnResult(sqlmock.NewResult(13, 1))

	u := userRepo.NewMysqlUserRepository(db)

	num := int64(13)
	anDeleteStatus, err := u.Delete(num)
	assert.NoError(t, err)
	assert.True(t, anDeleteStatus)
}