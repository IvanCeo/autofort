package migrate

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pressly/goose/v3"
)

// TestDisposablePostgresBootstrap never accepts a supplied DSN or loads .env.
// Explicit opt-in is required because it creates and removes its own container.
func TestDisposablePostgresBootstrap(t *testing.T) {
	if os.Getenv("AUTOFORT_DISPOSABLE_PG_TEST") != "1" {
		t.Skip("integration not run; set AUTOFORT_DISPOSABLE_PG_TEST=1 to create an isolated PostgreSQL 16 container")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Minute)
	defer cancel()
	if output, err := dockerCommand(ctx, "version", "--format", "{{.Server.Version}}"); err != nil {
		t.Fatalf("Docker unavailable; bootstrap NOT verified: %v: %s", err, output)
	}
	endpoint := os.Getenv("DOCKER_HOST")
	if os.Getenv("DOCKER_CONTEXT") != "" || endpoint == "" {
		args := []string{"context", "inspect", "--format", "{{.Endpoints.docker.Host}}"}
		if selected := os.Getenv("DOCKER_CONTEXT"); selected != "" {
			args = append(args, selected)
		}
		var err error
		endpoint, err = dockerCommand(ctx, args...)
		if err != nil {
			t.Fatalf("cannot verify local Docker endpoint: %v", err)
		}
	}
	if !isLocalDockerEndpoint(strings.TrimSpace(endpoint)) {
		t.Fatal("Docker endpoint is remote or unsupported; no container or database connection attempted")
	}
	var random [16]byte
	if _, err := rand.Read(random[:]); err != nil {
		t.Fatal(err)
	}
	token := hex.EncodeToString(random[:])
	var secret [32]byte
	if _, err := rand.Read(secret[:]); err != nil {
		t.Fatal(err)
	}
	password := hex.EncodeToString(secret[:])
	name := "autofort-bootstrap-" + token
	dbName := "bootstrap_" + token
	const user = "bootstrap"
	const label = "autofort.disposable-bootstrap"
	containerID, err := createBootstrapContainer(ctx, dockerCommand, name, token, []string{"run", "--detach", "--name", name,
		"--label", label + "=" + token,
		"--tmpfs", "/var/lib/postgresql/data:rw",
		"--publish", "127.0.0.1::5432",
		"--env", "POSTGRES_USER=" + user,
		"--env", "POSTGRES_PASSWORD=" + password,
		"--env", "POSTGRES_DB=" + dbName,
		"postgres:16"}, func(cleanup func() error) {
		t.Cleanup(func() {
			if err := cleanup(); err != nil {
				t.Error(err)
			}
		})
	})
	if err != nil {
		t.Fatalf("create disposable PostgreSQL: %v; no database connection attempted", err)
	}

	container, err := inspectContainer(ctx, dockerCommand, containerID)
	if err != nil {
		t.Fatal(err)
	}
	if container.ID != containerID || container.Name != "/"+name || container.Config.Labels[label] != token || container.Config.Image != "postgres:16" {
		t.Fatal("container identity/image/ownership verification failed; no database connection attempted")
	}
	if !hasIsolatedPostgresMounts(container) {
		t.Fatalf("unexpected container mounts: tmpfs=%v binds=%v mounts=%+v; no database connection attempted", container.HostConfig.Tmpfs, container.HostConfig.Binds, container.Mounts)
	}
	bindings := container.NetworkSettings.Ports["5432/tcp"]
	if len(bindings) != 1 || bindings[0].HostIP != "127.0.0.1" {
		t.Fatal("database port is not a single loopback binding; no database connection attempted")
	}
	port, err := strconv.Atoi(bindings[0].HostPort)
	if err != nil || port < 1 || port > 65535 {
		t.Fatal("invalid disposable database port; no database connection attempted")
	}
	_, sourceFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot locate real project migrations")
	}
	cfg := Config{
		MigrationsDir: filepath.Join(filepath.Dir(sourceFile), "..", "..", "migrations"),
		Timeout:       30 * time.Second,
		Host:          "127.0.0.1",
		Port:          bindings[0].HostPort,
		User:          user,
		Password:      password,
		DBName:        dbName,
		SSLMode:       "disable",
	}
	dsn, err := BuildPostgresDSN(cfg)
	if err != nil {
		t.Fatal(err)
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatal("open disposable database failed")
	}
	defer func() { _ = db.Close() }()
	readyCtx, readyCancel := context.WithTimeout(ctx, 60*time.Second)
	defer readyCancel()
	for db.PingContext(readyCtx) != nil {
		select {
		case <-readyCtx.Done():
			t.Fatal("disposable PostgreSQL was not ready within 60 seconds")
		case <-time.After(250 * time.Millisecond):
		}
	}
	var actualDB, actualUser, serverVersion string
	if err := db.QueryRowContext(ctx, "SELECT current_database(), current_user, current_setting('server_version_num')").Scan(&actualDB, &actualUser, &serverVersion); err != nil {
		t.Fatal(err)
	}
	version, err := strconv.Atoi(serverVersion)
	if err != nil || version < 160000 || version >= 170000 || actualDB != dbName || actualUser != user {
		t.Fatal("disposable database identity/PostgreSQL 16 verification failed; no migrations attempted")
	}
	var tableCount int
	if err := db.QueryRowContext(ctx, "SELECT count(*) FROM information_schema.tables WHERE table_schema = 'public'").Scan(&tableCount); err != nil || tableCount != 0 {
		t.Fatalf("database is not empty: tables=%d, err=%v; no migrations attempted", tableCount, err)
	}
	t.Logf("isolated target verified: container=%s postgres=%s tmpfs=/var/lib/postgresql/data loopback=127.0.0.1:%d no persistent mounts", name, serverVersion, port)
	// Reproduce the historical PostgreSQL 16 failure without leaving any changes.
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	_, oldSyntaxErr := tx.ExecContext(ctx, "CREATE TYPE IF NOT EXISTS user_role AS ENUM ('superuser', 'businessman');")
	if err := tx.Rollback(); err != nil {
		t.Fatal(err)
	}
	var pgErr *pgconn.PgError
	if !errors.As(oldSyntaxErr, &pgErr) || pgErr.Code != "42601" {
		t.Fatalf("expected PostgreSQL 16 syntax error for historical enum statement, got %v", oldSyntaxErr)
	}
	t.Logf("historical migration 7 reproduced in rolled-back transaction: SQLSTATE=%s message=%s", pgErr.Code, pgErr.Message)
	if err := Up(ctx, cfg); err != nil {
		t.Fatalf("fresh migrations 1-7 failed: %v", err)
	}
	assertBootstrapVersion(t, ctx, db, 7)
	assertBootstrapSchema(t, ctx, db)
	var migrationRowsBefore int
	if err := db.QueryRowContext(ctx, "SELECT count(*) FROM goose_db_version").Scan(&migrationRowsBefore); err != nil {
		t.Fatal(err)
	}
	if err := Up(ctx, cfg); err != nil {
		t.Fatalf("second Up must be a no-op: %v", err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	db, err = sql.Open("pgx", dsn)
	if err != nil {
		t.Fatal("reopen disposable database failed")
	}
	if err := Up(ctx, cfg); err != nil {
		t.Fatalf("Up after reopening must be a no-op: %v", err)
	}
	assertBootstrapVersion(t, ctx, db, 7)
	assertBootstrapSchema(t, ctx, db)
	var migrationRowsAfter int
	if err := db.QueryRowContext(ctx, "SELECT count(*) FROM goose_db_version").Scan(&migrationRowsAfter); err != nil || migrationRowsAfter != migrationRowsBefore {
		t.Fatalf("repeat Up changed migration history: before=%d after=%d err=%v", migrationRowsBefore, migrationRowsAfter, err)
	}
	if err := Down(ctx, cfg); err != nil {
		t.Fatalf("rollback migration 7 failed: %v", err)
	}
	assertBootstrapVersion(t, ctx, db, 6)
	var usersGone, roleGone bool
	if err := db.QueryRowContext(ctx, "SELECT to_regclass('public.users') IS NULL, NOT EXISTS (SELECT 1 FROM pg_type WHERE typname='user_role' AND typnamespace='public'::regnamespace)").Scan(&usersGone, &roleGone); err != nil || !usersGone || !roleGone {
		t.Fatalf("rollback did not remove users and enum: usersGone=%v roleGone=%v err=%v", usersGone, roleGone, err)
	}
	if err := Up(ctx, cfg); err != nil {
		t.Fatalf("reapply migration 7 failed: %v", err)
	}
	// Recreating an enum changes its OID; discard cached verification plans.
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	reopenedDB, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatal("reopen verification database after enum recreation failed")
	}
	db = reopenedDB
	assertBootstrapVersion(t, ctx, db, 7)
	assertBootstrapSchema(t, ctx, db)
	t.Log("PASS: fresh migrations 1-7, version 7, repeated/reopened Up no-op, rollback 7 to 6, reapply 7")
}

func assertBootstrapVersion(t *testing.T, ctx context.Context, db *sql.DB, expected int64) {
	t.Helper()
	version, err := goose.GetDBVersionContext(ctx, db)
	if err != nil || version != expected {
		t.Fatalf("Goose version: got=%d expected=%d err=%v", version, expected, err)
	}
}

func assertBootstrapSchema(t *testing.T, ctx context.Context, db *sql.DB) {
	t.Helper()
	var vehicleTypes, services int
	var roles string
	if err := db.QueryRowContext(ctx, "SELECT (SELECT count(*) FROM vehicle_types), (SELECT count(*) FROM services), array_to_string(enum_range(NULL::user_role), ',')").Scan(&vehicleTypes, &services, &roles); err != nil {
		t.Fatal(err)
	}
	if vehicleTypes != 130 || services != 50 || roles != "superuser,businessman" {
		t.Fatalf("unexpected seeded schema: vehicleTypes=%d services=%d roles=%q", vehicleTypes, services, roles)
	}
	rows, err := db.QueryContext(ctx, "SELECT id, nickname, u_role, pass_hash FROM users LIMIT 0")
	if err != nil {
		t.Fatal(err)
	}
	_ = rows.Close()
}

type disposableContainer struct {
	ID     string `json:"Id"`
	Name   string
	Config struct {
		Image  string
		Labels map[string]string
	}
	HostConfig struct {
		Tmpfs  map[string]string
		Binds  []string
		Mounts []json.RawMessage
	}
	Mounts []struct {
		Type        string
		Destination string
	}
	NetworkSettings struct {
		Ports map[string][]struct {
			HostIP   string `json:"HostIp"`
			HostPort string
		}
	}
}

func hasIsolatedPostgresMounts(container disposableContainer) bool {
	const dataDir = "/var/lib/postgresql/data"
	if _, ok := container.HostConfig.Tmpfs[dataDir]; !ok || len(container.HostConfig.Tmpfs) != 1 || len(container.HostConfig.Binds) != 0 {
		return false
	}
	// Docker may report --tmpfs only in HostConfig.Tmpfs, leaving Mounts empty.
	// Reject every volume/bind or unexpected target when Mounts is populated.
	for _, mount := range container.Mounts {
		if mount.Type != "tmpfs" || mount.Destination != dataDir {
			return false
		}
	}
	return len(container.Mounts) <= 1
}

func TestIsolatedPostgresMounts(t *testing.T) {
	for _, test := range []struct {
		name    string
		inspect string
		valid   bool
	}{
		{"tmpfs reported only in HostConfig", `{"HostConfig":{"Tmpfs":{"/var/lib/postgresql/data":"rw"},"Binds":null},"Mounts":[]}`, true},
		{"tmpfs also reported in Mounts", `{"HostConfig":{"Tmpfs":{"/var/lib/postgresql/data":"rw"}},"Mounts":[{"Type":"tmpfs","Destination":"/var/lib/postgresql/data"}]}`, true},
		{"missing tmpfs configuration", `{"Mounts":[{"Type":"tmpfs","Destination":"/var/lib/postgresql/data"}]}`, false},
		{"wrong tmpfs target", `{"HostConfig":{"Tmpfs":{"/var/lib/postgresql":"rw"}}}`, false},
		{"extra tmpfs target", `{"HostConfig":{"Tmpfs":{"/var/lib/postgresql/data":"rw","/extra":"rw"}}}`, false},
		{"host bind configured", `{"HostConfig":{"Tmpfs":{"/var/lib/postgresql/data":"rw"},"Binds":["/host:/extra"]}}`, false},
		{"anonymous or named volume", `{"HostConfig":{"Tmpfs":{"/var/lib/postgresql/data":"rw"}},"Mounts":[{"Type":"volume","Destination":"/var/lib/postgresql/data"}]}`, false},
		{"bind in Mounts", `{"HostConfig":{"Tmpfs":{"/var/lib/postgresql/data":"rw"}},"Mounts":[{"Type":"bind","Destination":"/extra"}]}`, false},
		{"unexpected mount target", `{"HostConfig":{"Tmpfs":{"/var/lib/postgresql/data":"rw"}},"Mounts":[{"Type":"tmpfs","Destination":"/extra"}]}`, false},
		{"duplicate mount entries", `{"HostConfig":{"Tmpfs":{"/var/lib/postgresql/data":"rw"}},"Mounts":[{"Type":"tmpfs","Destination":"/var/lib/postgresql/data"},{"Type":"tmpfs","Destination":"/var/lib/postgresql/data"}]}`, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			var container disposableContainer
			if err := json.Unmarshal([]byte(test.inspect), &container); err != nil {
				t.Fatal(err)
			}
			if got := hasIsolatedPostgresMounts(container); got != test.valid {
				t.Fatalf("mount isolation: got=%v want=%v", got, test.valid)
			}
		})
	}
}

func inspectContainer(ctx context.Context, run func(context.Context, ...string) (string, error), reference string) (disposableContainer, error) {
	output, err := run(ctx, "inspect", "--type", "container", reference)
	if err != nil {
		return disposableContainer{}, fmt.Errorf("inspect own disposable container: %w: %s", err, output)
	}
	var containers []disposableContainer
	if err := json.Unmarshal([]byte(output), &containers); err != nil {
		return disposableContainer{}, err
	}
	if len(containers) != 1 {
		return disposableContainer{}, fmt.Errorf("expected one disposable container, got %d", len(containers))
	}
	return containers[0], nil
}

func createBootstrapContainer(ctx context.Context, run func(context.Context, ...string) (string, error), name, token string, args []string, registerCleanup func(func() error)) (string, error) {
	var expectedID string
	// Register before creation: Docker may commit a container even if the client
	// receives EOF, startup fails, or stdout never delivers its ID.
	registerCleanup(func() error {
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		return cleanupBootstrapContainer(cleanupCtx, run, name, token, expectedID)
	})
	output, err := run(ctx, args...)
	if candidate := strings.TrimSpace(output); isDockerContainerID(candidate) {
		expectedID = candidate
	}
	if err != nil {
		return "", fmt.Errorf("Docker creation/startup failed: %w", err)
	}
	if expectedID == "" {
		return "", errors.New("Docker did not return a valid container ID")
	}
	return expectedID, nil
}

func cleanupBootstrapContainer(ctx context.Context, run func(context.Context, ...string) (string, error), name, token, expectedID string) error {
	owned, err := inspectContainer(ctx, run, name)
	if err != nil {
		// Verify absence independently instead of treating every inspect error as
		// "not created"; an unavailable daemon must remain a cleanup failure.
		listed, listErr := run(ctx, "container", "ls", "--all", "--filter", "name=^/"+name+"$", "--no-trunc", "--format", "{{.ID}}")
		if listErr == nil && strings.TrimSpace(listed) == "" {
			return nil
		}
		return fmt.Errorf("cannot verify own container %s for cleanup; no removal attempted: %w", name, errors.Join(err, listErr))
	}
	if owned.Name != "/"+name || owned.Config.Labels["autofort.disposable-bootstrap"] != token || !isDockerContainerID(owned.ID) || (expectedID != "" && owned.ID != expectedID) {
		return fmt.Errorf("ownership verification failed for %s; no removal attempted", name)
	}
	args := []string{"rm", "--force"}
	// Our creation command supplies no volume/bind mounts. With that still
	// confirmed, --volumes reclaims only image-created anonymous volumes if an
	// early failure occurs before the tmpfs gate. Explicit mounts are preserved;
	// Docker rm never removes named volumes and no volume prune is used.
	if len(owned.HostConfig.Binds) == 0 && len(owned.HostConfig.Mounts) == 0 {
		args = append(args, "--volumes")
	}
	args = append(args, owned.ID)
	if output, err := run(ctx, args...); err != nil {
		return fmt.Errorf("remove own disposable container %s: %w: %s", name, err, output)
	}
	return nil
}

func isDockerContainerID(id string) bool {
	return len(id) == 64 && strings.Trim(id, "0123456789abcdef") == ""
}

func dockerCommand(ctx context.Context, args ...string) (string, error) {
	// Keep image-pull progress on stderr separate from the returned container ID.
	output, err := exec.CommandContext(ctx, "docker", args...).Output()
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return string(output), fmt.Errorf("%w: %s", err, exitErr.Stderr)
	}
	return string(output), err
}

func isLocalDockerEndpoint(endpoint string) bool {
	return strings.HasPrefix(endpoint, "unix:///") || strings.HasPrefix(strings.ToLower(endpoint), "npipe:////./pipe/")
}

func TestLocalDockerEndpoint(t *testing.T) {
	for _, test := range []struct {
		endpoint string
		local    bool
	}{
		{"unix:///var/run/docker.sock", true},
		{"unix:///run/user/1000/docker.sock", true},
		{"npipe:////./pipe/dockerDesktopLinuxEngine", true},
		{"npipe:////remote/pipe/docker_engine", false},
		{"tcp://127.0.0.1:2375", false},
		{"ssh://remote", false},
		{"", false},
	} {
		if got := isLocalDockerEndpoint(test.endpoint); got != test.local {
			t.Errorf("endpoint %q: got local=%v, want %v", test.endpoint, got, test.local)
		}
	}
}
