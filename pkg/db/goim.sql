CREATE DATABASE goim;
USE goim;

-- START: USERS --
CREATE TABLE users (
    `id` bigint unsigned AUTO_INCREMENT NOT NULL,
    `name` varchar(100) NOT NULL,
    PRIMARY KEY (`id`)
);

-- END

INSERT INTO users(`userId`, `password`, `role`, `status`, `telephone`, `email`, `username`, `nickname`) VALUES (1, '1234567890', '0', '0', '13631210000', 'qwer@qq.com', 'TEST01', 'TEST01');
