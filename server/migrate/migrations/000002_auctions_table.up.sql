CREATE TABLE auctions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR NOT NULL,
    description TEXT,
    starting_price NUMERIC NOT NULL,
    current_price NUMERIC NOT NULL,
    type VARCHAR NOT NULL CHECK (type IN ('english', 'dutch', 'sealed')),
    status VARCHAR NOT NULL CHECK (status IN ('open', 'closed')),
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    seller_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (seller_id) REFERENCES users(id)
);