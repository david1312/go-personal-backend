-- Adminer 4.8.1 MySQL 8.0.27 dump

SET NAMES utf8;
SET time_zone = '+00:00';
SET foreign_key_checks = 0;
SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO';

SET NAMES utf8mb4;

INSERT INTO `customers` (`id`, `name`, `birthdate` ,`uid`, `password`, `email`, `email_verified_token`, `email_verified_at`, `gender`, `is_active`, `phone`, `phone_verified_at`, `avatar`, `created_at`, `updated_at`, `deleted_at`) VALUES
(1,	'Loid Forger', '1996-12-13', 'a08d620f-0c3c-44d1-a054-df4bb0e96d30',	'$2a$10$2s3xRCm1xJMYW7Pe3G5tBuGsqioRoKyOCZBsu/t8Vn7xiSTvLxIGG',	'tester123@gmail.com',	NULL,	'2022-05-07 16:22:50',	'LAKI-LAKI',	1,	'081217852333',	'2022-05-07 16:22:50',	NULL,	'2022-05-07 16:22:50',	'2022-05-07 16:22:50',	NULL),
(2,	'David Bernadi', '1996-12-25', '93c996d5-4aed-403f-b19f-ee6fae02a7c3',	'$2a$10$62nS3UJdS8QIB9He5z8plOZRD.PtbEyHl77xjL6gUy/DYYdpVt0xO',	'davidbernadi13@gmail.com',	'b172f7edf6e98a43e28ea6b5d5fc3c50a722e7dc20d710bb30c18061a177960b',	'2022-05-08 08:05:15',	NULL,	1,	NULL,	NULL,	NULL,	'2022-05-08 06:01:28',	'2022-05-08 08:05:15',	NULL);
-- 2022-05-08 08:05:39