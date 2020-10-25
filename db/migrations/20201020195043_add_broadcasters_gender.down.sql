BEGIN;
ALTER TABLE broadcasters DROP COLUMN gender;
DROP TYPE broadcaster_gender;
COMMIT;
