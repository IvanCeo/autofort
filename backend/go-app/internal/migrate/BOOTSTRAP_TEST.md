# Disposable PostgreSQL 16 bootstrap check

From `backend/go-app`, with Go 1.27.1 and a running local Docker daemon:

```powershell
$env:AUTOFORT_DISPOSABLE_PG_TEST = '1'
try {
    go test ./internal/migrate -run '^TestDisposablePostgresBootstrap$' -count=1 -v -timeout=5m
    $bootstrapTestExit = $LASTEXITCODE
} finally {
    Remove-Item Env:AUTOFORT_DISPOSABLE_PG_TEST
}
exit $bootstrapTestExit
```

This opt-in check creates its own uniquely named `postgres:16` container with
generated disposable credentials, tmpfs database storage, and a random loopback
port. It accepts no DSN and reads no `.env`. Before migrations, it verifies the
container identity, ownership label, image, tmpfs-only mounts, loopback binding,
PostgreSQL 16 server, database/user identity, and empty public schema. It removes
only that verified container on completion. Cleanup is registered before Docker
creation, so client EOF, startup failure, and missing or malformed container IDs
still trigger inspection by the unique name and ownership label. A captured ID
must match too. Cleanup uses a separate timeout even if creation was cancelled.
No Compose services or persistent volumes are requested. Docker may download
`postgres:16` if it is not cached.
The effective Docker endpoint must be a local Unix socket or Windows named pipe;
remote contexts and `DOCKER_HOST` endpoints are rejected before container creation.
Removal also reclaims image-created anonymous volumes when inspection confirms
there are no explicit bind/volume mounts. Explicit mounted volumes are preserved;
named volumes are never deleted, and no volume prune is used. If ownership cannot
be verified or Docker cleanup fails, the test reports an error instead of deleting
unverified resources.

The check reproduces the historical enum syntax error in a rolled-back
transaction, then executes the real migrations 1-7 through `migrate.Up`. It checks
Goose version 7, the users schema/enum, and seed counts; repeats Up; reconnects and
repeats Up again; checks unchanged migration history; rolls back migration 7 to
version 6; and reapplies it. It verifies startup migration reconnect/reapply
behavior, not a full HTTP application restart or authentication flow.

Without the opt-in variable, the integration test is skipped. A skipped test is
not evidence of database bootstrap success. With opt-in, Docker unavailability,
isolation failures, SQL failures, and cleanup failures fail the test.

Migration 7 Down drops `users` before its enum without `CASCADE`. Down is
destructive to users data and is exercised only in this disposable database.
Migrations 1-6 and existing deployment configuration are unchanged.
