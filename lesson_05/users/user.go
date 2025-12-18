package users

import (
	"errors"

	"lesson5/documentstore"
)

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrUserAlreadyExist = errors.New("user already exist")
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Service struct {
	coll *documentstore.Collection
}

func New() (*Service, error) {
	store := documentstore.NewStore()
	config := &documentstore.CollectionConfig{PrimaryKey: "ID"}

	userCollection, err := store.CreateCollection("users", config)
	if err != nil {
		return nil, err
	}

	service := Service{
		coll: userCollection,
	}
	return &service, nil
}

func (s *Service) CreateUser(id string, name string) (*User, error) {
	_, exist := s.coll.Get(id)
	if exist {
		return nil, ErrUserAlreadyExist
	}

	newUser := User{
		ID:   id,
		Name: name,
	}
	newUserDocument, err := documentstore.MarshalDocument(&newUser)
	if err != nil {
		return nil, err
	}

	s.coll.Put(*newUserDocument)
	return &newUser, nil
}

func (s *Service) ListUsers() ([]User, error) {
	userDocuments := s.coll.List()
	users := make([]User, 0, len(userDocuments))
	for _, doc := range userDocuments {
		var user User
		err := documentstore.UnmarshalDocument(&doc, &user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (s *Service) GetUser(userID string) (*User, error) {
	userDoc, found := s.coll.Get(userID)
	if !found {
		return nil, ErrUserNotFound
	}
	var user User
	err := documentstore.UnmarshalDocument(userDoc, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Service) DeleteUser(userID string) error {
	deleted := s.coll.Delete(userID)
	if !deleted {
		return ErrUserNotFound
	}
	return nil
}
