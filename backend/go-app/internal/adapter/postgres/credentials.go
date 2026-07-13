package postgres

import (
	er "autofort/internal/errors"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

// deprecated
func (p *Postgres) SignUp(ctx context.Context, nickname, role, passhash string) error {
	_, err := p.pool.Exec(
		ctx,
		`insert into users (nickname, u_role, passhash) values ($1, $2, $3)`,
		nickname,
		role,
		passhash,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return er.ErrUserAlreadyExists
		}
		return fmt.Errorf("pgx SignUp: %w", err)
	}

	return nil
}

func (p *Postgres) GetPassHash(
	ctx context.Context,
	userID *uuid.UUID,
	nickname string,
	hash *string,
) error {
	row := p.pool.QueryRow(
		ctx,
		`select id, pass_hash
		from users
		where nickname = $1`,
		nickname,
	)

	return row.Scan(userID, hash)
}
