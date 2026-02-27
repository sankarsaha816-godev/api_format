-- Create purchases table (VAT reclaim)
CREATE TABLE IF NOT EXISTS `vat_purchases` (
  `id` BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `tenant_info_id` BIGINT NOT NULL,
  `date` DATE NOT NULL,
  `supplier_id` BIGINT,
  `supplier_name` VARCHAR(255),
  `net_amount` DECIMAL(15,2) NOT NULL,
  `vat_amount` DECIMAL(15,2) NOT NULL,
  `gross_amount` DECIMAL(15,2) NOT NULL,
  `reclaimable` BOOLEAN DEFAULT TRUE,
  `category` ENUM('car_purchase','parts','supplies','overheads','other') DEFAULT 'other',
  `description` TEXT,
  `invoice_number` VARCHAR(100),
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY `idx_tenant_info_id` (`tenant_info_id`),
  KEY `idx_date` (`date`),
  KEY `idx_reclaimable` (`reclaimable`),
  CONSTRAINT `fk_vat_purchases_tenant` FOREIGN KEY (`tenant_info_id`) REFERENCES `tenant_infos` (`tenant_info_id`) ON DELETE CASCADE
);
