package store

import (
	"github.com/harborapp/harbor-api/model"
	"github.com/jinzhu/gorm"
	"golang.org/x/net/context"
)

// GetUsers retrieves all available users from the database.
func GetUsers(c context.Context) (*model.Users, error) {
	return FromContext(c).GetUsers()
}

// CreateUser creates a new user.
func CreateUser(c context.Context, record *model.User) error {
	return FromContext(c).CreateUser(record)
}

// UpdateUser updates a user.
func UpdateUser(c context.Context, record *model.User) error {
	return FromContext(c).UpdateUser(record)
}

// DeleteUser deletes a user.
func DeleteUser(c context.Context, record *model.User) error {
	return FromContext(c).DeleteUser(record)
}

// GetUser retrieves a specific user from the database.
func GetUser(c context.Context, id string) (*model.User, *gorm.DB) {
	return FromContext(c).GetUser(id)
}
