CREATE TABLE IF NOT EXISTS users (
    email character varying UNIQUE,
    password character varying
);