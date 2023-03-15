package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/mesirendon/go-tdd/internal/infrastructure/repository"
	"github.com/mesirendon/go-tdd/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type userDBClientMock struct {
	mock.Mock
}

func (m *userDBClientMock) PutItem(
	ctx context.Context,
	params *dynamodb.PutItemInput,
	optFns ...func(*dynamodb.Options),
) (*dynamodb.PutItemOutput, error) {
	args := m.Called(ctx, params)

	if err := args.Error(1); err != nil {
		return nil, err
	}
	return args.Get(0).(*dynamodb.PutItemOutput), args.Error(1)
}

type timeNowMock struct{}

func (m *timeNowMock) Now() time.Time {
	return now()
}

func now() time.Time {
	return time.Date(2023, 02, 15, 23, 35, 5, 0, time.Local)
}

type uuidMock struct {
	mock.Mock
}

func (m *uuidMock) UUID() string {
	args := m.Called()
	return args.Get(0).(string)
}

func dynamoDBMarshal(t *testing.T, v any) types.AttributeValue {
	var a types.AttributeValue
	var err error

	switch i := v.(type) {
	case time.Time:
		a, err = attributevalue.Marshal(i.UTC().Unix())
	default:
		a, err = attributevalue.Marshal(i)
	}

	if err != nil {
		t.Fatal(err)
	}

	return a
}

func TestUserDBRepository_Save(t *testing.T) {
	table := "users"
	type args struct {
		user models.User
	}
	tests := []struct {
		name    string
		args    args
		mocker  func(*userDBClientMock, *uuidMock, *testing.T)
		want    *models.User
		wantErr *string
	}{
		{
			name: "Successfully save a user",
			args: args{
				user: models.User{
					FirstName: "John",
					LastName:  "Doe",
					Phone:     "+12345678",
				},
			},
			mocker: func(
				db *userDBClientMock,
				uuid *uuidMock,
				t *testing.T,
			) {
				uuid.On("UUID").Return("usr-uuid-123")

				input := &dynamodb.PutItemInput{
					Item: map[string]types.AttributeValue{
						"id":         dynamoDBMarshal(t, "usr-uuid-123"),
						"first_name": dynamoDBMarshal(t, "John"),
						"last_name":  dynamoDBMarshal(t, "Doe"),
						"phone":      dynamoDBMarshal(t, "+12345678"),
						"created_at": dynamoDBMarshal(t, now()),
						"updated_at": dynamoDBMarshal(t, now()),
					},
					TableName: &table,
				}
				db.On("PutItem", context.Background(), input).
					Return(&dynamodb.PutItemOutput{}, nil)
			},
			want: &models.User{
				ID:        "usr-uuid-123",
				FirstName: "John",
				LastName:  "Doe",
				Phone:     "+12345678",
				CreatedAt: now(),
				UpdatedAt: now(),
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientMock := new(userDBClientMock)
			nowMock := new(timeNowMock)
			uuidMock := new(uuidMock)
			r := repository.NewUserDBRepository(
				clientMock,
				table,
				nowMock.Now,
				uuidMock.UUID,
			)
			tt.mocker(clientMock, uuidMock, t)
			got, err := r.Save(context.Background(), tt.args.user)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorContains(t, err, *tt.wantErr)
			} else {
				assert.NoError(t, err)
			}

			assert.EqualValues(t, tt.want, got)
			clientMock.AssertExpectations(t)
		})
	}
}
