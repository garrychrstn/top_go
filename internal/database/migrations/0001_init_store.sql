CREATE TABLE users (
    id  uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    username text NOT NULL unique,
    password    text not null
);

CREATE TABLE customers (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name    TEXT NOT NULL,
    phone_number TEXT NOT NULL,
    address TEXT
);
