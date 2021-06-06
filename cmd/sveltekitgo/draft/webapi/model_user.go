package webapi

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

// User -
type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
}

var (
	userStorage = "cmd/sveltekitgo/storage/users.json"
)

// LoadUsers -
func LoadUsers() ([]User, error) {
	b, err := os.ReadFile(userStorage)
	if err != nil {
		return nil, err
	}

	users := make([]User, 0)

	err = json.Unmarshal(b, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// SaveUsers -
func SaveUsers(users []User) error {
	b, err := json.Marshal(users)
	if err != nil {
		return err
	}

	return os.WriteFile(userStorage, b, 0644)
}

// CreateUser -
func CreateUser(user User) error {
	users, err := LoadUsers()
	if err != nil {
		return err
	}

	passhash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(passhash)

	uid, err := generateUUID()
	if err != nil {
		return err
	}
	user.ID = uid

	users = append(users, user)

	return SaveUsers(users)
}

// DeleteUser -
func DeleteUser(id string) error {
	users, err := LoadUsers()
	if err != nil {
		return err
	}

	for i, v := range users {
		if v.ID == id {
			users = removeItem(users, i)
			break
		}
	}

	return SaveUsers(users)
}

func removeItem(s []User, index int) []User {
	return append(s[:index], s[index+1:]...)
}

// UpdateUser -
func UpdateUser(user User) error {
	users, err := LoadUsers()
	if err != nil {
		return err
	}

	for i, v := range users {
		if v.ID == user.ID {
			passhash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
			if err != nil {
				return err
			}

			users[i] = User{
				ID:        user.ID,
				Email:     user.Email,
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Password:  string(passhash),
			}
			break
		}
	}

	return SaveUsers(users)
}

// GetUser -
func GetUser(email string) (*User, error) {
	users, err := LoadUsers()
	if err != nil {
		return nil, err
	}

	for _, v := range users {
		if v.Email == email {
			return &v, nil
		}
	}

	return nil, errors.New("not found")
}

// generateUUID for use as an random identifier.
func generateUUID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
}

func passwordMatch(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
