package models

import (
	"errors"
	"strings"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"-"`
}

// For this demo, we're storing the user list in memory
// We also have some users predefined.
// In a real application, this list will most likely be fetched
// from a database. Moreover, in production settings, you should
// store passwords securely by salting and hashing them instead
// of using them as we're doing in this demo
var UserList = []User{
	{Username: "user1", Password: "pass1"},
	{Username: "user2", Password: "pass2"},
	{Username: "user3", Password: "pass3"},
}

// Register a new user with the given username and password
func RegisterNewUser(username, password string) (*User, error) {
	if strings.TrimSpace(password) == "" {
		return nil, errors.New("the password can't be empty")
	} else if !IsUsernameAvailable(username) {
		return nil, errors.New("the username isn't available")
	}

	u := User{Username: username, Password: password}

	UserList = append(UserList, u)

	return &u, nil
}

// Check if the supplied username is available
func IsUsernameAvailable(username string) bool {
	for _, u := range UserList {
		if u.Username == username {
			return false
		}
	}
	return true
}

//Function to validate the login credentials
func IsUserValid(username, password string) bool {
	for _, u := range UserList {
		if u.Username == username && u.Password == password {
			return true
		}
	}
	return false
}
