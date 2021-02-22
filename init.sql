-- MySQL Workbench Forward Engineering

SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION';

-- -----------------------------------------------------
-- Schema projja
-- -----------------------------------------------------
DROP SCHEMA IF EXISTS `projja` ;

-- -----------------------------------------------------
-- Schema projja
-- -----------------------------------------------------
CREATE SCHEMA IF NOT EXISTS `projja` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci ;
USE `projja` ;

-- -----------------------------------------------------
-- Table `projja`.`users`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `projja`.`users` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(45) NOT NULL,
  `username` VARCHAR(45) NOT NULL,
  `telegram_id` VARCHAR(45) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `username_UNIQUE` (`username` ASC) VISIBLE,
  UNIQUE INDEX `telegram_id_UNIQUE` (`telegram_id` ASC) VISIBLE)
ENGINE = InnoDB
AUTO_INCREMENT = 4
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_0900_ai_ci;


-- -----------------------------------------------------
-- Table `projja`.`project`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `projja`.`project` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(45) NOT NULL,
  `owner` INT NOT NULL,
  `status` ENUM('opened', 'closed') NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_project_1_idx` (`owner` ASC) VISIBLE,
  CONSTRAINT `fk_project_1`
    FOREIGN KEY (`owner`)
    REFERENCES `projja`.`users` (`id`)
    ON DELETE CASCADE
    ON UPDATE RESTRICT)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_0900_ai_ci;


-- -----------------------------------------------------
-- Table `projja`.`member`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `projja`.`member` (
  `users` INT NOT NULL,
  `project` INT NOT NULL,
  INDEX `fk_member_1_idx` (`users` ASC) VISIBLE,
  INDEX `fk_member_2_idx` (`project` ASC) VISIBLE,
  CONSTRAINT `fk_member_1`
    FOREIGN KEY (`users`)
    REFERENCES `projja`.`users` (`id`)
    ON DELETE CASCADE
    ON UPDATE RESTRICT,
  CONSTRAINT `fk_member_2`
    FOREIGN KEY (`project`)
    REFERENCES `projja`.`project` (`id`)
    ON DELETE CASCADE
    ON UPDATE RESTRICT)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_0900_ai_ci;


-- -----------------------------------------------------
-- Table `projja`.`skill`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `projja`.`skill` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `skill` VARCHAR(45) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `skill_UNIQUE` (`skill` ASC) VISIBLE)
ENGINE = InnoDB
AUTO_INCREMENT = 5
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_0900_ai_ci;


-- -----------------------------------------------------
-- Table `projja`.`task_status`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `projja`.`task_status` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `status` VARCHAR(45) NOT NULL,
  `level` INT NOT NULL,
  `project` INT NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_task_status_1_idx` (`project` ASC) VISIBLE,
  CONSTRAINT `fk_task_status_1`
    FOREIGN KEY (`project`)
    REFERENCES `projja`.`project` (`id`)
    ON DELETE CASCADE
    ON UPDATE RESTRICT)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_0900_ai_ci;


-- -----------------------------------------------------
-- Table `projja`.`task`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `projja`.`task` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `description` VARCHAR(200) NOT NULL,
  `project` INT NOT NULL,
  `deadline` DATETIME NOT NULL,
  `priority` ENUM('critical', 'high', 'medium', 'low') NOT NULL,
  `status` INT NOT NULL,
  `is_closed` TINYINT NOT NULL,
  `executor` INT NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_task_1_idx` (`project` ASC) VISIBLE,
  INDEX `fk_task_2_idx` (`status` ASC) VISIBLE,
  INDEX `fk_task_3_idx` (`executor` ASC) VISIBLE,
  CONSTRAINT `fk_task_1`
    FOREIGN KEY (`project`)
    REFERENCES `projja`.`project` (`id`)
    ON DELETE CASCADE
    ON UPDATE RESTRICT,
  CONSTRAINT `fk_task_2`
    FOREIGN KEY (`status`)
    REFERENCES `projja`.`task_status` (`id`)
    ON DELETE CASCADE
    ON UPDATE RESTRICT,
  CONSTRAINT `fk_task_3`
    FOREIGN KEY (`executor`)
    REFERENCES `projja`.`users` (`id`)
    ON DELETE CASCADE
    ON UPDATE RESTRICT)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_0900_ai_ci;


-- -----------------------------------------------------
-- Table `projja`.`task_skill`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `projja`.`task_skill` (
  `task` INT NOT NULL,
  `skill` INT NOT NULL,
  INDEX `fk_task_skill_1_idx` (`task` ASC) VISIBLE,
  INDEX `fk_task_skill_2_idx` (`skill` ASC) VISIBLE,
  CONSTRAINT `fk_task_skill_1`
    FOREIGN KEY (`task`)
    REFERENCES `projja`.`task` (`id`)
    ON DELETE CASCADE
    ON UPDATE RESTRICT,
  CONSTRAINT `fk_task_skill_2`
    FOREIGN KEY (`skill`)
    REFERENCES `projja`.`skill` (`id`)
    ON DELETE CASCADE
    ON UPDATE RESTRICT)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_0900_ai_ci;


-- -----------------------------------------------------
-- Table `projja`.`users_skill`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `projja`.`users_skill` (
  `users` INT NOT NULL,
  `skill` INT NOT NULL,
  INDEX `fk_users_skill_1_idx` (`users` ASC) VISIBLE,
  INDEX `fk_users_skill_2_idx` (`skill` ASC) VISIBLE,
  CONSTRAINT `fk_users_skill_1`
    FOREIGN KEY (`users`)
    REFERENCES `projja`.`users` (`id`)
    ON DELETE CASCADE
    ON UPDATE RESTRICT,
  CONSTRAINT `fk_users_skill_2`
    FOREIGN KEY (`skill`)
    REFERENCES `projja`.`skill` (`id`)
    ON DELETE CASCADE
    ON UPDATE RESTRICT)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_0900_ai_ci;


SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;
