package postgres

import (
	"autofort/internal/entity"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (p *Postgres) AddVehicle(v *entity.Vehicle) error {
	const q = `
		INSERT INTO vehicles (id, customer_id, vehicle_type_id, vin, gov_number, mileage)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := p.pool.Exec(
		p.ctx,
		q,
		v.ID,
		v.CustomerId,
		v.VehicleType,
		v.Vin,
		v.GovNumber,
		v.Mileage,
	)

	return err
}

func (p *Postgres) ListVehiclesByCustomerIDs(ids []uuid.UUID) ([]*entity.Vehicle, error) {
	if len(ids) == 0 {
		return []*entity.Vehicle{}, nil
	}

	placeholders := make([]string, 0, len(ids))
	args := make([]any, 0, len(ids))

	for i, id := range ids {
		placeholders = append(placeholders, fmt.Sprintf("$%d", i+1))
		args = append(args, id)
	}

	q := fmt.Sprintf(`
		SELECT
			id,
			customer_id,
			vehicle_type_id,
			vin,
			gov_number,
			mileage
		FROM vehicles
		WHERE customer_id IN (%s)
		ORDER BY created_at DESC
	`, strings.Join(placeholders, ", "))

	rows, err := p.pool.Query(p.ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*entity.Vehicle, 0, 32)
	for rows.Next() {
		v := new(entity.Vehicle)
		if err := rows.Scan(
			&v.ID,
			&v.CustomerId,
			&v.VehicleType,
			&v.Vin,
			&v.GovNumber,
			&v.Mileage,
		); err != nil {
			return nil, err
		}
		items = append(items, v)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (p *Postgres) GetVehicle(id uuid.UUID) (*entity.Vehicle, error) {
	const q = `
		SELECT id, customer_id, vehicle_type_id, vin, gov_number, mileage
		FROM vehicles
		WHERE id = $1
	`

	row := p.pool.QueryRow(p.ctx, q, id)

	v := new(entity.Vehicle)
	if err := row.Scan(
		&v.ID,
		&v.CustomerId,
		&v.VehicleType,
		&v.Vin,
		&v.GovNumber,
		&v.Mileage,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("vehicle not found")
		}
		return nil, err
	}

	return v, nil
}

func (p *Postgres) UpdateVehicle(v *entity.Vehicle) error {
	const q = `
		UPDATE vehicles
		SET vin = $2,
		    gov_number = $3,
		    mileage = $4,
		    updated_at = now()
		WHERE id = $1
	`

	ct, err := p.pool.Exec(p.ctx, q, v.ID, v.Vin, v.GovNumber, v.Mileage)
	if err != nil {
		return err
	}

	if ct.RowsAffected() == 0 {
		return errors.New("vehicle not found")
	}

	return nil
}
