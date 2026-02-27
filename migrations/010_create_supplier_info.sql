-- Migration: create supplier_info table
CREATE TABLE IF NOT EXISTS `supplier_info` (
  `supplier_info_id` INT AUTO_INCREMENT PRIMARY KEY,
  `supplier_name` VARCHAR(150) NOT NULL,
  `first_name` VARCHAR(150) NOT NULL,
  `last_name` VARCHAR(150) NOT NULL,
  `email_id` VARCHAR(150) NOT NULL,
  `contact_no1` VARCHAR(100) NOT NULL,
  `contact_no2` VARCHAR(100) DEFAULT NULL,
  `tenant_info_id` INT NOT NULL,
  `tenant_name` VARCHAR(100) NOT NULL,
  `address_line1` TEXT NOT NULL,
  `address_line2` TEXT,
  `address_line3` TEXT,
  `city` VARCHAR(100) NOT NULL,
  `state` VARCHAR(100) NOT NULL,
  `country` VARCHAR(100) NOT NULL,
  `post_code` VARCHAR(50) NOT NULL,
  `status` INT DEFAULT 1,
  `created_by` VARCHAR(50) NOT NULL,
  `updated_by` VARCHAR(50) NOT NULL,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
