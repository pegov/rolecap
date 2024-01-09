package repo

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pegov/rolecap/backend/internal/entity"
	"github.com/pegov/rolecap/backend/internal/model"
	"github.com/pegov/rolecap/backend/internal/util"
	"github.com/pegov/rolecap/backend/pkg/private/editor"
	"github.com/redis/go-redis/v9"
)

type GenRepo interface {
	GetById(id int) (*entity.Gen, error)
	Create(user *model.AccessPayload, head editor.Head, body editor.Body, bodyMeta editor.BodyMeta) (int, error)
}

type genRepo struct {
	q     map[string]string
	db    *sqlx.DB
	cache *redis.Client
}

func NewGenRepo(db *sqlx.DB, cache *redis.Client) GenRepo {
	return &genRepo{
		util.LoadFromFile("internal/repo/sql/gen.sql"),
		db,
		cache,
	}
}

func (repo *genRepo) GetById(id int) (*entity.Gen, error) {
	var genRow entity.GenRow
	err := repo.db.QueryRowx(repo.q["get_by_id"], id).StructScan(&genRow)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return entity.GenFromGenRow(&genRow), nil
}

func (repo *genRepo) Create(
	user *model.AccessPayload,
	head editor.Head,
	body editor.Body,
	bodyMeta editor.BodyMeta,
) (int, error) {
	var id int
	now := time.Now()
	err := repo.db.QueryRow(repo.q["create"], user.Id, head.Title, head.Description, head.Access, body, bodyMeta, now, now).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
