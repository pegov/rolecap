//go:build !internal

package repo

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pegov/rolecap/backend/internal/model"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

var (
	TestDatabaseName = "rolecap_test"
)

type Connections struct {
	Db    *sqlx.DB
	Cache *redis.Client
}

var connections Connections
var repo AuthRepo

func SetupAuthDatabase() {
	postgresDb, err := sqlx.Open("pgx", "postgres://postgres:postgres@127.0.0.1:5432/")
	if err != nil {
		panic(err)
	}

	_, err = postgresDb.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %v WITH (FORCE)", TestDatabaseName))
	if err != nil {
		panic(err)
	}

	_, err = postgresDb.Exec(fmt.Sprintf("CREATE DATABASE %v", TestDatabaseName))
	if err != nil {
		panic(err)
	}
	postgresDb.Close()

	postgresDb, err = sqlx.Connect("pgx", fmt.Sprintf("postgres://postgres:postgres@127.0.0.1:5432/%v", TestDatabaseName))
	if err != nil {
		panic(err)
	}

	sqlInitContentBytes, err := os.ReadFile("../migrations/01_init.sql")
	if err != nil {
		panic(err)
	}

	sqlInitContent := string(sqlInitContentBytes)
	_, err = postgresDb.Exec(sqlInitContent)
	if err != nil {
		panic(err)
	}

	redisOptions, err := redis.ParseURL("redis://127.0.0.1:6379/3")
	if err != nil {
		log.Fatalln(err)
	}
	redisCache := redis.NewClient(redisOptions)
	if err := redisCache.Ping(context.Background()).Err(); err != nil {
		log.Fatalln(err)
	}
	defer redisCache.Close()

	redisCache.FlushAll(context.Background())

	repo = NewAuthRepo(postgresDb, redisCache)

	connections = Connections{postgresDb, redisCache}
}

func TeardownAuthDatabase() {
	connections.Db.Close()

	postgresDb, err := sqlx.Open("pgx", "postgres://postgres:postgres@127.0.0.1:5432/")
	if err != nil {
		panic(err)
	}

	_, err = postgresDb.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %v WITH (FORCE)", TestDatabaseName))
	if err != nil {
		panic(err)
	}

	connections.Cache.FlushAll(context.Background())
	defer connections.Cache.Close()
}

func TestMain(m *testing.M) {
	SetupAuthDatabase()
	code := m.Run()
	TeardownAuthDatabase()
	os.Exit(code)
}

var (
	TestEmail1    = "test@test.com"
	TestUsername1 = "testusername"
	TestPassword1 = "hunter2"
)

func BeforeEach(db *sqlx.DB, cache *redis.Client) {

}

func AfterEach() {
	connections.Db.Exec("DELETE FROM auth_user WHERE 1 = 1")
	connections.Db.Exec("ALTER SEQUENCE auth_user_id_seq RESTART WITH 0")
	connections.Cache.FlushAll(context.Background())
}

func TestCreate(t *testing.T) {
	entity, err := repo.Create(model.UserCreate{
		Email:    TestEmail1,
		Username: TestUsername1,
		Password: TestPassword1,
	})
	assert.Nil(t, err)
	assert.NotNil(t, entity)

	assert.Equal(t, entity.Id, 1)
	assert.Equal(t, TestUsername1, entity.Username)
	assert.Equal(t, TestEmail1, entity.Email)
	assert.Equal(t, TestPassword1, entity.Password)
	AfterEach()
}

func TestGetByUsername(t *testing.T) {
	_, err := repo.Create(model.UserCreate{
		Email:    TestEmail1,
		Username: TestUsername1,
		Password: TestPassword1,
	})
	assert.Nil(t, err)

	entity, err := repo.GetByUsername(TestUsername1)
	assert.Nil(t, err)
	assert.NotNil(t, entity)
	assert.Equal(t, TestUsername1, entity.Username)
	AfterEach()
}

func TestGetByEmail(t *testing.T) {
	_, err := repo.Create(model.UserCreate{
		Email:    TestEmail1,
		Username: TestUsername1,
		Password: TestPassword1,
	})
	assert.Nil(t, err)

	entity, err := repo.GetByEmail(TestEmail1)
	assert.Nil(t, err)
	assert.NotNil(t, entity)
	assert.Equal(t, TestEmail1, entity.Email)
	AfterEach()
}

func TestGetByLogin(t *testing.T) {
	_, err := repo.Create(model.UserCreate{
		Email:    TestEmail1,
		Username: TestUsername1,
		Password: TestPassword1,
	})
	assert.Nil(t, err)

	entity, err := repo.GetByLogin(TestEmail1)
	assert.Nil(t, err)
	assert.NotNil(t, entity)
	assert.Equal(t, TestEmail1, entity.Email)

	entity, err = repo.GetByLogin(TestUsername1)
	assert.Nil(t, err)
	assert.NotNil(t, entity)
	assert.Equal(t, TestUsername1, entity.Username)
	AfterEach()
}
