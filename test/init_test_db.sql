-- init_test_db.sql

-- Create the necessary table
CREATE TABLE IF NOT EXISTS adverts (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    photos JSONB,
    price DECIMAL(10, 2) NOT NULL
);

-- Insert test data
INSERT INTO adverts (title, description, photos, price)
VALUES
    ('Test Advert 1', 'Description 1', '["photo1.jpg", "photo2.jpg"]', 100.00),
    ('Test Advert 2', 'Description 2', '["photo3.jpg", "photo4.jpg"]', 150.00),
    ('Test Advert 3', 'Description 3', '["photo5.jpg", "photo6.jpg"]', 200.00);
