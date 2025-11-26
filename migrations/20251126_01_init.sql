-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS campaigns (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    template TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    scheduled_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS campaign_messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    campaign_id UUID REFERENCES campaigns(id) ON DELETE CASCADE,
    user_id TEXT NOT NULL,
    user_email TEXT NOT NULL,
    user_data JSONB NOT NULL,
    status TEXT DEFAULT 'pending',
    sent_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id TEXT UNIQUE NOT NULL,
    first_name TEXT NOT NULL,
    email TEXT NOT NULL,
    discount_percent INT DEFAULT 10,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Seed some real users
INSERT INTO users (user_id, first_name, email, discount_percent) VALUES
('alice', 'Alice Wonder', 'alice@company.com', 25),
('bob', 'Bob Builder', 'bob@company.com', 35),
('charlie', 'Charlie Chocolate', 'charlie@company.com', 15),
('ngeno', 'Ngeno', 'ngeno@company.com', 100),
('david', 'David Goliath', 'david@company.com', 40)
ON CONFLICT (user_id) DO NOTHING;
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS campaign_messages;
DROP TABLE IF EXISTS campaigns;