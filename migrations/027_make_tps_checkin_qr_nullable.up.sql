-- +goose Up
-- Allow TPS check-in without QR record (manual/registration token)

ALTER TABLE tps_checkins
    ALTER COLUMN qr_id DROP NOT NULL;
