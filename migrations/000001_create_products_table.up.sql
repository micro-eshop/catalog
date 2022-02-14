CREATE TABLE IF NOT EXISTS products (
  id INTEGER PRIMARY KEY UNIQUE NOT NULL,
  name TEXT NOT NULL,
  description TEXT,
  price DOUBLE PRECISION NOT NULL,
  brand TEXT,
  promotion_price DOUBLE PRECISION NULL
);