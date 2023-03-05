package main

import (
	"crud_app/internal/config"
	"crud_app/internal/repository/psql"
	"crud_app/internal/service"
	"crud_app/internal/transport/rest"
	"crud_app/pkg/database"
	"crud_app/pkg/hash"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"time"
)

const (
	CONFIG_DIR  = "configs"
	CONFIG_FILE = "main"
)

func main() {
	cfg, err := config.New(CONFIG_DIR, CONFIG_FILE)

	db, err := database.NewPostgresConnection(database.ConnectionInfo{
		Host:     "localhost",
		Port:     5432,
		Username: "postgres",
		DBName:   "postgres",
		SSLMode:  "disable",
		Password: "qwerty123",
	})
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	hasher := hash.NewSHA1Hasher("salt")

	userRep := psql.NewUsers(db)
	usersService := service.NewUsers(userRep, hasher, []byte("sample secret"), cfg.Auth.TokenTTL)
	handler := rest.NewHandler(usersService)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: handler.InitRouter(),
	}

	log.Println("SERVER STARTED AT", time.Now().Format(time.RFC3339))

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
