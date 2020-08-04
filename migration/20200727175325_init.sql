-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE `chat_rooms` (
  `chat_id` bigint(20) DEFAULT NULL,
  `group_name` varchar(30) DEFAULT NULL,
  `user_id` bigint(20) DEFAULT NULL,
  `user_name` varchar(30) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
