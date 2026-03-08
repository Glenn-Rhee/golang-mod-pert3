-- Buat Database name 
CREATE DATABASE TokoLepkom_npm;

-- Buat Table Products nya
CREATE TABLE products (
    id      VARCHAR(100) PRIMARY KEY,
    name    VARCHAR(255) NOT NULL,
    price   DECIMAL(15,2) NOT NULL
);