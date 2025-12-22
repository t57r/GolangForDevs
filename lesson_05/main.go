package main

import (
	"fmt"

	"lesson5/documentstore"
	"lesson5/users"
)

func main() {
	userServiceExample()

	// marshalExample()
	// unmarshalExample()
}

func userServiceExample() {
	usersService, err := users.New()
	if err != nil {
		fmt.Printf("Couldn't create users service: %v\n", err)
		return
	}

	testCreateUser(usersService, "1", "Alice")
	testCreateUser(usersService, "2", "Bob")
	testCreateUser(usersService, "3", "Caren")
	testCreateUser(usersService, "1", "Ann") // should be failed due id conflict

	testListUsers(usersService)

	caren, err := usersService.GetUser("3")
	if err != nil {
		fmt.Printf("can't get user with id=3: %v\n", err)
		return
	}
	fmt.Printf("\nfound user with id=3 -> %v\n", caren)

	err = usersService.DeleteUser("2")
	if err != nil {
		fmt.Printf("error deleting user: %v\n", err)
		return
	}

	fmt.Printf("\nAfter deleteting user with id=2\n")
	testListUsers(usersService)
}

func testCreateUser(service *users.Service, id string, name string) {
	user, err := service.CreateUser(id, name)
	if err != nil {
		fmt.Printf("Failed to create user %v\n", err)
		return
	}
	fmt.Printf("User created %v\n", user)
}

func testListUsers(service *users.Service) {
	users, err := service.ListUsers()
	if err != nil {
		fmt.Printf("error listing users: %v\n", err)
		return
	}
	for index, user := range users {
		fmt.Printf("User[%d] => %v\n", index, user)
	}
}

type MyStruct struct {
	X int
}

func marshalExample() {
	s := &MyStruct{X: 5}
	doc, err := documentstore.MarshalDocument(s)
	if err != nil {
		fmt.Printf("failed to marshal document: %v\n", err)
		return
	}

	fmt.Printf("marshaled document: %v\n", doc)
}

func unmarshalExample() {
	doc := &documentstore.Document{
		Fields: map[string]documentstore.DocumentField{
			"X": {Type: documentstore.DocumentFieldTypeNumber, Value: 5},
		},
	}

	s := &MyStruct{}
	err := documentstore.UnmarshalDocument(doc, s)
	if err != nil {
		fmt.Printf("failed to unmarshal document: %v\n", err)
		return
	}

	fmt.Printf("unmarshaled document: %v\n", s)
}