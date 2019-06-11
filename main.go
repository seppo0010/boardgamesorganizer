package main

import (
	"github.com/seppo0010/boardgamesorganizer/meetings"
	"github.com/seppo0010/boardgamesorganizer/users"
	"log"
	"os"
)

func main() {
	token := os.Getenv("BGO_TELEGRAM_BOT_TOKEN")
	postgresMeetingsURL := os.Getenv("BGO_MEETINGS_POSTGRES_URL")
	postgresUsersURL := os.Getenv("BGO_USERS_POSTGRES_URL")
	mf, err := meetings.NewPostgres(&meetings.PostgresConfig{URL: postgresMeetingsURL, MigrationsPath: "./meetings/migrations"})
	if err != nil {
		log.Fatalf("error starting meetings factory: %#v", err)
	}
	uf, err := users.NewPostgres(&users.PostgresConfig{URL: postgresUsersURL, MigrationsPath: "./users/migrations"})
	if err != nil {
		log.Fatalf("error starting users factory: %#v", err)
	}
	err = startTelegram(token, mf, uf)
	if err != nil {
		log.Fatalf("error running telegram: %#v", err)
	}
}
