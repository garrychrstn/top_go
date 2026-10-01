CREATE TABLE items (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    price NUMERIC NOT NULL
);

CREATE TABLE tx_rental (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid REFERENCES users(id),
    customer_id uuid REFERENCES customers(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE tx_rental_items (
    rental_id uuid REFERENCES tx_rental(id),
    item_id uuid REFERENCES items(id),
    status TEXT NOT NULL
);
