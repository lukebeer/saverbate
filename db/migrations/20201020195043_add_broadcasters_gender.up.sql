BEGIN;
CREATE TYPE broadcaster_gender AS ENUM (
  'female',
  'male',
  'trans',
  'couple'
);
ALTER TABLE broadcasters ADD COLUMN gender broadcaster_gender NOT NULL DEFAULT 'female';
COMMIT;
