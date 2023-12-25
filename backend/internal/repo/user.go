package repo

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pegov/rolecap/backend/internal/entity"
	"github.com/pegov/rolecap/backend/internal/model"
	"github.com/pegov/rolecap/backend/internal/util"
	"github.com/redis/go-redis/v9"
)

type AuthRepo interface {
	GetById(id int) (*entity.User, error)
	GetByEmail(email string) (*entity.User, error)
	GetByUsername(username string) (*entity.User, error)
	GetByLogin(login string) (*entity.User, error)
	Create(payload model.UserCreate) (*entity.User, error)
}

type authRepo struct {
	q     map[string]string
	db    *sqlx.DB
	cache *redis.Client
}

func NewAuthRepo(db *sqlx.DB, cache *redis.Client) AuthRepo {
	return &authRepo{
		util.LoadFromFile("internal/repo/sql/user.sql"),
		db,
		cache,
	}
}

func (repo *authRepo) GetById(id int) (*entity.User, error) {
	var userRow entity.UserRow
	err := repo.db.QueryRowx("SELECT * FROM auth_user WHERE id = $1", id).StructScan(&userRow)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return entity.UserFromUserRow(&userRow), nil
}

func (repo *authRepo) GetByEmail(email string) (*entity.User, error) {
	var userRow entity.UserRow
	err := repo.db.QueryRowx("SELECT * FROM auth_user WHERE email = $1", email).StructScan(&userRow)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return entity.UserFromUserRow(&userRow), nil
}

func (repo *authRepo) GetByUsername(email string) (*entity.User, error) {
	var userRow entity.UserRow
	err := repo.db.QueryRowx("SELECT * FROM auth_user WHERE username = $1", email).StructScan(&userRow)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return entity.UserFromUserRow(&userRow), nil
}

func (repo *authRepo) GetByLogin(login string) (*entity.User, error) {
	entity, err := repo.GetByUsername(login)
	if err != nil {
		return nil, err
	}

	if entity != nil {
		return entity, nil
	}

	return repo.GetByEmail(login)
}

func (repo *authRepo) Create(payload model.UserCreate) (*entity.User, error) {
	var id int
	now := time.Now()
	err := repo.db.QueryRowx(repo.q["create"], payload.Email, payload.Username, payload.Password, true, false, now, now).Scan(&id)
	if err != nil {
		return nil, err
	}

	return repo.GetById(id)
}
