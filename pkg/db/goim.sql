CREATE DATABASE goim;
USE goim;

-- START: USERS --
CREATE TABLE users (
    `id` bigint unsigned AUTO_INCREMENT NOT NULL,
    `name` varchar(100) NOT NULL,
    PRIMARY KEY (`id`)
);

-- END
