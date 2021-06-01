CREATE TABLE `merchant`
(
    `id`         bigint(20)   NOT NULL AUTO_INCREMENT,
    `name`       varchar(255) NOT NULL,
    `email`      varchar(255) NOT NULL,
    `password`   varchar(255) NOT NULL,
    `created_at` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime              DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `email` (`email`)
) ENGINE = InnoDB;


CREATE TABLE `product`
(
    `id`          bigint(20)   NOT NULL AUTO_INCREMENT,
    `name`        varchar(255) NOT NULL,
    `price`       float        NOT NULL,
    `created_at`  datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`  datetime              DEFAULT NULL,
    `merchant_id` bigint(20)   NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB;
