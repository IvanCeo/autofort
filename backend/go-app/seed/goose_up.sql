-- +goose Up

CREATE TABLE IF NOT EXISTS customers (
    id UUID PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    phone_number TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS vehicle_types (
    id UUID PRIMARY KEY,
    brand TEXT NOT NULL,
    model TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (brand, model)
);

CREATE TABLE IF NOT EXISTS vehicles (
    id UUID PRIMARY KEY,
    customer_id UUID NOT NULL REFERENCES customers(id),
    vehicle_type_id UUID NOT NULL REFERENCES vehicle_types(id),
    license_plate TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (license_plate)
);

CREATE TABLE IF NOT EXISTS parts (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS services (
    id UUID PRIMARY KEY,
    category TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    base_minutes INT NOT NULL CHECK (base_minutes > 0),
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (category, name)
);

CREATE TABLE IF NOT EXISTS vehicle_type_multipliers (
    vehicle_type_id UUID PRIMARY KEY REFERENCES vehicle_types(id),
    multiplier NUMERIC(4,2) NOT NULL DEFAULT 1.00 CHECK (multiplier > 0)
);

CREATE TABLE IF NOT EXISTS vehicle_type_service_overrides (
    vehicle_type_id UUID NOT NULL REFERENCES vehicle_types(id),
    service_id UUID NOT NULL REFERENCES services(id),
    minutes_override INT NOT NULL CHECK (minutes_override > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (vehicle_type_id, service_id)
);

CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY,
    customer_id UUID NOT NULL REFERENCES customers(id),
    vehicle_id UUID NOT NULL REFERENCES vehicles(id),
    status TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS order_services (
    order_id UUID NOT NULL REFERENCES orders(id),
    service_id UUID NOT NULL REFERENCES services(id),

    estimated_time_minutes INT NOT NULL CHECK (estimated_time_minutes > 0),
    actual_time_minutes INT CHECK (actual_time_minutes > 0),
    notes TEXT,

    rate_per_hour NUMERIC(10,2) NOT NULL CHECK (rate_per_hour > 0),

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY (order_id, service_id)
);

CREATE TABLE IF NOT EXISTS order_parts (
    order_id UUID NOT NULL REFERENCES orders(id),
    part_id UUID NOT NULL REFERENCES parts(id),

    quantity INT NOT NULL CHECK (quantity > 0),
    price_per_unit NUMERIC(10,2) NOT NULL CHECK (price_per_unit > 0),

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY (order_id, part_id)
);

CREATE TABLE IF NOT EXISTS labor_rates (
    id UUID PRIMARY KEY,
    rate_per_hour NUMERIC(10,2) NOT NULL CHECK (rate_per_hour > 0),
    valid_from TIMESTAMPTZ NOT NULL,
    valid_to TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CHECK (valid_to IS NULL OR valid_to > valid_from)
);

-- +goose Down

DROP TABLE IF EXISTS labor_rates;
DROP TABLE IF EXISTS order_parts;
DROP TABLE IF EXISTS order_services;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS vehicle_type_service_overrides;
DROP TABLE IF EXISTS vehicle_type_multipliers;
DROP TABLE IF EXISTS services;
DROP TABLE IF EXISTS parts;
DROP TABLE IF EXISTS vehicles;
DROP TABLE IF EXISTS vehicle_types;
DROP TABLE IF EXISTS customers;
