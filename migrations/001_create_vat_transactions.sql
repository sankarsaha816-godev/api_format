-- Create vat_transactions table (core transaction log for all VAT)
CREATE TABLE IF NOT EXISTS `vat_transactions` (
  `id` BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `tenant_info_id` BIGINT NOT NULL,
  `date` DATE NOT NULL,
  `transaction` VARCHAR(50) NOT NULL COMMENT 'sales, service, acquisition, overhead, cost, purchase',
  `vat_scheme` ENUM('standard','margin','zero','exempt') NOT NULL DEFAULT 'standard',
  `gross_amount` DECIMAL(15,2) NOT NULL DEFAULT 0,
  `net_amount` DECIMAL(15,2) NOT NULL DEFAULT 0,
  `purchase_price` DECIMAL(15,2) DEFAULT 0,
  `margin` DECIMAL(15,2) DEFAULT 0,
  `vat_amount` DECIMAL(15,2) NOT NULL DEFAULT 0,
  `name` VARCHAR(255) COMMENT 'regno for sales/service, overhead category for OH',
  `detail` VARCHAR(255) COMMENT 'stock num for sales, OH detail for OH',
  `description` TEXT,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY `idx_tenant_info_id` (`tenant_info_id`),
  KEY `idx_date` (`date`),
  KEY `idx_transaction` (`transaction`),
  CONSTRAINT `fk_vat_transactions_tenant` FOREIGN KEY (`tenant_info_id`) REFERENCES `tenant_infos` (`tenant_info_id`) ON DELETE CASCADE
);
