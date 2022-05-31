package usecase_test

import (
	"strconv"
	"testing"

	models "simple-grpc-go/user"
	"simple-grpc-go/user/repository/mocks"
	ucase "simple-grpc-go/user/usecase"

	"github.com/bxcodec/faker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFetch(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	var mockUser models.User
	err := faker.FakeData(&mockUser)
	assert.NoError(t, err)

	mockListUser := make([]*models.User, 0)
	mockListUser = append(mockListUser, &mockUser)
	mockUserRepo.On("Fetch", mock.AnythingOfType("string"), mock.AnythingOfType("int64")).Return(mockListUser, nil)
	u := ucase.NewUserUsecase(mockUserRepo)
	num := int64(1)
	cursor := "0"
	list, nextCursor, err := u.Fetch(cursor, num)
	cursorExpected := strconv.Itoa(int(mockUser.ID))
	assert.Equal(t, cursorExpected, nextCursor)
	assert.NotEmpty(t, nextCursor)
	assert.NoError(t, err)
	assert.Len(t, list, len(mockListUser))

	mockUserRepo.AssertCalled(t, "Fetch", mock.AnythingOfType("string"), mock.AnythingOfType("int64"))
}

func TestGetByID(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	var mockUser models.User
	err := faker.FakeData(&mockUser)
	assert.NoError(t, err)

	mockUserRepo.On("GetByID", mock.AnythingOfType("int64")).Return(&mockUser, nil)
	defer mockUserRepo.AssertCalled(t, "GetByID", mock.AnythingOfType("int64"))
	u := ucase.NewUserUsecase(mockUserRepo)

	a, err := u.GetByID(mockUser.ID)

	assert.NoError(t, err)
	assert.NotNil(t, a)
}

func TestStore(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	var mockUser models.User
	err := faker.FakeData(&mockUser)
	assert.NoError(t, err)
	//set to 0 because this is test from Client, and ID is an AutoIncreament
	tempMockUser := mockUser
	tempMockUser.ID = 0

	mockUserRepo.On("Store", &tempMockUser).Return(mockUser.ID, nil)
	defer mockUserRepo.AssertExpectations(t)

	u := ucase.NewUserUsecase(mockUserRepo)

	a, err := u.Store(&tempMockUser)

	assert.NoError(t, err)
	assert.NotNil(t, a)
	assert.Equal(t, mockUser.Name, tempMockUser.Name)
}

func TestUpdate(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	var mockUser models.User
	err := faker.FakeData(&mockUser)
	assert.NoError(t, err)

	tempMockUser := mockUser
	tempMockUser.ID = 0

	mockUserRepo.On("Update", &tempMockUser).Return(mockUser.ID, nil)
	defer mockUserRepo.AssertExpectations(t)

	u := ucase.NewUserUsecase(mockUserRepo)

	a, err := u.Update(&tempMockUser)

	assert.NoError(t, err)
	assert.NotNil(t, a)
	assert.Equal(t, mockUser.Name, tempMockUser.Name)
}

func TestDelete(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	var mockUser models.User
	err := faker.FakeData(&mockUser)
	assert.NoError(t, err)

	mockUserRepo.On("GetByID", mock.AnythingOfType("int64")).Return(&mockUser, models.NOT_FOUND_ERROR)
	defer mockUserRepo.AssertCalled(t, "GetByID", mock.AnythingOfType("int64"))

	mockUserRepo.On("Delete", mock.AnythingOfType("int64")).Return(true, nil)
	defer mockUserRepo.AssertCalled(t, "Delete", mock.AnythingOfType("int64"))

	u := ucase.NewUserUsecase(mockUserRepo)

	a, err := u.Delete(mockUser.ID)

	assert.NoError(t, err)
	assert.True(t, a)
}