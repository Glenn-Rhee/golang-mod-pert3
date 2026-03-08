-- Buat Database name 
CREATE DATABASE TokoLepkom_npm;

-- Buat Table Products nya
CREATE TABLE products (
    id      VARCHAR(100) PRIMARY KEY,
    name    VARCHAR(255) NOT NULL,
    price   DECIMAL(15,2) NOT NULL
);

-- Buat Table product_details
CREATE TABLE product_details (
    product_id  VARCHAR(100) PRIMARY KEY,
    stock       INT NOT NULL,
    is_active   BOOLEAN NOT NULL DEFAULT TRUE,
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    image       LONGBLOB,
    
    CONSTRAINT fk_product_details_product
        FOREIGN KEY (product_id) REFERENCES products(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);