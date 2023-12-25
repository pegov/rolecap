package mock_repo

import (
	"time"

	"github.com/pegov/rolecap/backend/internal/entity"
	"github.com/pegov/rolecap/backend/internal/model"
	"github.com/pegov/rolecap/backend/internal/repo"
)

type MockAuthRepo struct {
	Users []entity.User
}

func NewMockAuthRepo(users []entity.User) repo.AuthRepo {
	return &MockAuthRepo{
		Users: users,
	}
}

func (repo *MockAuthRepo) GetById(id int) (*entity.User, error) {
	for _, user := range repo.Users {
		if user.Id == id {
			return &user, nil
		}
	}

	return nil, nil
}

func (repo *MockAuthRepo) GetByEmail(email string) (*entity.User, error) {
	for _, user := range repo.Users {
		if user.Email == email {
			return &user, nil
		}
	}

	return nil, nil
}

func (repo *MockAuthRepo) GetByUsername(username string) (*entity.User, error) {
	for _, user := range repo.Users {
		if user.Email == username {
			return &user, nil
		}
	}

	return nil, nil
}

func (repo *MockAuthRepo) GetByLogin(login string) (*entity.User, error) {
	for _, user := range repo.Users {
		if user.Username == login {
			return &user, nil
		}
	}

	for _, user := range repo.Users {
		if user.Email == login {
			return &user, nil
		}
	}

	return nil, nil
}

func (repo *MockAuthRepo) Create(payload model.UserCreate) (*entity.User, error) {
	id := len(repo.Users)
	now := time.Now()
	entity := entity.User{
		Id:       id,
		Email:    payload.Email,
		Username: payload.Username,
		Password: payload.Password,

		Active:   true,
		Verified: true,

		CreatedAt: now,
		LastLogin: now,

		OAuth: nil,
	}
	repo.Users = append(repo.Users, entity)
	return &entity, nil
}
