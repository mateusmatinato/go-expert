package user

import (
	"strconv"

	"github.com/google/uuid"
)

type User struct {
	Id   string
	Name string
	Age  uint64
}

func NewUser(name string, ageStr string) (*User, error) {
	age, err := strconv.ParseUint(ageStr, 10, 64)
	if err != nil {
		return nil, err
	}

	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	return &User{
		Id:   id.String(),
		Name: name,
		Age:  age,
	}, nil
}
