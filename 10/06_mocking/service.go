package mocking

import (
	"errors"
	"fmt"
)

type Service struct {
	storage Storage
}

func NewService(storage Storage) *Service {
	return &Service{
		storage: storage,
	}
}

func (s *Service) CreateUser(user User) error {
	if user.Name == "Peter" {
		return errors.New("not Peter!")
	}

	err := s.storage.CreateUser(user)
	if err != nil {
		return err
	}

	fmt.Println("user was successfully created")

	return nil
}
