CREATE TABLE events (
  id  VARCHAR(64) NOT NULL,
  type VARCHAR(32) NOT NULL,
  created_at BIGINT NOT NULL,
  payload BYTEA ,
  mapper VARCHAR(32) NOT NULL
);
CREATE INDEX id_idx ON events (id);