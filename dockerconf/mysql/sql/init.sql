# ************************************************************
# Sequel Pro SQL dump
# Version 5446
#
# https://www.sequelpro.com/
# https://github.com/sequelpro/sequelpro
#
# Host: 127.0.0.1 (MySQL 8.0.16)
# Database: whispering_conner
# Generation Time: 2021-12-07 10:23:33 +0000
# ************************************************************


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
SET NAMES utf8mb4;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


# Dump of table social_accounts
# ------------------------------------------------------------

DROP TABLE IF EXISTS `social_accounts`;

CREATE TABLE `social_accounts` (
`id` bigint(64) NOT NULL AUTO_INCREMENT,
`social_id` varchar(64) COLLATE utf8mb4_general_ci NOT NULL,
`user_id` bigint(20) NOT NULL,
`type` tinyint(1) NOT NULL DEFAULT '0',
`created_date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY (`id`),
UNIQUE KEY `idx_social_id` (`social_id`),
UNIQUE KEY `idx_user_id` (`user_id`),
KEY `idx_date` (`created_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;



# Dump of table users
# ------------------------------------------------------------

DROP TABLE IF EXISTS `users`;

CREATE TABLE `users` (
`id` bigint(64) NOT NULL AUTO_INCREMENT,
`user_id` bigint(64) NOT NULL,
`account` varchar(20) COLLATE utf8mb4_general_ci NOT NULL,
`password` varchar(100) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
`status` tinyint(1) NOT NULL DEFAULT '0',
`name` varchar(32) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
`email` varchar(64) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
`remark` varchar(64) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
`created_date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY (`id`),
UNIQUE KEY `idx_account` (`account`),
UNIQUE KEY `idx_user_id` (`user_id`),
KEY `idx_date` (`created_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


# Dump of table cocktail_ingredients
# ------------------------------------------------------------

CREATE TABLE `cocktail_ingredients` (
`id` bigint(64) NOT NULL AUTO_INCREMENT,
`ingredient_id` bigint(64) NOT NULL,
`cocktail_id` bigint(64) NOT NULL,
`ingredient_name` varchar(16) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT ' 成分名稱',
`ingredient_amount` float NOT NULL DEFAULT '0' COMMENT '成分數量',
`ingredient_unit` varchar(16) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT ' 成分單位',
`created_date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY (`id`),
UNIQUE KEY `idx_ingredient_id` (`ingredient_id`),
UNIQUE KEY `idx_cocktail_id` (`cocktail_id`),
KEY `idx_date` (`created_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


# Dump of table cocktail_steps
# ------------------------------------------------------------

CREATE TABLE `cocktail_steps` (
`id` bigint(64) NOT NULL AUTO_INCREMENT,
`step_id` bigint(64) NOT NULL,
`cocktail_id` bigint(64) NOT NULL,
`step_number` int(2) unsigned NOT NULL DEFAULT '1' COMMENT ' 步驟',
`step_description` varchar(64) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT ' 步驟介紹',
`created_date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
PRIMARY KEY (`id`),
UNIQUE KEY `idx_recipe_id` (`step_id`),
UNIQUE KEY `idx_cocktail_id` (`cocktail_id`),
KEY `idx_date` (`created_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


# Dump of table cocktails
# ------------------------------------------------------------

CREATE TABLE `cocktails` (
 `id` bigint(64) NOT NULL AUTO_INCREMENT,
 `cocktail_id` bigint(64) NOT NULL,
 `title` varchar(16) COLLATE utf8mb4_general_ci NOT NULL COMMENT ' 調酒名稱',
 `description` varchar(512) COLLATE utf8mb4_general_ci NOT NULL COMMENT ' 調酒介紹',
 `user_id` bigint(64) NOT NULL COMMENT ' 作者id',
 `user_name` varchar(32) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT ' 作者名稱',
 `created_date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
 PRIMARY KEY (`id`),
 UNIQUE KEY `idx_cocktail_id` (`cocktail_id`),
 KEY `idx_user_id` (`user_id`),
 KEY `idx_date` (`created_date`)
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

LOCK TABLES `cocktails` WRITE;
/*!40000 ALTER TABLE `cocktails` DISABLE KEYS */;

INSERT INTO `cocktails` (`id`, `cocktail_id`, `title`, `description`, `user_id`, `user_name`, `created_date`)
VALUES
(1, 123456, 'Side Car', 'Good to drink', 1, 'static/my_image01.jpg', '2021-01-15 18:38:30'),
(2, 1111111, 'Old Fashion', 'Good to drink', 1, 'static/my_image01.jpg', '2021-02-15 18:38:30'),
(3, 222222, 'Gin tonic', 'Good to drink', 1, 'static/my_image01.jpg', '2021-03-15 18:38:30'),
(4, 333333, 'Very Impressive', 'Good to drink', 2, 'static/my_image01.jpg', '2021-04-15 18:38:30'),
(5, 444444, 'Pathetic', 'Good to drink', 2, 'static/my_image01.jpg', '2021-05-15 18:38:30'),
(6, 555555, 'Old Fashion', 'Good to drink', 3, 'static/my_image01.jpg', '2021-12-10 18:38:30'),
(7, 666666, 'Old Fashion', 'Good to drink', 3, 'static/my_image01.jpg', '2021-12-11 18:38:30'),
(8, 777777, 'Old Fashion', 'Good to drink', 4, 'static/my_image01.jpg', '2021-12-12 18:38:30'),
(9, 888888, 'Old Fashion', 'Good to drink', 5, 'static/my_image01.jpg', '2021-12-13 18:38:30'),
(10, 999999, 'Old Fashion', 'Good to drink', 6, 'static/my_image01.jpg', '2021-12-14 18:38:30'),
(11, 12121212, 'Old Fashion', 'Good to drink', 6, 'static/my_image01.jpg', '2021-12-15 18:38:30');

/*!40000 ALTER TABLE `cocktails` ENABLE KEYS */;
UNLOCK TABLES;



/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
