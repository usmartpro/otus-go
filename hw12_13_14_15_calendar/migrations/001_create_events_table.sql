-- +goose Up
CREATE TABLE IF NOT EXISTS events
(
    id                  UUID            NOT NULL,
    user_id             UUID            NOT NULL,
    title               VARCHAR(1024)   NOT NULL,
    started_at          timestamp       NOT NULL,
    finished_at         timestamp       NOT NULL,
    description         TEXT,
    notify_before_time  INT             NOT NULL default 0
);

-- +goose Down
DROP TABLE events;