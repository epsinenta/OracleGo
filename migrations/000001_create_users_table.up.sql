CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    email character varying UNIQUE,
    password character varying
);