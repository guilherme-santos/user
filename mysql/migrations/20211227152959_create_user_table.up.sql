CREATE TABLE `user` (
  `id` CHAR(20) NOT NULL,
  `first_name` VARCHAR(100) NOT NULL,
  `last_name` VARCHAR(100) NOT NULL,
  `nickname` VARCHAR(100) NOT NULL,
  `password` VARCHAR(255) NOT NULL,
  `email` VARCHAR(255) NOT NULL,
  `country` CHAR(2) NOT NULL,
  `created_at` TIMESTAMP(6)
    NOT NULL
    DEFAULT CURRENT_TIMESTAMP(6),
  `updated_at` TIMESTAMP(6)
    ON UPDATE CURRENT_TIMESTAMP(6)
    NOT NULL
    DEFAULT CURRENT_TIMESTAMP(6),
  `removed_at` TIMESTAMP(6) NULL,
  PRIMARY KEY (`id`),
  UNIQUE (`email`)
) ENGINE = InnoDB;
