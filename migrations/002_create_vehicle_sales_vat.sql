-- Create vehicle_sales table (margin scheme + VAT-qualifying car sales)
CREATE TABLE IF NOT EXISTS `vehicle_sales_vat` (
  `id` BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `tenant_info_id` BIGINT NOT NULL,
  `vehicle_id` BIGINT NOT NULL,
  `purchase_price` DECIMAL(15,2) NOT NULL,
  `selling_price` DECIMAL(15,2) NOT NULL,
  `margin` DECIMAL(15,2) NOT NULL COMMENT 'selling_price - purchase_price',
  `margin_vat` DECIMAL(15,2) NOT NULL COMMENT 'margin * 1/6',
  `vat_scheme` ENUM('margin','vat_qualifying','zero') NOT NULL,
  `sale_date` DATE NOT NULL,
  `customer_id` BIGINT,
  `invoice_id` BIGINT,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY `idx_tenant_info_id` (`tenant_info_id`),
  KEY `idx_vehicle_id` (`vehicle_id`),
  KEY `idx_sale_date` (`sale_date`),
  CONSTRAINT `fk_vehicle_sales_vat_tenant` FOREIGN KEY (`tenant_info_id`) REFERENCES `tenant_infos` (`tenant_info_id`) ON DELETE CASCADE
);
