CREATE TABLE IF NOT EXISTS `users` (
    `id` varchar(50) NOT NULL,
    `email` varchar(255) NOT NULL,
    `password` varchar(255) NOT NULL,
    `address` varchar(255) NOT NULL DEFAULT '',
    `deleted_time` bigint(20) unsigned DEFAULT NULL,
    `created_time` bigint(20) unsigned NOT NULL DEFAULT '0',
    `updated_time` bigint(20) unsigned NOT NULL DEFAULT '0',
    PRIMARY KEY (`id`),
    KEY `email_idx` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
