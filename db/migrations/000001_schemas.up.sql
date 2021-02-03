CREATE TABLE IF NOT EXISTS users
(
    id         varchar(36) PRIMARY KEY NOT NULL,
    username   varchar(45)             NOT NULL,
    password   varchar(97)             NOT NULL,
    email      varchar(45)             NOT NULL,
    created_at timestamp               NOT NULL DEFAULT now(),
    updated_at timestamp
);

CREATE TABLE IF NOT EXISTS profiles
(
    id        bigint PRIMARY KEY AUTO_INCREMENT NOT NULL,
    userId    varchar(36),
    INDEX `idxUser` (userId),
    CONSTRAINT `fkProfileId`
        FOREIGN KEY (userId)
            REFERENCES users (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    avatar    varchar(128) DEFAULT('') NOT NULL,
    firstName varchar(100) DEFAULT('') NOT NULL,
    lastName  varchar(100) DEFAULT('') NOT NULL
);

CREATE TABLE IF NOT EXISTS post_cards
(
    id     bigint PRIMARY KEY AUTO_INCREMENT NOT NULL,
    userId varchar(36),
    INDEX `idxUser` (userId),
    CONSTRAINT `fkMemePostId`
        FOREIGN KEY (userId)
            REFERENCES users (id) ON UPDATE CASCADE ON DELETE RESTRICT,
    image varchar(128) NOT NULL,
    wacks int(11) unsigned DEFAULT(0),
    cools int(11) unsigned DEFAULT(0),
    hearts int(11) unsigned DEFAULT(0)
)