-- +goose Up

CREATE TABLE IF NOT EXISTS vehicle_types (
    id UUID PRIMARY KEY,
    brand TEXT NOT NULL,
    model TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (brand, model)
);
