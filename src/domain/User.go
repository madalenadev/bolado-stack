package domain

import "github.com/google/uuid"

// User represents domain for user model
type User struct {
	ID        uuid.UUID
	FirstName string
	LastName  string
}

// GetFullName get fullname of an User
func (u User) GetFullName() string {
	return u.FirstName + " " + u.LastName
}
