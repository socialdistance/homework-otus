-- +goose Up
CREATE TABLE IF NOT EXISTS events
(
    id          uuid      not null,
    title       text      not null,
    started     timestamp not null,
    ended       timestamp not null,
    description text      not null default '',
    user_id     uuid      not null,
    notify      time      not null
);

-- +goose Down
DROP TABLE events;