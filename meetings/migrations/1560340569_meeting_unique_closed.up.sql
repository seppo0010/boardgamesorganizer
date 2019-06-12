ALTER TABLE meetings DROP CONSTRAINT meetings_group_id_key;
CREATE UNIQUE INDEX meetings_group_id_key ON meetings (group_id) WHERE (not closed);
