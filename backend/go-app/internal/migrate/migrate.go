package migrate

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

type Config struct {
	MigrationsDir string
	Timeout       time.Duration

	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func LoadConfigFromEnv() Config {
	return Config{
		MigrationsDir: getenv("MIGRATIONS_DIR", "migrations"),
		Timeout:       mustDuration("MIGRATIONS_TIMEOUT", 30*time.Second),

		Host:     mustEnv("PGHOST", "localhost"),
		Port:     mustEnv("PGPORT", "5432"),
		User:     mustEnv("PGUSER", "postgres"),
		Password: os.Getenv("PGPASSWORD"),
		DBName:   mustEnv("PGDATABASE", "postgres"),
		SSLMode:  mustEnv("PGSSLMODE", "disable"),
	}
}

func mustEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func mustDuration(key string, def time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		panic(fmt.Sprintf("invalid duration in %s: %s", key, v))
	}
	return d
}

func Up(ctx context.Context, cfg Config) error {
	if cfg.Timeout <= 0 {
		cfg.Timeout = 30 * time.Second
	}
	ctx, cancel := context.WithTimeout(ctx, cfg.Timeout)
	defer cancel()

	dsn, err := BuildPostgresDSN(cfg)
	if err != nil {
		return fmt.Errorf("build dsn: %w", err)
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("ping db: %w", err)
	}

	dirAbs, err := filepath.Abs(cfg.MigrationsDir)
	if err != nil {
		return fmt.Errorf("resolve migrations dir: %w", err)
	}

	goose.SetDialect("postgres")

	if err := goose.UpContext(ctx, db, dirAbs); err != nil {
		return fmt.Errorf("goose up: %w", err)
	}
	return nil
}

func Down(ctx context.Context, cfg Config) error {
	if cfg.Timeout <= 0 {
		cfg.Timeout = 30 * time.Second
	}
	ctx, cancel := context.WithTimeout(ctx, cfg.Timeout)
	defer cancel()

	dsn, err := BuildPostgresDSN(cfg)
	if err != nil {
		return fmt.Errorf("build dsn: %w", err)
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("ping db: %w", err)
	}

	dirAbs, err := filepath.Abs(cfg.MigrationsDir)
	if err != nil {
		return fmt.Errorf("resolve migrations dir: %w", err)
	}

	goose.SetDialect("postgres")
	if err := goose.DownContext(ctx, db, dirAbs); err != nil {
		return fmt.Errorf("goose down: %w", err)
	}
	return nil
}

func BuildPostgresDSN(cfg Config) (string, error) {
	if cfg.Host == "" || cfg.Port == "" || cfg.User == "" || cfg.DBName == "" {
		return "", fmt.Errorf("missing required postgres config (host/port/user/dbname)")
	}

	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.User, cfg.Password),
		Host:   fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Path:   cfg.DBName,
	}
	q := url.Values{}
	if cfg.SSLMode != "" {
		q.Set("sslmode", cfg.SSLMode)
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getenvDuration(key string, def time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return def
}
