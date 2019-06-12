DROP INDEX meetings_group_id_key;
ALTER TABLE meetings ADD CONSTRAINT meetings_group_id_key UNIQUE (group_id);
