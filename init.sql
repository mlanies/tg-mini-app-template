CREATE TABLE users (
    id BIGINT PRIMARY KEY,
    username TEXT,
    first_name TEXT,
    last_name TEXT,
    language_code TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT unique_user UNIQUE (id, username)
);

CREATE TABLE appointments (
    id SERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    service_name TEXT,
    appointment_time TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT unique_appointment UNIQUE (user_id, appointment_time)
);

CREATE TABLE services (
    id SERIAL PRIMARY KEY,
    name TEXT,
    price NUMERIC(10, 2),
    created_at TIMESTAMP DEFAULT NOW()
);

-- Индексы
CREATE INDEX idx_appointments_user_id ON appointments(user_id);
CREATE INDEX idx_appointments_appointment_time ON appointments(appointment_time);
CREATE INDEX idx_services_name ON services(name);
