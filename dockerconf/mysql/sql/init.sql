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


/*!40101 SET @OLD_CHARACTER_SET_CLIENT = @@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS = @@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION = @@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
SET NAMES utf8mb4;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS = @@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS = 0 */;
/*!40101 SET @OLD_SQL_MODE = @@SQL_MODE, SQL_MODE = 'NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES = @@SQL_NOTES, SQL_NOTES = 0 */;


# Dump of table social_accounts
# ------------------------------------------------------------

DROP TABLE IF EXISTS `social_accounts`;

CREATE TABLE `social_accounts`
(
    `id`           bigint(64)                             NOT NULL AUTO_INCREMENT,
    `social_id`    varchar(64) COLLATE utf8mb4_general_ci NOT NULL,
    `user_id`      bigint(20)                             NOT NULL,
    `type`         tinyint(1)                             NOT NULL DEFAULT '0',
    `created_date` timestamp                              NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_social_id` (`social_id`),
    UNIQUE KEY `idx_user_id` (`user_id`),
    KEY `idx_date` (`created_date`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci;



# Dump of table users
# ------------------------------------------------------------

DROP TABLE IF EXISTS `users`;

CREATE TABLE `users`
(
    `id`                   bigint(64)                              NOT NULL AUTO_INCREMENT,
    `account`              varchar(20) COLLATE utf8mb4_general_ci  NOT NULL,
    `password`             varchar(100) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
    `status`               tinyint(1)                              NOT NULL DEFAULT '0',
    `name`                 varchar(32) COLLATE utf8mb4_general_ci  NOT NULL DEFAULT '',
    `email`                varchar(64) COLLATE utf8mb4_general_ci  NOT NULL DEFAULT '',
    `photo`                varchar(128) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
    `number_of_post`       int(10) unsigned                        NOT NULL DEFAULT '0' COMMENT ' 貼文數',
    `number_of_collection` int(10) unsigned                        NOT NULL DEFAULT '0' COMMENT ' 收藏數',
    `number_of_draft`      int(10) unsigned                        NOT NULL DEFAULT '0' COMMENT ' 草稿數',
    `is_collection_public` tinyint(1)                              NOT NULL DEFAULT '0' COMMENT ' 是否公開收藏 0=不公開, 1=公開',
    `remark`               varchar(64) COLLATE utf8mb4_general_ci  NOT NULL DEFAULT '',
    `created_date`         timestamp                               NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_account` (`account`),
    KEY `idx_date` (`created_date`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci;

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users`
    DISABLE KEYS */;

INSERT INTO `users` (`id`, `account`, `password`, `status`, `name`, `email`, `photo`, `number_of_post`,
                     `number_of_collection`, `number_of_draft`, `is_collection_public`, `remark`, `created_date`)
VALUES (1, 'mockAccount1', '', 1, 'mockUser1', '', '', 0, 0, 0, 1, '', '2022-03-12 13:36:24'),
       (2, 'mockAccount2', '', 1, 'mockUser2', '', '', 0, 0, 0, 1, '', '2022-03-12 13:36:24'),
       (3, 'mockAccount3', '', 1, 'mockUser3', '', '', 0, 0, 0, 1, '', '2022-03-12 13:36:24'),
       (4, 'mockAccount4', '', 1, 'mockUser4', '', '', 0, 0, 0, 1, '', '2022-03-12 13:36:24'),
       (5, 'mockAccount5', '', 1, 'mockUser5', '', '', 0, 0, 0, 1, '', '2022-03-12 13:36:24'),
       (6, 'mockAccount6', '', 1, 'mockUser6', '', '', 0, 0, 0, 1, '', '2022-03-12 13:36:24');

/*!40000 ALTER TABLE `users`
    ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table cocktail_ingredients
# ------------------------------------------------------------

CREATE TABLE `cocktail_ingredients`
(
    `id`                bigint(64)                             NOT NULL AUTO_INCREMENT,
    `cocktail_id`       bigint(64)                             NOT NULL,
    `ingredient_name`   varchar(16) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT ' 成分名稱',
    `ingredient_amount` varchar(16)                            NOT NULL DEFAULT '' COMMENT '成分數量',
    `created_date`      timestamp                              NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_cocktail_id` (`cocktail_id`),
    KEY `idx_date` (`created_date`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci;

LOCK TABLES `cocktail_ingredients` WRITE;
/*!40000 ALTER TABLE `cocktail_ingredients`
    DISABLE KEYS */;

INSERT INTO `cocktail_ingredients` (`id`, `cocktail_id`, `ingredient_name`, `ingredient_amount`, `created_date`)
VALUES (1, 123456, 'gin tonic', '1 liter', '2022-03-12 13:36:24'),
       (2, 123456, 'cherry', '2 basket', '2022-03-12 13:36:24'),
       (3, 1111111, 'gin tonic', '1 liter', '2022-03-12 13:36:24'),
       (4, 1111111, 'cherry', '2 basket', '2022-03-12 13:36:24'),
       (5, 222222, 'gin tonic', '1 liter', '2022-03-12 13:36:24'),
       (6, 222222, 'cherry', '2 basket', '2022-03-12 13:36:24'),
       (7, 333333, 'gin tonic', '1 liter', '2022-03-12 13:36:24'),
       (8, 333333, 'cherry', '2 basket', '2022-03-12 13:36:24'),
       (9, 444444, 'gin tonic', '1 liter', '2022-03-12 13:36:24'),
       (10, 444444, 'cherry', '2 basket', '2022-03-12 13:36:24'),
       (11, 555555, 'gin tonic', '1 liter', '2022-03-12 13:36:24'),
       (12, 555555, 'cherry', '2 basket', '2022-03-12 13:36:24'),
       (13, 666666, 'gin tonic', '1 liter', '2022-03-12 13:36:24'),
       (14, 666666, 'cherry', '2 basket', '2022-03-12 13:36:24'),
       (15, 777777, 'gin tonic', '1 liter', '2022-03-12 13:36:24'),
       (16, 777777, 'cherry', '2 basket', '2022-03-12 13:36:24'),
       (17, 888888, 'gin tonic', '1 liter', '2022-03-12 13:36:24'),
       (18, 888888, 'cherry', '2 basket', '2022-03-12 13:36:24'),
       (19, 999999, 'gin tonic', '1 liter', '2022-03-12 13:36:24'),
       (20, 999999, 'cherry', '2 basket', '2022-03-12 13:36:24'),
       (21, 12121212, 'gin tonic', '1 liter', '2022-03-12 13:36:24'),
       (22, 12121212, 'cherry', '2 basket', '2022-03-12 13:36:24');

/*!40000 ALTER TABLE `cocktail_ingredients`
    ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table cocktail_steps
# ------------------------------------------------------------

CREATE TABLE `cocktail_steps`
(
    `id`               bigint(64)                             NOT NULL AUTO_INCREMENT,
    `cocktail_id`      bigint(64)                             NOT NULL,
    `step_number`      int(2) unsigned                        NOT NULL DEFAULT '1' COMMENT ' 步驟',
    `step_description` varchar(64) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT ' 步驟介紹',
    `created_date`     timestamp                              NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_cocktail_id` (`cocktail_id`),
    KEY `idx_date` (`created_date`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci;

LOCK TABLES `cocktail_steps` WRITE;
/*!40000 ALTER TABLE `cocktail_steps`
    DISABLE KEYS */;

INSERT INTO `cocktail_steps` (`id`, `cocktail_id`, `step_number`, `step_description`, `created_date`)
VALUES (1, 123456, 1, 'mix', '2022-03-12 13:34:56'),
       (2, 123456, 2, 'shack', '2022-03-12 13:35:13'),
       (3, 123456, 3, 'drink', '2022-03-12 13:44:06'),
       (4, 1111111, 1, 'mix', '2022-03-12 13:34:56'),
       (5, 1111111, 2, 'shack', '2022-03-12 13:35:13'),
       (6, 222222, 1, 'mix', '2022-03-12 13:34:56'),
       (7, 222222, 2, 'shack', '2022-03-12 13:35:13'),
       (8, 333333, 1, 'mix', '2022-03-12 13:34:56'),
       (9, 333333, 2, 'shack', '2022-03-12 13:35:13'),
       (10, 444444, 1, 'mix', '2022-03-12 13:34:56'),
       (11, 444444, 2, 'shack', '2022-03-12 13:35:13'),
       (12, 555555, 1, 'mix', '2022-03-12 13:34:56'),
       (13, 555555, 2, 'shack', '2022-03-12 13:35:13'),
       (14, 666666, 1, 'mix', '2022-03-12 13:34:56'),
       (15, 666666, 2, 'shack', '2022-03-12 13:35:13'),
       (16, 777777, 1, 'mix', '2022-03-12 13:34:56'),
       (17, 777777, 2, 'shack', '2022-03-12 13:35:13'),
       (18, 888888, 1, 'mix', '2022-03-12 13:34:56'),
       (19, 888888, 2, 'shack', '2022-03-12 13:35:13'),
       (20, 999999, 1, 'mix', '2022-03-12 13:34:56'),
       (21, 999999, 2, 'shack', '2022-03-12 13:35:13'),
       (22, 12121212, 1, 'mix', '2022-03-12 13:34:56'),
       (23, 12121212, 2, 'shack', '2022-03-12 13:35:13');

/*!40000 ALTER TABLE `cocktail_steps`
    ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table cocktail_steps
# ------------------------------------------------------------

CREATE TABLE `cocktail_photos`
(
    `id`             bigint(64)                              NOT NULL AUTO_INCREMENT,
    `cocktail_id`    bigint(64)                              NOT NULL,
    `photo`          varchar(128) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '1' COMMENT ' 照片',
    `is_cover_photo` tinyint(1)                              NOT NULL DEFAULT '0' COMMENT ' 是否為封面照 0=否, 1=是',
    `created_date`   timestamp                               NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_cocktail_id` (`cocktail_id`),
    KEY `idx_date` (`created_date`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci;

LOCK TABLES `cocktail_photos` WRITE;
/*!40000 ALTER TABLE `cocktail_photos`
    DISABLE KEYS */;

INSERT INTO `cocktail_photos` (`id`, `cocktail_id`, `photo`, `is_cover_photo`, `created_date`)
VALUES (1, 123456, 'static/my_image01.jpg', 1, '2022-02-11 13:42:48'),
       (2, 1111111, 'static/my_image01.jpg', 1, '2022-02-11 13:42:48'),
       (3, 222222, 'static/my_image01.jpg', 1, '2022-02-11 13:42:48'),
       (4, 333333, 'static/my_image01.jpg', 1, '2022-02-11 13:42:48'),
       (5, 444444, 'static/my_image01.jpg', 1, '2022-02-11 13:42:48'),
       (6, 555555, 'static/my_image01.jpg', 1, '2022-02-11 13:42:48'),
       (7, 666666, 'static/my_image01.jpg', 1, '2022-02-11 13:42:48'),
       (8, 777777, 'static/my_image01.jpg', 1, '2022-02-11 13:42:48'),
       (9, 888888, 'static/my_image01.jpg', 1, '2022-02-11 13:42:48'),
       (10, 999999, 'static/my_image01.jpg', 1, '2022-02-11 13:42:48'),
       (11, 12121212, 'static/my_image01.jpg', 1, '2022-02-11 13:42:48');

/*!40000 ALTER TABLE `cocktail_photos`
    ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table cocktails
# ------------------------------------------------------------

CREATE TABLE `cocktails`
(
    `id`           bigint(64)                              NOT NULL AUTO_INCREMENT,
    `cocktail_id`  bigint(64)                              NOT NULL,
    `user_id`      bigint(64)                              NOT NULL COMMENT ' 作者id',
    `title`        varchar(16) COLLATE utf8mb4_general_ci  NOT NULL COMMENT ' 調酒名稱',
    `description`  varchar(512) COLLATE utf8mb4_general_ci NOT NULL COMMENT ' 調酒介紹',
    `category`     tinyint(1)                              NOT NULL DEFAULT '0' COMMENT ' 類型 0=草稿, 1=正式',
    `created_date` timestamp                               NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_cocktail_id` (`cocktail_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_date` (`created_date`)
) ENGINE = InnoDB
  AUTO_INCREMENT = 12
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci;

LOCK TABLES `cocktails` WRITE;
/*!40000 ALTER TABLE `cocktails`
    DISABLE KEYS */;

INSERT INTO `cocktails` (`id`, `cocktail_id`, `title`, `description`, `user_id`, `category`, `created_date`)
VALUES (1, 123456, 'Side Car', 'Good to drink', 1, 1, '2021-01-15 18:38:30'),
       (2, 1111111, 'Old Fashion', 'Good to drink', 1, 1, '2021-02-15 18:38:30'),
       (3, 222222, 'Gin tonic', 'Good to drink', 1, 1, '2021-03-15 18:38:30'),
       (4, 333333, 'Very Impressive', 'Good to drink', 2, 1, '2021-04-15 18:38:30'),
       (5, 444444, 'Pathetic', 'Good to drink', 2, 1, '2021-05-15 18:38:30'),
       (6, 555555, 'Old Fashion', 'Good to drink', 3, 1, '2021-12-10 18:38:30'),
       (7, 666666, 'Old Fashion', 'Good to drink', 3, 1, '2021-12-11 18:38:30'),
       (8, 777777, 'Old Fashion', 'Good to drink', 4, 1, '2021-12-12 18:38:30'),
       (9, 888888, 'Old Fashion', 'Good to drink', 5, 1, '2021-12-13 18:38:30'),
       (10, 999999, 'Old Fashion', 'Good to drink', 6, 1, '2021-12-14 18:38:30'),
       (11, 12121212, 'Old Fashion', 'Good to drink', 6, 1, '2021-12-15 18:38:30');

/*!40000 ALTER TABLE `cocktails`
    ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table favorite_cocktails
# ------------------------------------------------------------

CREATE TABLE `favorite_cocktails`
(
    `id`           bigint(64) NOT NULL AUTO_INCREMENT,
    `cocktail_id`  bigint(64) NOT NULL COMMENT ' 調酒id',
    `user_id`      bigint(64) NOT NULL COMMENT ' 作者id',
    `created_date` timestamp  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_favorite_cocktail` (`user_id`, `cocktail_id`),
    KEY `idx_date` (`created_date`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci;



/*!40111 SET SQL_NOTES = @OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE = @OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS = @OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT = @OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS = @OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION = @OLD_COLLATION_CONNECTION */;
