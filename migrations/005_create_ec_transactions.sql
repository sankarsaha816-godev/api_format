-- Create ec_transactions table (EU acquisitions & supplies)
CREATE TABLE IF NOT EXISTS `ec_transactions` (
  `id` BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `tenant_info_id` BIGINT NOT NULL,
  `date` DATE NOT NULL,
  `type` ENUM('acquisition','supply') NOT NULL COMMENT 'Goods IN or OUT of EU',
  `net_amount` DECIMAL(15,2) NOT NULL,
  `vat_due` DECIMAL(15,2) NOT NULL,
  `customer_or_supplier` BIGINT,
  `customer_or_supplier_name` VARCHAR(255),
  `country_code` CHAR(2),
  `description` TEXT,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY `idx_tenant_info_id` (`tenant_info_id`),
  KEY `idx_date` (`date`),
  KEY `idx_type` (`type`),
  CONSTRAINT `fk_ec_transactions_tenant` FOREIGN KEY (`tenant_info_id`) REFERENCES `tenant_infos` (`tenant_info_id`) ON DELETE CASCADE
);
