
-- +migrate Up
ALTER TABLE `team` ADD `mix_team` tinyint(1) DEFAULT NULL ;

-- +migrate Down

ALTER TABLE `team` DROP `mix_team` ;