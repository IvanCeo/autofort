package postgres

import (
	e "autofort/internal/errors"
	"errors"
	"strings"

	"autofort/internal/entity"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (p *Postgres) AddCustomer(c *entity.Customer) error {
	const q = `
		INSERT INTO customers (id, first_name, last_name, phone_number)
		VALUES ($1, $2, $3, $4)
	`
	_, err := p.pool.Exec(p.ctx, q, c.ID, c.FirstName, c.LastName, c.PhoneNumber)
	return err
}

func (p *Postgres) GetCustomer(id uuid.UUID) (*entity.Customer, error) {
	const q = `
		SELECT id, first_name, last_name, phone_number
		FROM customers
		WHERE id = $1
	`
	row := p.pool.QueryRow(p.ctx, q, id)

	var c entity.Customer
	if err := row.Scan(&c.ID, &c.FirstName, &c.LastName, &c.PhoneNumber); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, e.ErrNoCustomer
		}
		return nil, err
	}
	return &c, nil
}

func (p *Postgres) UpdateCustomer(c *entity.Customer) error {
	const q = `
		UPDATE customers
		SET first_name = $2,
		    last_name = $3,
		    phone_number = $4,
		    updated_at = now()
		WHERE id = $1
	`
	ct, err := p.pool.Exec(p.ctx, q, c.ID, c.FirstName, c.LastName, c.PhoneNumber)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return e.ErrNoCustomer
	}
	return nil
}

func (p *Postgres) SearchCustomers(q string, limit, offset int) ([]*entity.Customer, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	q = strings.TrimSpace(q)
	pattern := "%" + q + "%"

	const sql = `
		SELECT id, first_name, last_name, phone_number
		FROM customers
		WHERE ($1 = '' OR first_name ILIKE $2 OR last_name ILIKE $2 OR phone_number ILIKE $2)
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := p.pool.Query(p.ctx, sql, q, pattern, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*entity.Customer, 0, limit)
	for rows.Next() {
		c := new(entity.Customer)
		if err := rows.Scan(&c.ID, &c.FirstName, &c.LastName, &c.PhoneNumber); err != nil {
			return nil, err
		}
		items = append(items, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
