CREATE TABLE IF NOT EXISTS users
(
    id         bigint PRIMARY KEY AUTO_INCREMENT NOT NULL,
    username   varchar(45)                       NOT NULL,
    password   varchar(97)                       NOT NULL,
    email      varchar(45)                       NOT NULL,
    created_at timestamp                         NOT NULL DEFAULT now(),
    updated_at timestamp
);

CREATE TABLE IF NOT EXISTS profiles
(
    id     bigint PRIMARY KEY AUTO_INCREMENT NOT NULL,
    userId bigint,
    INDEX `idxUser`(userId),
    CONSTRAINT `fkProfileId`
    FOREIGN KEY (userId)
    REFERENCES users(id) ON UPDATE CASCADE ON DELETE RESTRICT
);