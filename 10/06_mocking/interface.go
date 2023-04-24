package mocking

type User struct {
	Name string
}

type Storage interface {
	CreateUser(user User) error
	ReadUser(name string) (User, error)
	DeleteUser(name string) error
}
