-- +goose Down

ALTER TABLE tps_checkins
    ALTER COLUMN qr_id SET NOT NULL;
