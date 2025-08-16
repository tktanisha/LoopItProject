-- ----- USER -----
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    full_name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    phone_number VARCHAR(20),
    address VARCHAR(255),
    password_hash VARCHAR(255) NOT NULL,
    society_id INT NOT NULL,
    role VARCHAR(20) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- ----- SOCIETY -----
CREATE TABLE IF NOT EXISTS societies (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    location VARCHAR(255),
    pincode VARCHAR(20)
);

-- ----- LENDER -----
CREATE TABLE IF NOT EXISTS lenders (
    id SERIAL PRIMARY KEY,
    is_verified BOOLEAN DEFAULT FALSE,
    total_earnings NUMERIC(12,2) DEFAULT 0
);

-- ----- CATEGORY -----
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    price NUMERIC(12,2) NOT NULL,
    security NUMERIC(12,2) NOT NULL
);

-- ----- PRODUCT -----
CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    lender_id INT NOT NULL,
    category_id INT NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    duration INT NOT NULL,
    is_available BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (lender_id) REFERENCES lenders(id),
    FOREIGN KEY (category_id) REFERENCES categories(id)
);

-- ----- PRODUCT IMAGE -----
CREATE TABLE IF NOT EXISTS product_images (
    id SERIAL PRIMARY KEY,
    product_id INT NOT NULL,
    image_url TEXT NOT NULL,
    uploaded_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (product_id) REFERENCES products(id)
);

-- ----- FEEDBACK -----
CREATE TABLE IF NOT EXISTS feedbacks (
    id SERIAL PRIMARY KEY,
    given_by INT NOT NULL,
    given_to INT NOT NULL,
    text TEXT,
    rating INT CHECK(rating >= 0 AND rating <= 5),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (given_by) REFERENCES users(id),
    FOREIGN KEY (given_to) REFERENCES users(id)
);

-- ----- ORDER -----
CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    product_id INT NOT NULL,
    user_id INT NOT NULL,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    total_amount NUMERIC(12,2) NOT NULL,
    security_amount NUMERIC(12,2) NOT NULL,
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (product_id) REFERENCES products(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- ----- BUYING REQUEST -----
CREATE TABLE IF NOT EXISTS buying_requests (
    id SERIAL PRIMARY KEY,
    product_id INT NOT NULL,
    requested_by INT NOT NULL,
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (product_id) REFERENCES products(id),
    FOREIGN KEY (requested_by) REFERENCES users(id)
);

-- ----- RETURN REQUEST -----
CREATE TABLE IF NOT EXISTS return_requests (
    id SERIAL PRIMARY KEY,
    order_id INT NOT NULL,
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (order_id) REFERENCES orders(id)
);
