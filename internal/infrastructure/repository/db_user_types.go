package repository

import (
	"time"

	"github.com/mesirendon/go-tdd/internal/models"
)

type user struct {
	ID        string    `dynamodbav:"id"`
	FirstName string    `dynamodbav:"first_name"`
	LastName  string    `dynamodbav:"last_name"`
	Phone     string    `dynamodbav:"phone"`
	CreatedAt time.Time `dynamodbav:"created_at,unixtime"`
	UpdatedAt time.Time `dynamodbav:"updated_at,unixtime"`
}

func (u user) toModel() models.User {
	return models.User{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Phone:     u.Phone,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
