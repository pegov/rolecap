package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pegov/rolecap/backend/internal/handler"
	"github.com/pegov/rolecap/backend/internal/repo"
	"github.com/pegov/rolecap/backend/internal/router"
	"github.com/pegov/rolecap/backend/internal/util"
	"github.com/redis/go-redis/v9"
)

func main() {
	postgresUrl, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		postgresUrl = "postgres://postgres:postgres@127.0.0.1:5432/rolecap"
	}
	postgresDb, err := sqlx.Connect("pgx", postgresUrl)
	if err != nil {
		log.Fatalln(err)
	}
	defer postgresDb.Close()

	redisUrl, ok := os.LookupEnv("CACHE_URL")
	if !ok {
		redisUrl = "redis://127.0.0.1:6379/2"
	}
	redisOptions, err := redis.ParseURL(redisUrl)
	if err != nil {
		log.Fatalln(err)
	}
	redisCache := redis.NewClient(redisOptions)
	if err := redisCache.Ping(context.Background()).Err(); err != nil {
		log.Fatalln(err)
	}
	defer redisCache.Close()

	authRepo := repo.NewAuthRepo(postgresDb, redisCache)

	privateKeyBytes, err := os.ReadFile("./ed25519_1.key")
	if err != nil {
		log.Fatalln(err)
	}
	publicKeyBytes, err := os.ReadFile("./ed25519_1.pub")
	if err != nil {
		log.Fatalln(err)
	}
	jwtBackend := util.NewJwtBackend(privateKeyBytes, publicKeyBytes, "1")

	passwordHasher := util.NewBcryptPasswordHasher()

	authHandler := handler.NewAuthHandler(authRepo, jwtBackend, passwordHasher)

	r := gin.Default()
	router.SetupUserRouter(r, authHandler)

	host, ok := os.LookupEnv("HOST")
	if !ok {
		host = "0.0.0.0"
	}

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}
	bindAddr := fmt.Sprintf("%v:%v", host, port)
	r.Run(bindAddr)
}
