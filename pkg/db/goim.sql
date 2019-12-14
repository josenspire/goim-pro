CREATE DATABASE goim;
USE goim;

-- START - for testing
CREATE TABLE tests (
    `id` bigint unsigned AUTO_INCREMENT NOT NULL,
    `name` varchar(100) NOT NULL,
    PRIMARY KEY (`id`)
);
INSERT INTO `tests` (`name`) VALUES('JAMES_01'),('JAMES_02');
SELECT * FROM tests;
-- END
