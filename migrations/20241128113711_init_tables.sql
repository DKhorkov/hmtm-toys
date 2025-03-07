-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS masters
(
    id         SERIAL PRIMARY KEY,
    user_id    INTEGER   NOT NULL UNIQUE,
    info       TEXT      NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS categories
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS toys
(
    id          SERIAL PRIMARY KEY,
    master_id   INTEGER     NOT NULL,
    category_id INTEGER     NOT NULL,
    name        VARCHAR(50) NOT NULL,
    description TEXT        NOT NULL,
    price       FLOAT       NOT NULL,
    quantity    INTEGER     NOT NULL,
    created_at  TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (master_id) REFERENCES masters (id),
    FOREIGN KEY (category_id) REFERENCES categories (id)
);

CREATE TABLE IF NOT EXISTS tags
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(50) NOT NULL UNIQUE,
    created_at TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS toys_tags_associations
(
    id         SERIAL PRIMARY KEY,
    toy_id     INTEGER   NOT NULL,
    tag_id     INTEGER   NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (toy_id) REFERENCES toys (id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS toys_attachments
(
    id         SERIAL PRIMARY KEY,
    toy_id     INTEGER   NOT NULL,
    link       VARCHAR   NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (toy_id) REFERENCES toys (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS toys_tags_associations;
DROP TABLE IF EXISTS toys_attachments;
DROP TABLE IF EXISTS toys;
DROP TABLE IF EXISTS masters;
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS categories;
-- +goose StatementEnd
