ALTER TABLE users
  ADD COLUMN recover_selector varchar(1024),
  ADD COLUMN recover_verifier varchar(1024),
  ADD COLUMN recover_expiry timestamp with time zone;
