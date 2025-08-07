CREATE TABLE IF NOT EXISTS services.services(
    uuid VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    price INT4 NOT NULL
);

CREATE TABLE IF NOT EXISTS services.subscription(
    id SERIAL4 PRIMARY KEY,
    user_uuid VARCHAR(255) NOT NULL,
    service_uuid VARCHAR(255) NOT NULL,
    start_date DATE NOT NULL,
    stop_date DATE NOT NULL,
    CONSTRAINT fk_service FOREIGN KEY (service_uuid) REFERENCES services.services(uuid)
);