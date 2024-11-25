CREATE DATABASE food_delivery_db;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE role AS ENUM ('admin', 'employee', 'customer');

CREATE TYPE gender AS ENUM ('male', 'female');

CREATE TYPE transaction_type AS ENUM ('credit', 'debit'); 

CREATE TYPE menu_type AS ENUM ('main dish', 'side dish', 'dessert', 'beverage');

CREATE TYPE unit_type AS ENUM ('piece', 'portion', 'packet', 'cup')

CREATE TYPE order_status AS ENUM ('preparing', 'out for delivery', 'delivered');

CREATE TABLE users(
  id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  username VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL,
  password VARCHAR(255) NOT NULL,
  role role NOT NULL,
  gender gender NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP
);

ALTER TABLE users ADD CONSTRAINT unique_user_username UNIQUE (username);
ALTER TABLE users ADD CONSTRAINT unique_user_email UNIQUE (email);

CREATE TABLE token_blacklists (
  id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  token VARCHAR(512) NOT NULL,
  expires_at TIMESTAMP NOT NULL
);

create table menus(
  id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  type menu_type NOT NULL,
  description TEXT NOT NULL,
  unit_type unit_type NOT NULL,
  price DOUBLE PRECISION NOT NULL,
  created_by uuid NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP,
  FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE
);

ALTER TABLE menus ADD CONSTRAINT unique_menu_name UNIQUE (name);
ALTER TABLE menus ADD CONSTRAINT unique_menu_description UNIQUE (description);

create table balances(
  id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  customer_id uuid NOT NULL,
  transaction_type transaction_type NOT NULL,
  amount DOUBLE PRECISION NOT NULL,
  description TEXT NOT NULL,
  balance DOUBLE PRECISION NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (customer_id) REFERENCES users(id) ON DELETE CASCADE
);

create table promos(
  id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  employee_id uuid,
  promo_code VARCHAR(20) UNIQUE NOT NULL,
  discount NUMERIC(5, 2) NOT NULL,
  is_percentage BOOLEAN NOT NULL DEFAULT TRUE,
  start_date DATE,
  end_date DATE,
  description TEXT,
  FOREIGN KEY (employee_id) REFERENCES users(id) ON DELETE CASCADE
);

create table orders(
  id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  customer_id uuid NOT NULL,
  address TEXT NOT NULL,
  promo_code VARCHAR(255),
  promo_used BOOLEAN DEFAULT FALSE,
  order_status order_status NOT NULL,
  note TEXT,
  date DATE,
  total_price DOUBLE PRECISION NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (customer_id) REFERENCES users(id) ON DELETE CASCADE
);

create table order_items(
  id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  order_id uuid NOT NULL,
  menu_id uuid NOT NULL,
  quantity int NOT NULL,
  FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
  FOREIGN KEY (menu_id) REFERENCES menus(id) ON DELETE CASCADE
);

create table reviews(
  id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
  customer_id uuid NOT NULL,
  menu_id uuid NOT NULL,
  order_id uuid NOT NULL,
  rating INT NOT null,
  comment TEXT NOT NULL,
  buy_date TIMESTAMP,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP,
  FOREIGN KEY (customer_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (menu_id) REFERENCES menus(id) ON DELETE CASCADE,
  FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE
);

ALTER TABLE reviews
ADD CONSTRAINT unique_review_per_purchase UNIQUE (customer_id, menu_id, order_id);
