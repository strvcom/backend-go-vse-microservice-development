package mocking

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockStorage struct {
	mock.Mock
}

func (m *mockStorage) CreateUser(user User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *mockStorage) ReadUser(name string) (User, error) {
	args := m.Called(name)
	return args.Get(0).(User), args.Error(1)
}

func (m *mockStorage) DeleteUser(name string) error {
	args := m.Called(name)
	return args.Error(0)
}

func Test_Service_CreateUser(t *testing.T) {
	cases := []struct {
		name           string
		mockStorage    *mockStorage
		user           User
		expectingError bool
	}{
		{
			name: "Alice ok",
			mockStorage: func() *mockStorage {
				m := &mockStorage{}
				m.On("CreateUser", User{Name: "Alice"}).Return(nil)
				return m
			}(),
			user: User{
				Name: "Alice",
			},
			expectingError: false,
		},
		{
			name: "Bob ok",
			mockStorage: func() *mockStorage {
				m := &mockStorage{}
				m.On("CreateUser", User{Name: "Bob"}).Return(nil)
				return m
			}(),
			user: User{
				Name: "Bob",
			},
			expectingError: false,
		},
		{
			name:        "Peter error",
			mockStorage: &mockStorage{},
			user: User{
				Name: "Peter",
			},
			expectingError: true,
		},
		{
			name: "Bob error",
			mockStorage: func() *mockStorage {
				m := &mockStorage{}
				m.On("CreateUser", User{Name: "Bob"}).Return(errors.New("unexpected error"))
				return m
			}(),
			user: User{
				Name: "Bob",
			},
			expectingError: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			service := NewService(c.mockStorage)
			err := service.CreateUser(c.user)
			if c.expectingError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			c.mockStorage.AssertExpectations(t)
		})
	}
}
