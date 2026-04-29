-- +goose Up

CREATE TABLE IF NOT EXISTS vehicles (
    id UUID PRIMARY KEY,
    customer_id UUID NOT NULL REFERENCES customers(id),
    vehicle_type_id UUID NOT NULL REFERENCES vehicle_types(id),
    vin TEXT NOT NULL,
    gov_number TEXT NOT NULL,
    mileage INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (vin)
);