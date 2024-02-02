CREATE TABLE IF NOT EXISTS History (
    id SERIAL PRIMARY KEY,
    time TIMESTAMP NOT NULL,
    sender VARCHAR(255) NOT NULL,
    recipient VARCHAR(255) NOT NULL,
    amount DOUBLE PRECISION NOT NULL
);

CREATE INDEX idx_sender_recipient ON History (sender, recipient);