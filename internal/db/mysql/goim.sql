CREATE DATABASE goim;

set global time_zone = '+8:00';

USE goim;

-- START: USERS --
CREATE TABLE users (
                       `userId` bigint unsigned AUTO_INCREMENT NOT NULL,
                       `password` varchar(100) not null,
                       `role` varchar(10) not null,
                       `status` enum('ACTIVE', 'INACTIVE') not null,
                       `telephone` varchar(255) not null,
                       `email` varchar(255) not null,
                       `username` varchar(100) not null,
                       `nickname` varchar(100) not null,
                       `description` varbinary(255) null,
                       PRIMARY KEY (`userId`)
);
-- END

SELECT * FROM users;

INSERT INTO users(`userId`, `password`, `role`, `status`, `telephone`, `email`, `username`, `nickname`) VALUES (1, '1234567890', '1', 'ACTIVE', '13631210000', 'qwer@qq.com', 'TEST01', 'TEST01');
