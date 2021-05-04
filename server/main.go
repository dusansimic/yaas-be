package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	"github.com/dusansimic/yaas/engine"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/komkom/toml"

	_ "github.com/lib/pq"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	c, err := readTomlConfig()
	if err != nil {
		panic(err)
	}

	fmt.Println(c)

	// pgConnStr := os.Getenv("DB_PARAMS")

	db, err := sql.Open("postgres", c.DBConnStr)
	if err != nil {
		panic(err)
	}

	// migrationsPath := os.Getenv("MIGRATIONS_PATH")

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		panic(err)
	}

	m, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s", c.MigrationsPath), "postgres", driver)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if err != migrate.ErrNoChange {
			panic(err)
		}
		fmt.Println(err)
	}

	cookieStore := cookie.NewStore([]byte("verysecretkey"))
	serverStore := memstore.NewStore([]byte("verysecretkey"))

	// frontendBaseURL := os.Getenv("FRONTEND_BASE_URL")

	r := engine.New(
		engine.WithDB(db),
		engine.WithCookieStore(cookieStore),
		engine.WithServerStore(serverStore),
		engine.WithParam("frontBase", c.FrontendBase),
	)

	// port := os.Getenv("PORT")
	r.Run(fmt.Sprintf(":%d", c.Port))
}

func readTomlConfig() (config, error) {
	var c config

	f, err := os.Open("./config.toml")
	if err != nil {
		return config{}, err
	}

	if err := json.NewDecoder(toml.New(f)).Decode(&c); err != nil {
		return config{}, err
	}

	return c, nil
}

type config struct {
	DBConnStr      string `json:"data_source"`
	MigrationsPath string `json:"migrations_source"`
	FrontendBase   string `json:"frontend_base"`
	Port           int    `json:"port"`
}
