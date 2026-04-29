package postgres

import (
	"errors"

	"autofort/internal/entity"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (p *Postgres) ListVehicleTypes() ([]*entity.VehicleType, error) {
	const q = `
		SELECT id, brand, model
		FROM vehicle_types
		ORDER BY brand ASC, model ASC
	`

	rows, err := p.pool.Query(p.ctx, q)
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

func (p *Postgres) GetVehicleTypeByID(id uuid.UUID) (*entity.VehicleType, error) {
	const q = `
		SELECT id, brand, model
		FROM vehicle_types
		WHERE id = $1
	`

	row := p.pool.QueryRow(p.ctx, q, id)

	vt := new(entity.VehicleType)
	if err := row.Scan(&vt.ID, &vt.Brand, &vt.Model); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("vehicle type not found")
		}
		return nil, err
	}

	return vt, nil
}

func (p *Postgres) AddVehicleType(vt *entity.VehicleType) error {
	const q = `
		INSERT INTO vehicle_types (id, brand, model)
		VALUES ($1, $2, $3)
		ON CONFLICT (brand, model) DO NOTHING
	`

	// используем vt.ID, а не генерим новый
	_, err := p.pool.Exec(p.ctx, q, vt.ID, vt.Brand, vt.Model)
	return err
}
