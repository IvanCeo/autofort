CREATE TABLE customers (
    id UUID PRIMARY KEY,

    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    phone_number TEXT NOT NULL UNIQUE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE vehicle_types (
    id UUID PRIMARY KEY,

    brand TEXT NOT NULL,
    model TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE (brand, model)
);

CREATE TABLE vehicles (
    id UUID PRIMARY KEY,

    customer_id UUID NOT NULL,
    vehicle_type_id UUID NOT NULL,

    license_plate TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    FOREIGN KEY (customer_id) REFERENCES customers(id),
    FOREIGN KEY (vehicle_type_id) REFERENCES vehicle_types(id),

    UNIQUE (license_plate)
);

CREATE TABLE services (
    id UUID PRIMARY KEY,

    name TEXT NOT NULL,
    description TEXT,

    is_active BOOLEAN NOT NULL DEFAULT true

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()

    UNIQUE (name)
);

CREATE TABLE parts (
    id UUID PRIMARY KEY,

    name TEXT NOT NULL,
    description TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE vehicle_type_services (
    vehicle_type_id UUID NOT NULL,
    service_id UUID NOT NULL,

    estimated_time_minutes INTEGER NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY (vehicle_type_id, service_id),

    FOREIGN KEY (vehicle_type_id) REFERENCES vehicle_types(id),
    FOREIGN KEY (service_id) REFERENCES services(id)
);

CREATE TABLE orders (
    id UUID PRIMARY KEY,

    customer_id UUID NOT NULL,
    vehicle_id UUID NOT NULL,

    status TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    FOREIGN KEY (customer_id) REFERENCES customers(id),
    FOREIGN KEY (vehicle_id) REFERENCES vehicles(id)
);

CREATE TABLE order_services (
    order_id UUID NOT NULL,
    service_id UUID NOT NULL,

    estimated_time_minutes INTEGER NOT NULL,
    rate_per_hour NUMERIC(10,2) NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY (order_id, service_id),

    FOREIGN KEY (order_id) REFERENCES orders(id),
    FOREIGN KEY (service_id) REFERENCES services(id)
);

CREATE TABLE order_parts (
    order_id UUID NOT NULL,
    part_id UUID NOT NULL,

    quantity INTEGER NOT NULL,
    price_per_unit NUMERIC(10,2) NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY (order_id, part_id),

    FOREIGN KEY (order_id) REFERENCES orders(id),
    FOREIGN KEY (part_id) REFERENCES parts(id)
);

CREATE TABLE labor_rates (
    id UUID PRIMARY KEY,

    rate_per_hour NUMERIC(10,2) NOT NULL,

    valid_from TIMESTAMPTZ NOT NULL,
    valid_to TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);