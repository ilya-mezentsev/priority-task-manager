-- noinspection SqlNoDataSourceInspectionForFile

DROP TABLE IF EXISTS account;
CREATE TABLE account(
    id SERIAL PRIMARY KEY,
    hash VARCHAR(32) UNIQUE NOT NULL,
    role VARCHAR(42) NOT NULL
);

INSERT INTO account(hash, role)
VALUES
    ('86d3f3a95c324c9479bd8986968f4327', 'role-1'),
    ('11c9beec53034beb3a6687891c9e248a', 'role-2'),
    ('eafae37058a254f5dfdfe22ede8cca1f', 'role-3');

DROP TABLE IF EXISTS task_stat;
CREATE TABLE task_stat(
    id SERIAL PRIMARY KEY,
    uuid VARCHAR(36) UNIQUE NOT NULL,
    account_hash VARCHAR(32) REFERENCES account(hash),
    status VARCHAR(42) NOT NULL,
    type VARCHAR(42) NOT NULL,
    added_to_queue TIMESTAMP NOT NULL,
    extracted_from_queue TIMESTAMP DEFAULT NULL,
    completed TIMESTAMP DEFAULT NULL
);
