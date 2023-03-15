package repository

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/mesirendon/go-tdd/internal/models"
	"github.com/mesirendon/go-tdd/internal/utils"
)

type userDBClient interface {
	PutItem(
		ctx context.Context,
		params *dynamodb.PutItemInput,
		optFns ...func(*dynamodb.Options),
	) (*dynamodb.PutItemOutput, error)
}

type UserDBRepository struct {
	dbClient userDBClient
	table    string
	now      utils.Now
	uuid     utils.UUID
}

func NewUserDBRepository(
	dbClient userDBClient,
	table string,
	now utils.Now,
	uuid utils.UUID,
) *UserDBRepository {
	return &UserDBRepository{
		dbClient: dbClient,
		table:    table,
		now:      now,
		uuid:     uuid,
	}
}

func (r *UserDBRepository) Save(
	ctx context.Context,
	user models.User,
) (*models.User, error) {
	now := r.now()
	user.ID = r.uuid()
	u := r.toEntity(user)
	u.CreatedAt = now
	u.UpdatedAt = now

	item, _ := attributevalue.MarshalMap(u)
	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: &r.table,
	}
	_, _ = r.dbClient.PutItem(ctx, input)

	usr := u.toModel()

	return &usr, nil
}

func (r *UserDBRepository) toEntity(u models.User) user {
	return user{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Phone:     u.Phone,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
