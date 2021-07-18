CREATE TABLE IF NOT EXISTS `doctor` (
  `id` INTEGER PRIMARY KEY,
  `name` VARCHAR(100) NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS `doctor_name_UNIQUE` ON `doctor` (`name` ASC);

CREATE TABLE IF NOT EXISTS `doctor_schedule` (
  `id` INTEGER PRIMARY KEY,
  `doctor_id` INT NOT NULL,
  `start_time` TIMESTAMP NULL,
  `end_time` TIMESTAMP NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS `schedule_doctor_id_INDEX` ON `doctor_schedule` (`doctor_id` ASC);

CREATE TABLE IF NOT EXISTS `patient` (
  `id` INTEGER PRIMARY KEY,
  `name` VARCHAR(100) NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS `patient_name_UNIQUE` ON `patient` (`name` ASC);

CREATE TABLE IF NOT EXISTS `appointments` (
  `id` INTEGER PRIMARY KEY,
  `doctor_id` INT NOT NULL,
  `patient_id` INT NOT NULL,
  `start_time` TIMESTAMP NOT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` TIMESTAMP NULL,
  `is_active` INT
);

CREATE INDEX IF NOT EXISTS `doctor_id_active_st_INDEX` ON `appointments` (`doctor_id` ASC, `is_active` ASC, `start_time` ASC);