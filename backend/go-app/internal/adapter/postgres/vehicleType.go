package postgres

import (
	"autofort/internal/entity"
	e "autofort/internal/errors"
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (p *Postgres) ListVehicleTypes(ctx context.Context) ([]*entity.VehicleType, error) {
	const q = `
		SELECT id, brand, model
		FROM vehicle_types
		ORDER BY brand ASC, model ASC
	`

	rows, err := p.pool.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*entity.VehicleType, 0, 32)

	for rows.Next() {
		vt := new(entity.VehicleType)
		if err := rows.Scan(&vt.ID, &vt.Brand, &vt.Model); err != nil {
			return nil, err
		}
		items = append(items, vt)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (p *Postgres) GetVehicleTypeByID(ctx context.Context, id uuid.UUID) (*entity.VehicleType, error) {
	const q = `
		SELECT id, brand, model
		FROM vehicle_types
		WHERE id = $1
	`

	row := p.pool.QueryRow(ctx, q, id)

	vt := new(entity.VehicleType)
	if err := row.Scan(&vt.ID, &vt.Brand, &vt.Model); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, e.ErrNoVehicleType
		}
		return nil, err
	}

	return vt, nil
}

func (p *Postgres) AddVehicleType(ctx context.Context, vt *entity.VehicleType) error {
	const q = `
		INSERT INTO vehicle_types (id, brand, model)
		VALUES ($1, $2, $3)
		ON CONFLICT (brand, model) DO NOTHING
	`

	// используем vt.ID, а не генерим новый
	_, err := p.pool.Exec(ctx, q, vt.ID, vt.Brand, vt.Model)
	return err
}
