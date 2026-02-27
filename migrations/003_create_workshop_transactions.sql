-- Create workshop_transactions table (labour + parts)
CREATE TABLE IF NOT EXISTS `workshop_transactions` (
  `id` BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `tenant_info_id` BIGINT NOT NULL,
  `date` DATE NOT NULL,
  `labour_net` DECIMAL(15,2) NOT NULL DEFAULT 0,
  `labour_vat` DECIMAL(15,2) NOT NULL DEFAULT 0,
  `parts_net` DECIMAL(15,2) NOT NULL DEFAULT 0,
  `parts_vat` DECIMAL(15,2) NOT NULL DEFAULT 0,
  `zero_rated_net` DECIMAL(15,2) NOT NULL DEFAULT 0 COMMENT 'Zero-rated work e.g. repairs',
  `customer_id` BIGINT,
  `invoice_id` BIGINT,
  `description` TEXT,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY `idx_tenant_info_id` (`tenant_info_id`),
  KEY `idx_date` (`date`),
  CONSTRAINT `fk_workshop_transactions_tenant` FOREIGN KEY (`tenant_info_id`) REFERENCES `tenant_infos` (`tenant_info_id`) ON DELETE CASCADE
);
